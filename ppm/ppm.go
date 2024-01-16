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

type Point struct {
	X, Y int
}

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

func main() {
	/*var file string
	 fmt.Print("Entrer le nom du document (attention à ne pas faire d'erreur) : ")
	fmt.Scanf("%s", &file) */
	ppm, err := ReadPPM("test.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	/* ppm.Size()
	fmt.Println(ppm.Size()) */

	/* pixelValue := ppm.At(2, 2)
	fmt.Printf("Valeur du pixel à l'indice (2, 2): (%d, %d, %d)\n", pixelValue.R, pixelValue.G, pixelValue.B) */

	/* newPixelValue := Pixel{R: 100, G: 150, B: 200}
	ppm.Set(2, 2, newPixelValue)
	fmt.Printf("Nouvelle valeur du pixel à l'indice (2, 2): (%d, %d, %d)\n", newPixelValue.R, newPixelValue.G, newPixelValue.B)
	ppm.Save("change.ppm") */

	/* ppm.Invert()
	ppm.Save("invertPPM.ppm") */

	/* ppm.Flip()
	ppm.Save("flipPPM.ppm") */

	/* ppm.Flop()
	ppm.Save("flopPPM.ppm") */

	/* ppm.SetMaxValue(10)
	ppm.Save("maxValuePPM.ppm")
	*/
	/* ppm.SetMagicNumber("P6")
	ppm.Save("maxvaluechange.ppm") */

	/* ppm.Rotate90CW()
	ppm.Save("90PPM.ppm") */

	/* ppm.DrawLine(Point{0, 0}, Point{0, 5}, Pixel{R: 255, G: 0, B: 0})
	ppm.Save("drawLine.ppm") */

	ppm.DrawRectangle(Point{2, 2}, 4, 3, Pixel{R: 0, G: 255, B: 0})
	ppm.Save("drawRectangle.ppm")

	ppm.DrawFilledRectangle(Point{8, 2}, 3, 4, Pixel{R: 0, G: 0, B: 255})

	ppm.Save("drawFRectangle.ppm")
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

		/* fmt.Println("Data:")
		for _, row := range ppm.data {
			for _, pixel := range row {
				fmt.Printf("(%d, %d, %d) ", pixel.R, pixel.G, pixel.B)
			}
			fmt.Println()
		} */

		/* fmt.Printf("Width: %d, ", ppm.width)
		fmt.Printf("Height: %d\n", ppm.height)
		fmt.Printf("Magic Number: %s\n", ppm.magicNumber)
		fmt.Printf("Max: %d\n", ppm.max) */

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

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	if maxValue <= 255 || maxValue >= 1 {
		multiplicator := float64(maxValue) / float64(ppm.max)
		ppm.max = int(maxValue)

		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {
				ppm.data[i][j].R = uint8(math.Round(float64(ppm.data[i][j].R) * float64(multiplicator)))
				ppm.data[i][j].G = uint8(math.Round(float64(ppm.data[i][j].G) * float64(multiplicator)))
				ppm.data[i][j].B = uint8(math.Round(float64(ppm.data[i][j].B) * float64(multiplicator)))
			}
		}
	} else {
		fmt.Println("Erreur, le maximum doit être différent de zéro.")
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

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	//...
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{X: p1.X + width - 1, Y: p1.Y}
	p3 := Point{X: p1.X, Y: p1.Y + height - 1}
	p4 := Point{X: p1.X + width - 1, Y: p1.Y + height - 1}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p4, color)
	ppm.DrawLine(p4, p3, color)
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			ppm.Set(p1.X+j, p1.Y+i, color)
		}
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// ...
}

// DrawFilledCircle draws a filled circle in the PPM image.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// ...
}

// DrawTriangle draws a triangle in the PPM image.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// ...
}

// DrawFilledTriangle draws a filled triangle in the PPM image.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// ...
}

// DrawPolygon draws a polygon in the PPM image.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// ...
}

// DrawFilledPolygon draws a filled polygon in the PPM image.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// ...
}
