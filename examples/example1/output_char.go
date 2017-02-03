package main

import (
	"fmt"
	"log"
	"time"

	"github.com/d2r2/go-max7219"
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

	// Output a sequence of ascii codes in a loop
	font = max7219.FontCP437
	fmt.Printf("Char count: %d\n", len(font.GetLetterPatterns()))
	for i := 0; i < len(font.GetLetterPatterns()); i++ {
		mtx.OutputAsciiCode(0, font, i, true)
		time.Sleep(100 * time.Millisecond)
	}

	// Output non-latin national char (russian for example).
	// You must be sure, that your national char should match font code page.
	mtx.OutputChar(0, max7219.FontZXSpectrumRus, 'Я', true)
}
