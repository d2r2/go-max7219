package main

import (
	"log"
	"time"

	"github.com/d2r2/go-max7219/max7219"
)

func main() {
	mtx := max7219.NewMatrix(1)
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
	font = max7219.FontCP437
	mtx.SlideMessage("Hello world!!! Hey Beavis, let's rock!",
		font, true, 50*time.Millisecond)
	time.Sleep(1 * time.Second)
	font = max7219.FontZXSpectrumRus
	mtx.SlideMessage("Привет мир!!! Шарик - ты балбес!!!",
		font, true, 50*time.Millisecond)
}
