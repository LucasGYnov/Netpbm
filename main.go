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

	for i := 0; i < len(allLine); i++ {
		fmt.Println(line[i])
	}
}

func lireFichier(fichier string) (string, error) {
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		return "", err
	}
	return string(contenu), nil
}

/* file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture")
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	file.Close()
} */

/* image, err := lireFichier(filename)
if err != nil {
	log.Fatal(err)
}
magicNumber := image[0:2]
width := image[2:6]
height := image[6:10]
fmt.Println(image)
fmt.Println(magicNumber)
fmt.Println("L'image à une largeur de ", width, " et une hauteur de ", height)
fmt.Println(len(image)) */
/* if magicNumber[0] != 'P' || magicNumber[1] != 'X' {
	return nil, fmt.Errorf("%s n'est pas un fichier PBM", filename)
} */

/* data := [][]bool{}
var width int
var height int

fmt.Println("L'image à une largeur de ", width, " et une hauteur de ", height)
fmt.Println("Il y a deux types : ", data[0], "et", data[1]) */

/* Lire le fichier, la 3eme ligne donnée => w et h séparé par un espace " ",
data[0] = premier élément de 4eme ligne et data[1] = premier élément différent de la 4eme ligne et différent d'un espace */

/* }
func lireFichier(fichier string) (string, error) {
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		return "", err
	}
	return string(contenu), nil
} */

func main() {
	ReadPBM("test.pbm")
}
