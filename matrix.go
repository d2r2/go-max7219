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
	Device *Device
}

func NewMatrix(cascaded int) *Matrix {
	this := &Matrix{}
	this.Device = NewDevice(cascaded)
	return this
}

func (this *Matrix) Open(spibus int, spidevice int, brightness byte) error {
	return this.Device.Open(spibus, spidevice, brightness)
}

func (this *Matrix) Close() {
	this.Device.Close()
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
	var patrns [][]byte
	var limits [][]int
	totalWidth := 0
	for _, c := range text {
		pattern := font[c]
		start, end := getLetterPatternLimits(pattern)
		totalWidth += end - start + 1
		patrns = append(patrns, pattern)
		limits = append(limits, []int{start, end})
	}
	averageWidth := totalWidth / len(text)
	log.Debug("Average width: %d\n", averageWidth)
	var buf []byte
	for i := 0; i < len(patrns); i++ {
		if condenseLetterPattern {
			var startC = getLineCondense(patrns[i][limits[i][0]])
			var endC int = 0
			if i > 0 {
				endC = getLineCondense(patrns[i-1][limits[i-1][1]])
			}
			// In case of space char...
			if isEmpty(patrns[i]) {
				// ... specify average char width + extra line.
				limits[i][1] = averageWidth - 1 - 1
			}
			// ... + extra lines.
			if endC+startC == 0 {
			} else if endC+startC <= 2 || i == 0 {
				buf = append(buf, 0)
			} else if endC+startC <= 10 {
				buf = append(buf, 0, 0)
			} else {
				buf = append(buf, 0, 0, 0)
			}
			buf = append(buf, patrns[i][limits[i][0]:limits[i][1]+1]...)
		} else {
			buf = append(buf, patrns[i]...)
		}
	}
	if condenseLetterPattern {
		buf = append(buf, 0)
	}
	return buf
}

func repeat(b byte, count int) []byte {
	buf := make([]byte, count)
	for i := 0; i < len(buf); i++ {
		buf[i] = b
	}
	return buf
}

func isEmpty(pattern []byte) bool {
	for _, b := range pattern {
		if b != 0 {
			return false
		}
	}
	return true
}

// Output unicode char to the led matrix.
// Unicode char transforms to ascii code based on
// information taken from font.GetCodePage() call.
func (this *Matrix) OutputChar(cascadeId int, font Font,
	char rune, redraw bool) error {
	text := string(char)
	b := convertUnicodeToAscii(text, font.GetCodePage())
	if len(b) != 1 {
		return fmt.Errorf("One char expected: \"%s\"", text)
	}
	buf := preparePatterns(b, font.GetLetterPatterns(),
		false)
	for i, value := range buf {
		//fmt.Printf("value: %#x\n", value)
		err := this.Device.SetBufferLine(cascadeId, i, value, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

// Output ascii code to the led matrix.
func (this *Matrix) OutputAsciiCode(cascadeId int, font Font,
	asciiCode int, redraw bool) error {
	patterns := font.GetLetterPatterns()
	b := patterns[asciiCode]
	for i, value := range b {
		//fmt.Printf("value: %#x\n", value)
		err := this.Device.SetBufferLine(cascadeId, i, value, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

// Convert unicode text to ASCII text
// using specific codepage mapping.
func convertUnicodeToAscii(text string,
	codepage encoding.Encoding) []byte {
	b := []byte(text)
	// fmt.Printf("Text length: %d\n", len(b))
	var buf bytes.Buffer
	if codepage == nil {
		codepage = charmap.Windows1252
	}
	w := transform.NewWriter(&buf, codepage.NewEncoder())
	defer w.Close()
	w.Write(b)
	// fmt.Printf("Buffer length: %d\n", len(buf.Bytes()))
	return buf.Bytes()
}

// Show message sliding it by led matrix from the right to left.
func (this *Matrix) SlideMessage(text string, font Font,
	condensePattern bool, pixelDelay time.Duration) error {
	b := convertUnicodeToAscii(text, font.GetCodePage())
	buf := preparePatterns(b, font.GetLetterPatterns(),
		condensePattern)
	for _, b := range buf {
		time.Sleep(pixelDelay)
		err := this.Device.ScrollLeft(true)
		if err != nil {
			return err
		}
		err = this.Device.SetBufferLine(
			this.Device.GetCascadeCount()-1,
			this.Device.GetLedLineCount()-1, b, true)
		if err != nil {
			return err
		}
	}
	return nil
}
