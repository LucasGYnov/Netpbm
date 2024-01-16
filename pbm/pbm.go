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
	/*var file string
	 fmt.Print("Entrer le nom du document (attention à ne pas faire d'erreur) : ")
	fmt.Scanf("%s", &file) */
	pbm, err := ReadPBM("1t.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	/* pbm.Size()
	fmt.Println(pbm.Size()) */

	/* pbm.At(3, 3)
	fmt.Printf("Valeur à l'indice (3, 3): %t\n", pbm.At(3, 3)) */

	/* pbm.Set(3, 3, true)
	fmt.Printf("Nouvelle valeur à l'indice (3, 3): %t\n", pbm.At(3, 3))
	pbm.Save("change.pbm") */

	/* pbm.Invert()
	pbm.Save("invertePBM.pbm") */

	/* pbm.Flip()
	pbm.Save("flippedPBM.pbm") */

	/* pbm.Flop()
	pbm.Save("floppedPBM.pbm") */

	pbm.SetMagicNumber("P4")
	pbm.Save("p4top1.pbm")
}

// ReadPBM reads a PBM file and returns a PBM struct.
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

	if len(lines) < 3 {
		return nil, fmt.Errorf("Invalid PBM file format: insufficient number of lines")
	}

	magicNumber := lines[0]

	elmtInLine := strings.Fields(lines[1])
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

	if strings.HasPrefix(lines[0], "P1") {
		// Handle P1 ASCII format
		for i := 0; i < height; i++ {
			bitImg := strings.Fields(lines[i+2])
			pbm.data[i] = make([]bool, width)
			for j, val := range bitImg {
				bit, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("Error during conversion: %v", err)
				}
				pbm.data[i][j] = bit == 1
			}
		}
	} else if strings.HasPrefix(lines[0], "P4") {
		binaryStart := len(lines[0]) + len(lines[1]) + 2 // 2 accounts for the newline characters
		data := content[binaryStart:]
		binaryIndex := 0

		for i := 0; i < height; i++ {
			pbm.data[i] = make([]bool, width)
			for j := 0; j < width; j++ {
				byteIndex := binaryIndex / 8
				bitIndex := 7 - binaryIndex%8
				bit := (data[byteIndex] >> uint(bitIndex)) & 1
				pbm.data[i][j] = int(bit) == 1
				binaryIndex++
			}
		}
	} else {
		return nil, fmt.Errorf("Unsupported PBM format: %s", lines[0])
	}

	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
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

	/* if pbm.magicNumber == "P1" { */
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
	/* } else if pbm.magicNumber == "P4" {
		for _, row := range pbm.data {
			byteRow := make([]byte, (pbm.width+7)/8)
			for i, value := range row {
				if value {
					byteRow[i/8] |= 1 << uint(7-i%8)
				}
			}
			file.Write(byteRow)
		}
	} */
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

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
