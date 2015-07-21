package max7219

import (
	"fmt"

	"golang.org/x/text/encoding"

	"github.com/fulr/spidev"
)

// General interface of ASCII char set bit pattern
// for drawing on the LED matrix.
type Font interface {
	// Return font code page.
	// This function allow implement national font support.
	GetCodePage() encoding.Encoding
	// Return font char's bit pattern.
	// Font height is always equal to 8 pixel.
	// Font width may vary from one font
	// to another, but ordinary not exceed 8 pixel.
	GetLetterPatterns() [][]byte
}

type Max7219Reg byte

const (
	MAX7219_REG_NOOP   Max7219Reg = 0
	MAX7219_REG_DIGIT0            = iota
	MAX7219_REG_DIGIT1
	MAX7219_REG_DIGIT2
	MAX7219_REG_DIGIT3
	MAX7219_REG_DIGIT4
	MAX7219_REG_DIGIT5
	MAX7219_REG_DIGIT6
	MAX7219_REG_DIGIT7
	MAX7219_REG_DECODEMODE
	MAX7219_REG_INTENSITY
	MAX7219_REG_SCANLIMIT
	MAX7219_REG_SHUTDOWN
	MAX7219_REG_DISPLAYTEST = 0x0F
	MAX7219_REG_LASTDIGIT   = MAX7219_REG_DIGIT7
)

const MAX7219_DIGIT_COUNT = MAX7219_REG_LASTDIGIT -
	MAX7219_REG_DIGIT0 + 1

type Device struct {
	cascaded int
	buffer   []byte
	spi      *spidev.SPIDevice
}

func NewDevice(cascaded int) *Device {
	buf := make([]byte, MAX7219_DIGIT_COUNT*cascaded)
	this := &Device{cascaded: cascaded, buffer: buf}
	return this
}

func (this *Device) GetCascadeCount() int {
	return this.cascaded
}

func (this *Device) GetLedLineCount() int {
	return MAX7219_DIGIT_COUNT
}

func (this *Device) Open(spibus int, spidevice int, brightness byte) error {
	devstr := fmt.Sprintf("/dev/spidev%d.%d", spibus, spidevice)
	spi, err := spidev.NewSPIDevice(devstr)
	if err != nil {
		return err
	}
	this.spi = spi
	// Initialize Max7219 led driver.
	this.Command(MAX7219_REG_SCANLIMIT, 7)   // show all 8 digits
	this.Command(MAX7219_REG_DECODEMODE, 0)  // use matrix (not digits)
	this.Command(MAX7219_REG_DISPLAYTEST, 0) // no display test
	this.Command(MAX7219_REG_SHUTDOWN, 1)    // not shutdown mode
	this.Brightness(brightness)
	this.ClearAll(true)
	return nil
}

func (this *Device) Close() {
	this.spi.Close()
}

func (this *Device) Brightness(intensity byte) error {
	return this.Command(MAX7219_REG_INTENSITY, intensity)
}

func (this *Device) Command(reg Max7219Reg, value byte) error {
	buf := []byte{byte(reg), value}
	for i := 0; i < this.cascaded; i++ {
		_, err := this.spi.Xfer(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) sendBufferLine(position int) error {
	reg := MAX7219_REG_DIGIT0 + position
	//fmt.Printf("Register: %#x\n", reg)
	buf := make([]byte, this.cascaded*2)
	for i := 0; i < this.cascaded; i++ {
		b := this.buffer[i*MAX7219_DIGIT_COUNT+position]
		//fmt.Printf("Buffer value: %#x\n", b)
		buf[i*2] = byte(reg)
		buf[i*2+1] = b
	}
	log.Debug("Send to bus: %v\n", buf)
	_, err := this.spi.Xfer(buf)
	if err != nil {
		return err
	}
	return nil
}

func (this *Device) SetBufferLine(cascadeId int,
	position int, value byte, redraw bool) error {
	this.buffer[cascadeId*MAX7219_DIGIT_COUNT+position] = value
	if redraw {
		err := this.sendBufferLine(position)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) Flush() error {
	for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
		err := this.sendBufferLine(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) Clear(cascadeId int, redraw bool) error {
	if cascadeId >= 0 {
		for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
			this.buffer[cascadeId*MAX7219_DIGIT_COUNT+i] = 0
		}
	} else {
		for i := 0; i < this.cascaded*MAX7219_DIGIT_COUNT; i++ {
			this.buffer[i] = 0
		}
	}
	if redraw {
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) ClearAll(redraw bool) error {
	for i := 0; i < this.cascaded; i++ {
		err := this.Clear(i, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) ScrollLeft(redraw bool) error {
	this.buffer = append(this.buffer[1:], 0)
	log.Debug("Buffer: %v\n", this.buffer)
	if redraw {
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) ScrollRight(redraw bool) error {
	this.buffer = append([]byte{0}, this.buffer[:len(this.buffer)-1]...)
	if redraw {
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}
