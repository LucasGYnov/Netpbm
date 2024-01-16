package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func main() {
	/*var file string
	 fmt.Print("Entrer le nom du document (attention à ne pas faire d'erreur) : ")
	fmt.Scanf("%s", &file) */
	pgm, err := ReadPGM("test.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	/* pgm.Size()
	fmt.Println(pgm.Size()) */

	/* fmt.Printf("Valeur à l'indice (4, 4): %d\n", pgm.At(4, 4)) */

	/* pgm.Set(3, 3, 11)
	fmt.Printf("Nouvelle valeur à l'indice (3, 3): %d\n", pgm.At(3, 3))
	pgm.Save("change.pgm") */

	/* pgm.Save("savePGM.pbm") */

	/* pgm.Invert()
	pgm.Save("invertPGM.pbm") */

	/* pgm.Flip()
	pgm.Save("flipPGM.pbm") */

	/* pgm.Flop()
	pgm.Save("flopPGM.pbm") */

	pgm.SetMaxValue(5)
	pgm.Save("maxValuePGM.pbm")

	/* pgm.Rotate90CW()
	pgm.Save("90PGM.pbm") */

}

func ReadPGM(filename string) (*PGM, error) {
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

	if strings.HasPrefix(lines[0], "P2") || strings.HasPrefix(lines[0], "P5") {
		magicNumber := lines[0]
		max, err := strconv.Atoi(lines[2])
		if err != nil {
			return nil, err
		}

		elmtInLine := strings.Fields(lines[1]) //sépare en différents éléments la string lorsqu'un espace est rencontré
		width, err := strconv.Atoi(elmtInLine[0])
		if err != nil {
			return nil, err
		}
		height, err := strconv.Atoi(elmtInLine[1])
		if err != nil {
			return nil, err
		}

		pgm := &PGM{
			data:        make([][]uint8, height),
			width:       width,
			height:      height,
			magicNumber: magicNumber,
			max:         max,
		}

		for i := 0; i < height; i++ {
			bitImg := strings.Fields(lines[i+3])
			pgm.data[i] = make([]uint8, width)
			for j, val := range bitImg {
				bit, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("Erreur lors de la conversion: %v", err)
				}
				pgm.data[i][j] = uint8(bit)
			}
		}

		return pgm, nil
	} else {
		return nil, fmt.Errorf("Le format de l'image n'est pas PGM.")
	}
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	for _, row := range pgm.data {
		for _, value := range row {
			if value >= 100 {
				fmt.Fprint(file, value, " ")
			} else if value <= 10 {
				fmt.Fprint(file, value, "   ")
			} else {
				fmt.Fprint(file, value, "  ")
			}
		}
		fmt.Fprintln(file)
	}

	return nil
}

func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue int) {
	if maxValue >= 1 && maxValue <= 255 { // On vérifie que maxValue soit compris entre 1 et 255 inclus
		pgm.max = int(maxValue)

		for i := 0; i < pgm.height; i++ { // On parcours la matrice
			for j := 0; j < pgm.width; j++ {
				pgm.data[i][j] = uint8(math.Round(float64(pgm.data[i][j]) / float64(pgm.max) * 255)) // Chaque valeur de data est modifié selon le coefficient multiplicateur de maxValue
			}
		}
	}
}

func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			newData[j][pgm.height-i-1] = pgm.data[i][j]
		}
	}

	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = newData
}
