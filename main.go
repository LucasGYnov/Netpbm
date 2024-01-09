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

	ReadPBM(file)
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
