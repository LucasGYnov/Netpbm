package Netpbm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PBM is a struct that represents a PBM image.
type PBM struct {
	data          [][]bool // 2D slice to store the binary image data
	width, height int      // Width and height of the image
	magicNumber   string   // PBM file format identifier ("P1" for ASCII, "P4" for binary)
}

// ReadPBM reads a PBM image from a file and returns a struct representing the image.
func ReadPBM(filename string) (*PBM, error) {
	content, err := os.ReadFile(filename) // Read the entire file content into memory
	if err != nil {                       // Check for file read errors
		return nil, err
	}

	var lines []string // Variable representing the lines of the file

	// Remove comments and create a slice of lines
	for _, line := range strings.Split(string(content), "\n") { // Iterate through lines in the content
		if !strings.HasPrefix(line, "#") { // Exclude lines starting with '#' (comments)
			lines = append(lines, line)
		}
	}

	// Check if there are enough lines to extract necessary information
	if len(lines) < 3 {
		return nil, fmt.Errorf("Invalid PBM file format: insufficient number of lines")
	}

	// Extract magic number, width, and height information
	magicNumber := lines[0]                   // The first line corresponds to the magic number
	elmtInLine := strings.Fields(lines[1])    // Split the elements of the second line
	width, err := strconv.Atoi(elmtInLine[0]) // Retrieve the width from the second line
	if err != nil {                           //error handling
		return nil, err
	}

	height, err := strconv.Atoi(elmtInLine[1]) // Retrieve the height from the second line
	if err != nil {                            //error handling
		return nil, err
	}

	// Initialize PBM struct
	pbm := &PBM{
		data:        make([][]bool, height),
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}

	// Parse image data based on the PBM format
	if strings.HasPrefix(lines[0], "P1") { // P1 format (ASCII)
		for row := 0; row < height; row++ { // Iterate through each row
			bitRow := strings.Fields(lines[row+2]) // Split the bits in the current row
			pbm.data[row] = make([]bool, width)    // Initialize the row in the data slice
			for col, val := range bitRow {         // Iterate through each bit in the row
				// Convert ASCII string to boolean values
				bit, err := strconv.ParseBool(val)
				if err != nil { //error handling
					return nil, fmt.Errorf("Error during conversion: %v", err)
				}
				pbm.data[row][col] = bit
			}
		}
	} else if strings.HasPrefix(lines[0], "P4") { // P4 format (binary)
		binaryStart := len(lines[0]) + len(lines[1]) + 2 // Calculate the starting index of binary data
		binaryData := content[binaryStart:]              // Extract binary data from content
		binaryIndex := 0                                 // Initialize binary data index

		for row := 0; row < height; row++ { // Iterate through each row
			pbm.data[row] = make([]bool, width) // Initialize the row in the data slice
			for col := 0; col < width; col++ {  // Iterate through each bit in the row
				// Extract individual bits from binary data
				byteIndex := binaryIndex / 8                         // Calculate the byte index in binary data
				bitIndex := 7 - binaryIndex%8                        // Calculate the bit index in the byte
				bit := (binaryData[byteIndex] >> uint(bitIndex)) & 1 // Extract the bit from the byte
				pbm.data[row][col] = int(bit) == 1                   // Convert the bit to boolean and store in data
				binaryIndex++                                        // Move to the next bit
			}

			// Ensure byte alignment
			binaryIndex = (binaryIndex + 7) / 8 * 8 // Move to the next byte boundary
		}
	} else {
		// Unsupported PBM format
		return nil, fmt.Errorf("Unsupported PBM format: %s", lines[0])
	}
	return pbm, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height // Return the width and height of the image
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x] // Return the value of the pixel at the specified coordinates
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value // Set the value of the pixel at the specified coordinates
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename) // Create or open the file for writing
	if err != nil {                  // Check for file creation errors
		return err
	}
	defer file.Close() // Close the file when the function completes

	// Write PBM header information
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Write image data based on PBM format
	if pbm.magicNumber == "P1" { // P1 format (ASCII)
		for row := 0; row < pbm.height; row++ { // Iterate through each row in the data
			for col := 0; col < pbm.width; col++ { // Iterate through each value in the row
				// Convert boolean values to ASCII string
				if pbm.data[row][col] {
					fmt.Fprint(file, "1 ")
				} else {
					fmt.Fprint(file, "0 ")
				}
			}
			fmt.Fprintln(file) // Move to the next line after writing a row
		}
	} else if pbm.magicNumber == "P4" { // P4 format (binary)
		for row := 0; row < pbm.height; row++ { // Iterate through each row in the data
			byteRow := make([]byte, (pbm.width+7)/8) // Create a byte slice for the row
			for col := 0; col < pbm.width; col++ {   // Iterate through each value in the row
				if pbm.data[row][col] { // Convert boolean values to binary data
					byteRow[col/8] |= 1 << uint(7-col%8) // Set the corresponding bit in the byte
				}
			}
			file.Write(byteRow) // Write the binary row to the file
		}
	}
	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for row := 0; row < pbm.height; row++ { // Iterate through each row
		for col := 0; col < pbm.width; col++ { // Iterate through each value in the row
			pbm.data[row][col] = !pbm.data[row][col] // Invert the boolean value
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for row := 0; row < pbm.height; row++ { // Iterate through each row
		for col := 0; col < pbm.width/2; col++ { // Iterate through each half of the row
			pbm.data[row][col], pbm.data[row][pbm.width-col-1] = pbm.data[row][pbm.width-col-1], pbm.data[row][col] // Swap pixels horizontally
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for row := 0; row < pbm.height/2; row++ { // Iterate through each half of the rows
		pbm.data[row], pbm.data[pbm.height-row-1] = pbm.data[pbm.height-row-1], pbm.data[row] // Swap rows vertically
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber // Set the magic number of the image
}
