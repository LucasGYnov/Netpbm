package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool // position en x et y, (bool => 0 ou 1 - blanc ou noir)
	width, height int      // taille de l'image
	magicNumber   string   //pX
}

func main() {
	ReadPBM("test.pbm")
}

func ReadPBM(filename string) /*  (*PBM, error) */ {
	image, err := lireFichier(filename)
	if err != nil {
		log.Fatal(err)
	}

	line := strings.Split(image, "\n")

	var allLine []string

	for _, elmt := range line {
		entireImage := fmt.Sprintf("%s\n", elmt)
		allLine = append(allLine, entireImage)
	}

	/* for i := 0; i < len(allLine); i++ {
		fmt.Println(line[i])
	} */
	/* fmt.Println("Le magic number est : ", line[0], ".L'image est de largeur", line[2][0:2], "et d'hauteur", line[2][3:5], ".") */
	magicNumber := line[0][0:2]
	largeur := line[2][0:2]
	hauteur := line[2][3:5]

	fmt.Println(magicNumber, largeur, hauteur)
}

func lireFichier(fichier string) (string, error) {
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		return "", err
	}
	return string(contenu), nil
}
