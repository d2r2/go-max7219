Foreword
--------
This project is mainly a fork of respective functionality originally written by Richard Hull in python: <https://github.com/rm-hull/max7219>. Newetheless it differs in some parts: refuse some functionality (works only with matrix led), include extra functionality (extra fonts, support of national languages).

MAX7219 driver
--------------
This library intended to output text messages to 8x8 LED display (pdf reference) via MAX7219 driver chip (pdf reference):

Today there are many master devices which can drive MAX7219 chip, but this lib intended to work on Raspberry PI and theire clones (tested with Raspberry PI and Banana PI). It may works with any Raspberry PI clone, which support Kernel SPI bus, and you should carry out all necessary preparations to make SPI bus device present in /dev/ list.

Golang usage
------------
```go
func main() {
  dev.SlideMessage("Hello world :)", font, true, 50*time.Millisecond)
}
```
