package main

import (
	"log"
	"time"

	"github.com/d2r2/go-max7219"
)

func main() {
	// Create new LED matrix with number of cascaded devices is equal to 1.
	mtx := max7219.NewMatrix(3)
	// Open SPI device with spibus and spidev equal to 0 and 0.
	// Set brightness equal to 7.
	err := mtx.Open(0, 0, 7)
	if err != nil {
		log.Fatal(err)
	}
	defer mtx.Close()
	var font max7219.Font
	// font = max7219.FontCP437
	// font = max7219.FontLCD
	// font = max7219.FontMSXRus
	// font = max7219.FontZXSpectrumRus
	// font = max7219.FontSinclair
	// font = max7219.FontTiny
	// font = max7219.FontVestaPK8000Rus
	font = max7219.FontLCD
	// Output text message to LED matrix.
	mtx.SlideMessage("Hello world!!! Hey Beavis, let's rock!",
		font, true, 50*time.Millisecond)
	// Wait 1 sec, then continue output new text.
	time.Sleep(1 * time.Second)
	font = max7219.FontZXSpectrumRus
	// Output national text to LED matrix.
	mtx.SlideMessage("Привет мир!!! Шарик, ты - балбес!!!",
		font, true, 50*time.Millisecond)
	/*
		var b byte = 0x55
		for i := 0; i < mtx.Device.GetCascadeCount(); i++ {
			mtx.Device.SetBufferLine(i, 0, b, true)
			b ^= 0xFF
			time.Sleep(3 * time.Second)
		}
		mtx.Device.ScrollLeft(true)
		time.Sleep(3 * time.Second)
		mtx.Device.ScrollLeft(true)
		time.Sleep(3 * time.Second)
		mtx.Device.ScrollLeft(true)
		time.Sleep(3 * time.Second)
	*/
}
