/*
SMS Palette
A simple program that generates a .PNG image file and a .GPL palette file, for GIMP importing. of the 64-color Sega Master System palette.

Copyright 2018 Brian Puthuff

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
)

const BLOCK_DIM = 64
const ROWS = 8
const COLS = 8

func main() {
	var out_image *image.RGBA
	var out_png *os.File
	var out_gpl *os.File
	var err error
	var bounds image.Rectangle

	// create output image
	bounds = image.Rect(0, 0, BLOCK_DIM*COLS, BLOCK_DIM*ROWS)
	out_image = image.NewRGBA(bounds)

	// draw image (Sega Master System Palette)
	drawPalette(out_image)

	// create and save output file as .PNG format
	out_png, err = os.Create("SegaMasterSystem_Palette.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out_png.Close()
	err = png.Encode(out_png, out_image)
	if err != nil {
		out_png.Close()
		log.Fatal(err)
	}

	// create palette file for GIMP (.gpl)
	out_gpl, err = os.Create("SegaMasterSystem.gpl")
	if err != nil {
		log.Fatal(err)
	}
	defer out_gpl.Close()
	writePalette(out_gpl)
}

func writePalette(f *os.File) {
	var r, g, b uint8
	var rs, gs, bs string

	f.WriteString("GIMP Palette\n")
	f.WriteString("Name: Sega Master System\n")
	f.WriteString("#\n")
	for b = 0; b < 4; b++ {
		for g = 0; g < 4; g++ {
			for r = 0; r < 4; r++ {
				rs = fmt.Sprintf("%3d", r*85)
				gs = fmt.Sprintf("%3d", g*85)
				bs = fmt.Sprintf("%3d", b*85)
				f.WriteString(strings.Join([]string{rs, gs, bs}, " "))
				f.WriteString("\n")
			}
		}
	}
}

func drawPalette(i *image.RGBA) {
	var x, y int
	var pixel_color color.Color
	var pixel image.Rectangle

	// initialize drawing rectangle (pixel)
	pixel = image.Rect(0, 0, BLOCK_DIM, BLOCK_DIM)

	// main draw loop
	for y = 0; y < ROWS; y++ {
		for x = 0; x < COLS; x++ {
			pixel_color = getColor(x, y)
			pixel.Min.X = x * BLOCK_DIM
			pixel.Max.X = x*BLOCK_DIM + BLOCK_DIM
			pixel.Min.Y = y * BLOCK_DIM
			pixel.Max.Y = y*BLOCK_DIM + BLOCK_DIM
			draw.Draw(i, pixel, &image.Uniform{pixel_color}, image.ZP, draw.Src)
		}
	}
}

func getColor(x, y int) color.Color {
	var r, g, b uint8

	/*
		There are four quadrants each of size 4 x 4.
		For each quadrant the RGB values are organized as follows.

		    r00 r01 r10 r11
		g00 bxx bxx bxx bxx
		g01 bxx bxx bxx bxx
		g10 bxx bxx bxx bxx
		g11 bxx bxx bxx bxx

		Blue component values per each quandrant:

		[Q1 | b00] [Q2 | b01]
		[Q3 | b10] [Q4 | b11]
	*/
	r = uint8(x % 4)
	g = uint8(y % 4)
	if x < COLS/2 {
		if y < ROWS/2 {
			b = 0
		} else {
			b = 2
		}
	} else {
		if y < ROWS/2 {
			b = 1
		} else {
			b = 3
		}
	}
	r, g, b = r*85, g*85, b*85
	return color.RGBA{r, g, b, 0xFF}
}
