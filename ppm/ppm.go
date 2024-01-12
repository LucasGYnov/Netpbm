package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Pixel struct {
	R, G, B uint8
}

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

func main() {
	var file string
	fmt.Print("Entrer le nom du document (attention à ne pas faire d'erreur) : ")
	fmt.Scanf("%s", &file)

	ppm, err := ReadPPM(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	ppm.Size()

	pixelValue := ppm.At(2, 2)
	fmt.Printf("Valeur du pixel à l'indice (2, 2): (%d, %d, %d)\n", pixelValue.R, pixelValue.G, pixelValue.B)

	newPixelValue := Pixel{R: 100, G: 150, B: 200}
	ppm.Set(2, 2, newPixelValue)
	fmt.Printf("Nouvelle valeur du pixel à l'indice (2, 2): (%d, %d, %d)\n", newPixelValue.R, newPixelValue.G, newPixelValue.B)

	ppm.Invert()
	ppm.Save("invertPPM.ppm")

	ppm.Flip()
	ppm.Save("flipPPM.ppm")

	ppm.Flop()
	ppm.Save("flopPPM.ppm")

	ppm.SetMaxValue(200)
	ppm.Save("maxValuePPM.ppm")

	ppm.Rotate90CW()
	ppm.Save("90PPM.ppm")
}

func ReadPPM(filename string) (*PPM, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var lines []string

	for _, line := range strings.Split(string(content), "\n") {
		if !strings.HasPrefix(line, "#") {
			lines = append(lines, strings.TrimRight(line, "\r"))
		}
	}

	if strings.HasPrefix(lines[0], "P3") || strings.HasPrefix(lines[0], "P6") {
		magicNumber := lines[0]
		max, err := strconv.Atoi(lines[2])
		if err != nil {
			return nil, err
		}

		elmtInLine := strings.Fields(lines[1])
		width, err := strconv.Atoi(elmtInLine[0])
		if err != nil {
			return nil, err
		}
		height, err := strconv.Atoi(elmtInLine[1])
		if err != nil {
			return nil, err
		}

		ppm := &PPM{
			data:        make([][]Pixel, height),
			width:       width,
			height:      height,
			magicNumber: magicNumber,
			max:         max,
		}

		for i := 0; i < height; i++ {
			pixelLine := strings.Fields(lines[i+3])
			ppm.data[i] = make([]Pixel, width)
			for j := 0; j < width; j++ {
				r, _ := strconv.Atoi(pixelLine[j*3])
				g, _ := strconv.Atoi(pixelLine[j*3+1])
				b, _ := strconv.Atoi(pixelLine[j*3+2])
				ppm.data[i][j] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}

		fmt.Println("Data:")
		for _, row := range ppm.data {
			for _, pixel := range row {
				fmt.Printf("(%d, %d, %d) ", pixel.R, pixel.G, pixel.B)
			}
			fmt.Println()
		}

		fmt.Printf("Width: %d, ", ppm.width)
		fmt.Printf("Height: %d\n", ppm.height)
		fmt.Printf("Magic Number: %s\n", ppm.magicNumber)
		fmt.Printf("Max: %d\n", ppm.max)

		return ppm, nil
	} else {
		return nil, fmt.Errorf("Le format de l'image n'est pas PPM.")
	}
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	for _, row := range ppm.data {
		for _, pixel := range row {
			fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Fprintln(file)
	}

	return nil
}

func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-j-1] = ppm.data[i][ppm.width-j-1], ppm.data[i][j]
		}
	}
}

func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ {
		ppm.data[i], ppm.data[ppm.height-i-1] = ppm.data[ppm.height-i-1], ppm.data[i]
	}
}

func (ppm *PPM) SetMaxValue(maxValue int) {
	if maxValue <= 255 {
		multiplicator := float64(maxValue) / float64(ppm.max)
		ppm.max = maxValue

		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {
				ppm.data[i][j].R = uint8(math.Round(float64(ppm.data[i][j].R) * multiplicator))
				ppm.data[i][j].G = uint8(math.Round(float64(ppm.data[i][j].G) * multiplicator))
				ppm.data[i][j].B = uint8(math.Round(float64(ppm.data[i][j].B) * multiplicator))
			}
		}
	} else {
		fmt.Println("Erreur, le maximum doit être inférieur ou égal à 255.")
	}
}

func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.height)
	for i := range newData {
		newData[i] = make([]Pixel, ppm.width)
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			newData[i][j] = ppm.data[ppm.height-j-1][i]
		}
	}

	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = newData
}
