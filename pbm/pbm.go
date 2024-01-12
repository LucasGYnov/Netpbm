package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool // différentes valeurs possibles (0 ou 1 - noir ou blanc => true ou false)
	width, height int      // taille de l'image (largeur et hauteur)
	magicNumber   string   // Px (de 1 à 4, spécifique au format PBM)
}

func main() {
	var file string
	fmt.Print("Entrer le nom du document (attention à ne pas faire d'erreur) : ")
	fmt.Scanf("%s", &file)
	pbm, err := ReadPBM(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	pbm.Size()

	/* pbm.At(3, 3)
	fmt.Printf("Valeur à l'indice (3, 3): %t\n", pbm.At(3, 3))

	pbm.Set(3, 3, true)
	fmt.Printf("Nouvelle valeur à l'indice (3, 3): %t\n", pbm.At(3, 3))

	pbm.Invert()
	pbm.Save("invertePBM.pbm")
	*/
	pbm.Flip()
	pbm.Save("flippedPBM.pbm")

	pbm.Flop()
	pbm.Save("floppedPBM.pbm")

	/* pbm.SetMagicNumber("P4")
	pbm.Save("changedMagicNumberPBM.pbm") */
}

func ReadPBM(filename string) (*PBM, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var lines []string

	for _, line := range strings.Split(string(content), "\n") {
		if !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if strings.HasPrefix(lines[0], "P1") || strings.HasPrefix(lines[0], "P4") {
		magicNumber := lines[0]

		elmtInLine := strings.Fields(lines[1]) //sépare en différents éléments la string lorsqu'un espace est rencontré
		width, err := strconv.Atoi(elmtInLine[0])
		if err != nil {
			return nil, err
		}
		height, err := strconv.Atoi(elmtInLine[1])
		if err != nil {
			return nil, err
		}

		pbm := &PBM{
			data:        make([][]bool, height),
			width:       width,
			height:      height,
			magicNumber: magicNumber,
		}

		for i := 0; i < height; i++ {
			bitImg := strings.Fields(lines[i+2])
			pbm.data[i] = make([]bool, width)
			for j, val := range bitImg {
				bit, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("Erreur lors de la conversion: %v", err)
				}
				pbm.data[i][j] = bit == 1
			}
		}

		fmt.Println("Data:")
		for _, row := range pbm.data {
			for _, value := range row {
				if value {
					fmt.Print("1 ")
				} else {
					fmt.Print("0 ")
				}
			}
			fmt.Println()
		}

		fmt.Printf("Width: %d, ", pbm.width)
		fmt.Printf("Height: %d\n", pbm.height)
		fmt.Printf("Magic Number: %s\n", pbm.magicNumber)

		return pbm, nil
	} else {
		return nil, fmt.Errorf("Le format de l'image n'est pas PBM.")
	}
}

func (pbm *PBM) Size() (int, int) {
	fmt.Printf("Taille de l'image: %d x %d\n", pbm.width, pbm.height)
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	for _, row := range pbm.data {
		for _, value := range row {
			if value {
				fmt.Fprint(file, "1 ")
			} else {
				fmt.Fprint(file, "0 ")
			}
		}
		fmt.Fprintln(file)
	}

	return nil
}

func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

/* func (pbm *PBM) SetMagicNumber(magicNumber string) error {

} */
