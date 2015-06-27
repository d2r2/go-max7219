package main

import (
	"log"
	"time"

	"github.com/d2r2/go-max7219/max7219"
)

func main() {
	mtx := max7219.NewMatrix(1)
	err := mtx.Open(7)
	if err != nil {
		log.Fatal(err)
	}
	defer mtx.Close()
	var font max7219.Font
	// font = max7219.FontCP437
	// font = max7219.FontLCD
	// font = max7219.FontBoldCyrillic
	// font = max7219.FontMSXRus
	// font = max7219.FontZXSpectrumRus
	// font = max7219.FontVestaPK8000Rus
	// font = max7219.FontSinclair
	// font = max7219.FontTiny
	//dev.SlideMessage("Привет мир!!! Шарик - ты балбес!!!",
	//	font, true, 50*time.Millisecond)
	// time.Sleep(3 * time.Second)
	font = max7219.FontBoldCyrillic
	mtx.SlideMessage("Щоб тебе підняло та гепнуло.",
		font, true, 50*time.Millisecond)
	mtx.SlideMessage("Столик на одного человека/двух человек, пожалуйста.",
		font, true, 50*time.Millisecond)
}
