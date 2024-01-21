package Netpbm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM is a struct that represents a PGM image.
type PGM struct {
	data          [][]uint8 // 2D slice to store the image data
	width, height int       // Width and height of the image
	magicNumber   string    // PGM file format identifier ("P2" for ASCII, "P5" for binary)
	max           int       // Maximum pixel value
}

// ReadPGM reads a PGM image from a file and returns a struct representing the image.
func ReadPGM(filename string) (*PGM, error) {
	content, err := os.ReadFile(filename) // Read the entire file content into memory
	if err != nil {                       // Check for file read errors
		return nil, err
	}

	var fileLines []string // Variable representing the lines of the file

	// Remove comments and create a slice of lines
	for _, line := range strings.Split(string(content), "\n") { // Iterate through lines in the content
		if !strings.HasPrefix(line, "#") { // Exclude lines starting with '#' (comments)
			fileLines = append(fileLines, strings.TrimRight(line, "\r"))
		}
	}

	// Check if there are enough lines to extract necessary information
	if len(fileLines) < 3 {
		return nil, fmt.Errorf("Invalid PGM file format: insufficient number of lines")
	}

	// Check if the image format is P2 (ASCII) or P5 (binary)
	if strings.HasPrefix(fileLines[0], "P2") || strings.HasPrefix(fileLines[0], "P5") {
		magicNumber := fileLines[0] // The first line corresponds to the magic number
		max, err := strconv.Atoi(fileLines[2])
		if err != nil { //error handling
			return nil, err
		}

		elmtsInLine := strings.Fields(fileLines[1]) // Split the elements of the second line
		width, err := strconv.Atoi(elmtsInLine[0])
		if err != nil { //error handling
			return nil, err
		}
		height, err := strconv.Atoi(elmtsInLine[1])
		if err != nil { //error handling
			return nil, err
		}

		// Initialize PGM struct
		pgm := &PGM{
			data:        make([][]uint8, height),
			width:       width,
			height:      height,
			magicNumber: magicNumber,
			max:         max,
		}

		if strings.HasPrefix(fileLines[0], "P2") { // Handle P2 ASCII format
			for rowIndex := 0; rowIndex < height; rowIndex++ {
				bitRow := strings.Fields(fileLines[rowIndex+3]) // Split the bits in the current row
				pgm.data[rowIndex] = make([]uint8, width)       // Initialize the row in the data slice
				for colIndex, val := range bitRow {             // Iterate through each bit in the row
					bit, err := strconv.Atoi(val)
					if err != nil { //error handling
						return nil, fmt.Errorf("Error during conversion: %v", err)
					}
					pgm.data[rowIndex][colIndex] = uint8(bit) // Convert and store the bit in data
				}
			}
		} else if strings.HasPrefix(fileLines[0], "P5") { // Handle P5 binary format
			binaryStart := len(fileLines[0]) + len(fileLines[1]) + len(fileLines[2]) + 3 // Calculate the starting index of binary data
			data := content[binaryStart:]                                                // Extract binary data from content

			for rowIndex := 0; rowIndex < height; rowIndex++ { // Iterate through each row
				pgm.data[rowIndex] = make([]uint8, width)         // Initialize the row in the data slice
				for colIndex := 0; colIndex < width; colIndex++ { // Iterate through each bit in the row
					pgm.data[rowIndex][colIndex] = uint8(data[rowIndex*width+colIndex]) // Extract and store the bit in data
				}
			}
		}

		return pgm, nil
	} else {
		return nil, fmt.Errorf("Unsupported PGM format.")
	}
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the pixel value at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the pixel value at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename) // Create or open the file for writing
	if err != nil {                  // Check for file creation errors
		return err
	}
	defer file.Close() // Close the file when the function completes

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	if pgm.magicNumber == "P2" {
		// Handle P2 ASCII format
		for _, row := range pgm.data { // Iterate through each row in the data slice
			for _, value := range row { // Iterate through each value in the row
				if value >= 100 {
					fmt.Fprint(file, value, " ") // Print value followed by a space for values greater than or equal to 100
				} else if value <= 10 {
					fmt.Fprint(file, value, "   ") // Print value followed by three spaces for values less than or equal to 10
				} else {
					fmt.Fprint(file, value, "  ") // Print value followed by two spaces for other values
				}
			}
			fmt.Fprintln(file) // Move to the next line after processing each row
		}
	} else if pgm.magicNumber == "P5" {
		// Handle P5 binary format
		for _, row := range pgm.data { // Iterate through each row in the data slice
			for _, value := range row { // Iterate through each value in the row
				file.Write([]byte{value}) // Write the byte value to the file for binary format
			}
		}
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for rowIndex := 0; rowIndex < pgm.height; rowIndex++ { // Iterate through each row in the image
		for colIndex := 0; colIndex < pgm.width; colIndex++ { // Iterate through each column in the row
			pgm.data[rowIndex][colIndex] = uint8(pgm.max) - pgm.data[rowIndex][colIndex] // Invert the color by subtracting each pixel value from the maximum value
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for rowIndex := 0; rowIndex < pgm.height; rowIndex++ { // Iterate through each row in the image
		for colIndex := 0; colIndex < pgm.width/2; colIndex++ { // Iterate through each column up to the middle of the row
			// Swap pixel values between the current column and its corresponding column on the other side
			pgm.data[rowIndex][colIndex], pgm.data[rowIndex][pgm.width-colIndex-1] = pgm.data[rowIndex][pgm.width-colIndex-1], pgm.data[rowIndex][colIndex]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for rowIndex := 0; rowIndex < pgm.height/2; rowIndex++ { // Iterate through each row up to the middle of the image
		// Swap the entire row with its corresponding row on the other side
		pgm.data[rowIndex], pgm.data[pgm.height-rowIndex-1] = pgm.data[pgm.height-rowIndex-1], pgm.data[rowIndex]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber // Set the magic number of the image
}

// SetMaxValue sets the maximum pixel value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue int) {
	if maxValue >= 1 && maxValue <= 255 { // Check if the specified maximum value is within a valid range
		oldMax := pgm.max
		pgm.max = int(maxValue)

		maxFloat := float64(pgm.max)
		oldMaxFloat := float64(oldMax)
		for rowIndex := 0; rowIndex < pgm.height; rowIndex++ { // Iterate through each row in the image
			for colIndex := 0; colIndex < pgm.width; colIndex++ { // Iterate through each column in the row
				// Rescale each pixel value based on the new and old maximum values
				pgm.data[rowIndex][colIndex] = uint8(float64(pgm.data[rowIndex][colIndex]) * maxFloat / oldMaxFloat)
			}
		}
	}
}

// Rotate90CW rotates the PGM image 90 degrees clockwise.
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width) // Create a new 2D slice with swapped width and height
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
	}

	for rowIndex := 0; rowIndex < pgm.height; rowIndex++ { // Iterate through each row in the original image
		for colIndex := 0; colIndex < pgm.width; colIndex++ { // Iterate through each column in the original image
			// Rotate the pixel values by 90 degrees clockwise
			newData[colIndex][pgm.height-rowIndex-1] = pgm.data[rowIndex][colIndex]
		}
	}

	pgm.width, pgm.height = pgm.height, pgm.width // Swap the width and height of the image
	pgm.data = newData                            // Set the image data to the rotated data
}

// ToPBM converts a PGM image to a PBM image.
func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		data:        make([][]bool, pgm.height), // Create a new PBM image with the same dimensions
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	for rowIndex := 0; rowIndex < pgm.height; rowIndex++ { // Iterate through each row in the PGM image
		pbm.data[rowIndex] = make([]bool, pgm.width)
		for colIndex := 0; colIndex < pgm.width; colIndex++ { // Iterate through each column in the PGM image
			// Convert each pixel value to a boolean value based on a threshold
			pbm.data[rowIndex][colIndex] = pgm.data[rowIndex][colIndex] < uint8(pgm.max/2)
		}
	}
	return pbm // Return the resulting PBM image
}
