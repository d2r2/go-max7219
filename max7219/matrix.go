package max7219

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Matrix struct {
	device *Device
}

func NewMatrix(cascaded int) *Matrix {
	this := &Matrix{}
	this.device = NewDevice(cascaded)
	return this
}

func (this *Matrix) Open(brightness byte) error {
	return this.device.Open(brightness)
}

func (this *Matrix) Close() {
	this.device.Close()
}

func getLineCondense(line byte) int {
	var condense int = 0
	var i uint
	for i = 0; i < 8; i++ {
		if line&(1<<i) > 0 {
			condense += 1
		}
	}
	return condense
}

func getLetterPatternLimits(pattern []byte) (start int, end int) {
	startIndex := -1
	endIndex := -1
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != 0 && startIndex == -1 {
			startIndex = i
		}
	}
	if startIndex == -1 {
		startIndex = 0
	}
	for i := len(pattern) - 1; i >= 0; i-- {
		if pattern[i] != 0 && endIndex == -1 {
			endIndex = i
		}
	}
	if endIndex == -1 {
		endIndex = len(pattern) - 1
	}
	return startIndex, endIndex
}

func preparePatterns(text []byte, font [][]byte,
	condenseLetterPattern bool) []byte {
	var temp [][]byte
	var limits [][]int
	for _, c := range text {
		fmt.Printf("Letter: %d\n", c)
		pattern := font[c]
		start, end := getLetterPatternLimits(pattern)
		temp = append(temp, pattern)
		limits = append(limits, []int{start, end})
	}
	var buf []byte
	for i := 0; i < len(temp); i++ {
		if condenseLetterPattern {
			if i == 0 {
				buf = append(buf, 0)
			} else {
				endC := getLineCondense(temp[i-1][limits[i-1][1]])
				startC := getLineCondense(temp[i][limits[i][0]])
				if endC+startC == 0 {
				} else if endC+startC <= 2 {
					buf = append(buf, 0)
				} else if endC+startC <= 10 {
					buf = append(buf, 0, 0)
				} else {
					buf = append(buf, 0, 0, 0)
				}
			}
			buf = append(buf, temp[i][limits[i][0]:limits[i][1]+1]...)
		} else {
			buf = append(buf, temp[i]...)
		}
	}
	if condenseLetterPattern {
		buf = append(buf, 0)
	}
	return buf
}

func (this *Matrix) Letter(device int, font [][]byte,
	asciiCode byte, redraw bool) error {
	for i, value := range font[asciiCode] {
		//fmt.Printf("value: %#x\n", value)
		err := this.device.SetBufferLine(device, i, value, redraw)
		if err != nil {
			return err
		}
		// time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func convertUnicodeToAscii(text string,
	codepage encoding.Encoding) []byte {
	b := []byte(text)
	fmt.Printf("Text length: %d\n", len(b))
	var buf bytes.Buffer
	if codepage == nil {
		codepage = charmap.Windows1252
	}
	w := transform.NewWriter(&buf, codepage.NewEncoder())
	defer w.Close()
	w.Write(b)
	fmt.Printf("Buffer length: %d\n", len(buf.Bytes()))
	return buf.Bytes()
}

func (this *Matrix) SlideMessage(text string, font Font,
	condensePattern bool, pixelDelay time.Duration) error {
	b := convertUnicodeToAscii(text, font.GetCodePage())
	buf := preparePatterns(b, font.GetLetterPatterns(),
		condensePattern)
	for _, b := range buf {
		time.Sleep(pixelDelay)
		err := this.device.ScrollLeft(true)
		if err != nil {
			return err
		}
		err = this.device.SetBufferLine(
			this.device.GetCascadeCount()-1,
			this.device.GetLedLineCount()-1, b, true)
		if err != nil {
			return err
		}
	}
	return nil
}
