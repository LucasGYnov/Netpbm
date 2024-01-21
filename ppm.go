package Netpbm

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// PPM is a struct representing a Portable Pixmap image.
type PPM struct {
	data          [][]Pixel // Pixel data for each position in the image
	width, height int       // Width and height of the image
	magicNumber   string    // Format identifier ("P3"for ASCII, "P6" for binary)
	max           int       // Maximum color value in the image
}

// Pixel represents a color with red (R), green (G), and blue (B) channels.
type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a PPM struct.
func ReadPPM(filename string) (*PPM, error) {
	content, err := os.ReadFile(filename) // Read the entire file content into memory
	if err != nil {                       // Check for file read errors
		return nil, err
	}

	var fileLines []string // Variable representing the lines of the file

	// Read lines from the content, excluding comments
	for _, line := range strings.Split(string(content), "\n") { // Iterate through lines in the content
		if !strings.HasPrefix(line, "#") { // Exclude lines starting with '#' (comments)
			fileLines = append(fileLines, strings.TrimRight(line, "\r"))
		}
	}

	// Check if there are enough lines to extract necessary information
	if len(fileLines) < 3 {
		return nil, fmt.Errorf("Invalid PGM file format: insufficient number of lines")
	}

	// Check if the image format is P3 (ASCII) or P6 (binary)
	if strings.HasPrefix(fileLines[0], "P3") || strings.HasPrefix(fileLines[0], "P6") {
		magicNumber := fileLines[0] // The first line corresponds to the magic number
		max, err := strconv.Atoi(fileLines[2])
		if err != nil { //error handling
			return nil, err
		}

		elmtInLine := strings.Fields(fileLines[1]) // Split the elements of the second line
		width, err := strconv.Atoi(elmtInLine[0])
		if err != nil { //error handling
			return nil, err
		}
		height, err := strconv.Atoi(elmtInLine[1])
		if err != nil { //error handling
			return nil, err
		}

		// Initialize PPM struct
		ppm := &PPM{
			data:        make([][]Pixel, height),
			width:       width,
			height:      height,
			magicNumber: magicNumber,
			max:         max,
		}

		if magicNumber == "P3" { // Parse ASCII data for P3 format
			for y := 0; y < height; y++ { // Extract the pixel values from the current line
				pixelLine := strings.Fields(fileLines[y+3])

				ppm.data[y] = make([]Pixel, width) // Initialize a row in the PPM data for the current line

				for x := 0; x < width; x++ { // Iterate through each pixel in the line

					// Convert the ASCII values to integers for each channel
					red, _ := strconv.Atoi(pixelLine[x*3])
					green, _ := strconv.Atoi(pixelLine[x*3+1])
					blue, _ := strconv.Atoi(pixelLine[x*3+2])

					// Assign the RGB values to the current pixel in the PPM data
					ppm.data[y][x] = Pixel{R: uint8(red), G: uint8(green), B: uint8(blue)}
				}
			}
		} else if magicNumber == "P6" { // Parse binary data for P6 format

			binaryStart := len(fileLines[0]) + len(fileLines[1]) + len(fileLines[2]) + 3 // Calculate the starting position of binary data in the content

			data := content[binaryStart:] // Extract binary data from the content

			for y := 0; y < height; y++ { // Initialize a row in the PPM data for the current line
				ppm.data[y] = make([]Pixel, width)
				for x := 0; x < width; x++ { // Iterate through each pixel in the line

					startPos := (y*width + x) * 3 // Calculate the starting position of the current pixel in the binary data

					// Assign the RGB values to the current pixel in the PPM data
					ppm.data[y][x].R = uint8(data[startPos])
					ppm.data[y][x].G = uint8(data[startPos+1])
					ppm.data[y][x].B = uint8(data[startPos+2])
				}
			}
		}

		return ppm, nil
	} else {
		return nil, fmt.Errorf("Image format is not PPM.")
	}
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the color (Pixel) at a specified position (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set updates the color (Pixel) at a specified position (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save saves the PPM image to a file with the specified filename.
func (ppm *PPM) Save(filename string) error {

	// Open the file in write mode, create it if it doesn't exist
	file, err := os.Create(filename)
	if err != nil { // Check for file creation errors
		return err
	}
	defer file.Close() // Close the file when the function exits

	// Write header information to the file
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	if ppm.magicNumber == "P3" {
		// Save in ASCII P3 format
		for _, row := range ppm.data {
			for _, pixel := range row {
				fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B) // Write RGB color values separated by spaces
			}
			fmt.Fprintln(file) // New line after each row of pixels
		}
	} else if ppm.magicNumber == "P6" { // Save in binary P6 format
		for _, row := range ppm.data {
			for _, pixel := range row {
				file.Write([]byte{pixel.R, pixel.G, pixel.B}) // Write RGB color values directly in binary
			}
		}
	} else {
		return fmt.Errorf("Unsupported PPM format: %s", ppm.magicNumber) // Return an error if the PPM format is not supported
	}

	return nil
}

// Invert inverts the colors of the PPM image.
// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for row := 0; row < ppm.height; row++ {
		for col := 0; col < ppm.width; col++ {
			// Invert the color of each RGB channel by subtracting it from the maximum value.
			ppm.data[row][col].R = uint8(ppm.max) - ppm.data[row][col].R
			ppm.data[row][col].G = uint8(ppm.max) - ppm.data[row][col].G
			ppm.data[row][col].B = uint8(ppm.max) - ppm.data[row][col].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for row := 0; row < ppm.height; row++ {
		for col := 0; col < ppm.width/2; col++ {
			ppm.data[row][col], ppm.data[row][ppm.width-col-1] = ppm.data[row][ppm.width-col-1], ppm.data[row][col] // Swap the pixel values between the current column and its corresponding column on the other side.
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for row := 0; row < ppm.height/2; row++ {
		ppm.data[row], ppm.data[ppm.height-row-1] = ppm.data[ppm.height-row-1], ppm.data[row] // Swap the entire row with its corresponding row on the other side.
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the maximum color value of the PPM image and scales the pixel values accordingly to fit within the new maximum.
func (ppm *PPM) SetMaxValue(newMaxValue uint8) {
	if newMaxValue >= 1 && newMaxValue <= 255 { // Ensure that the new maximum value is within the valid range (1 to 255).

		// Calculate the scaling factor to adjust pixel values based on the new maximum.
		scalingFactor := float64(newMaxValue) / float64(ppm.max)
		ppm.max = int(newMaxValue)

		for row := 0; row < ppm.height; row++ { // Apply scaling to each pixel in the image.
			for col := 0; col < ppm.width; col++ {

				// Round the result of the scaling operation and update each color channel.
				ppm.data[row][col].R = uint8(math.Round(float64(ppm.data[row][col].R) * scalingFactor))
				ppm.data[row][col].G = uint8(math.Round(float64(ppm.data[row][col].G) * scalingFactor))
				ppm.data[row][col].B = uint8(math.Round(float64(ppm.data[row][col].B) * scalingFactor))
			}
		}
	} else {
		fmt.Println("Error: The maximum must be between 1 and 255.")
	}
}

// Rotate90CW rotates the PPM image 90 degrees clockwise.
func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.height)
	for i := range newData {
		newData[i] = make([]Pixel, ppm.width)
	}

	// Rotate pixel values by 90 degrees clockwise.
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			newData[i][j] = ppm.data[ppm.height-j-1][i]
		}

	}

	// Update width, height, and data with the rotated values.
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = newData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Create a new PGM instance with the same dimensions and maximum value as the original PPM
	pgm := &PGM{
		data:        make([][]uint8, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
	}

	// Iterate through each pixel in the PPM image
	for y := 0; y < ppm.height; y++ {
		// Initialize a new row in the PGM data
		pgm.data[y] = make([]uint8, ppm.width)

		for x := 0; x < ppm.width; x++ {
			// Convert RGB to grayscale using a simple average method
			// Calculate the average of the red, green, and blue values for each pixel
			gray := (int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3

			// Store the calculated grayscale value in the corresponding position in the PGM data
			pgm.data[y][x] = uint8(gray)
		}
	}

	// Return the resulting PGM instance
	return pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Create a new PBM instance with the same dimensions and "P1" magic number
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	// Initialize the PBM data with a 2D boolean array
	pbm.data = make([][]bool, ppm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
	}

	// Calculate a threshold value for converting RGB to binary (black or white)
	threshold := uint8(ppm.max / 2)

	// Iterate through each pixel in the PPM image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Calculate the average of the red, green, and blue values for each pixel
			average := (uint16(ppm.data[y][x].R) + uint16(ppm.data[y][x].G) + uint16(ppm.data[y][x].B)) / 3

			// Set the corresponding position in the PBM data to true if average is below the threshold
			pbm.data[y][x] = average < uint16(threshold)
		}
	}

	// Return the resulting PBM instance
	return pbm
}

// Point represents a 2D point with X and Y coordinates.
type Point struct {
	X, Y int
}

// DrawLine draws a line between two points in the PPM image.
func (ppm *PPM) DrawLine(startPoint, endPoint Point, lineColor Pixel) {
	// Calculate the differences in x and y coordinates.
	deltaX := endPoint.X - startPoint.X
	deltaY := endPoint.Y - startPoint.Y

	// Initialize the starting point (x, y) and the increments for each step.
	x, y := startPoint.X, startPoint.Y
	xIncrement, yIncrement := 1, 1

	// Adjust increments and deltas based on the signs of deltaX and deltaY.
	if deltaX < 0 {
		xIncrement = -1
		deltaX = -deltaX
	}
	if deltaY < 0 {
		yIncrement = -1
		deltaY = -deltaY
	}

	// Check whether the slope is less than 1 or greater than/equal to 1.
	if deltaX > deltaY {
		// Slope is less than 1
		decisionParameter := 2*deltaY - deltaX
		for i := 0; i <= deltaX; i++ {
			// Set the color of the pixel at the current (x, y) position.
			if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
				ppm.Set(x, y, lineColor)
			}

			// Update y and the decision parameter based on Bresenham's algorithm.
			if decisionParameter > 0 {
				y += yIncrement
				decisionParameter -= 2 * deltaX
			}
			x += xIncrement
			decisionParameter += 2 * deltaY
		}
	} else {
		// Slope is greater than or equal to 1
		decisionParameter := 2*deltaX - deltaY
		for i := 0; i <= deltaY; i++ {
			// Set the color of the pixel at the current (x, y) position.
			if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
				ppm.Set(x, y, lineColor)
			}

			// Update x and the decision parameter based on Bresenham's algorithm.
			if decisionParameter > 0 {
				x += xIncrement
				decisionParameter -= 2 * deltaY
			}
			y += yIncrement
			decisionParameter += 2 * deltaX
		}
	}
}

// DrawRectangle draws a rectangle in the PPM image.
func (ppm *PPM) DrawRectangle(topLeft Point, width, height int, rectangleColor Pixel) {
	// Ensure that topLeft is within bounds
	topLeft.X = max(topLeft.X, 0)
	topLeft.Y = max(topLeft.Y, 0)

	// Adjust width and height to fit within bounds
	width = min(width, ppm.width-topLeft.X)
	height = min(height, ppm.height-topLeft.Y)

	// Calculate other three corner points of the rectangle
	topRight := Point{topLeft.X + width, topLeft.Y}
	bottomRight := Point{topLeft.X + width, topLeft.Y + height}
	bottomLeft := Point{topLeft.X, topLeft.Y + height}

	// Draw lines connecting the four corner points to form a rectangle
	ppm.DrawLine(topLeft, topRight, rectangleColor)
	ppm.DrawLine(topRight, bottomRight, rectangleColor)
	ppm.DrawLine(bottomRight, bottomLeft, rectangleColor)
	ppm.DrawLine(bottomLeft, topLeft, rectangleColor)
}

// DrawFilledRectangle draws a filled rectangle in the PPM image.
func (ppm *PPM) DrawFilledRectangle(topLeft Point, width, height int, fillPixel Pixel) {
	// Draw the outline of the rectangle using DrawRectangle method
	ppm.DrawRectangle(topLeft, width, height, fillPixel)

	// Iterate through each row of the PPM image
	for row := 0; row < ppm.height; row++ {
		var coloredPixelPositions []int
		var coloredPixelCount int

		// Iterate through each column of the current row
		for col := 0; col < ppm.width; col++ {
			// Check if the current pixel has the specified fill color
			if ppm.data[row][col] == fillPixel {
				coloredPixelCount++
				coloredPixelPositions = append(coloredPixelPositions, col)
			}
		}

		// Fill the space between the leftmost and rightmost colored pixels in the current row
		if coloredPixelCount > 1 {
			leftmostCol := coloredPixelPositions[0] + 1
			rightmostCol := coloredPixelPositions[len(coloredPixelPositions)-1]

			for col := leftmostCol; col < rightmostCol; col++ {
				ppm.data[row][col] = fillPixel
			}
		}

		// If the specified width or height is greater than the image dimensions, fill the entire row
		if height > ppm.height && width > ppm.width {
			for col := 0; col < ppm.width; col++ {
				ppm.data[row][col] = fillPixel
			}
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// Iterate through each row of the image
	for row := 0; row < ppm.height; row++ {
		// Iterate through each column of the image
		for col := 0; col < ppm.width; col++ {
			// Calculate the horizontal distance from the current pixel to the center of the circle
			dx := float64(col) - float64(center.X)
			// Calculate the vertical distance from the current pixel to the center of the circle
			dy := float64(row) - float64(center.Y)
			// Calculate the distance from the current pixel to the center of the circle using the Pythagorean theorem
			distance := math.Sqrt(dx*dx + dy*dy)

			// Check if the current pixel is approximately on the circumference of the circle
			if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
				// Set the color of the current pixel to the specified color
				ppm.Set(col, row, color)
			}
		}
	}

	// Mark key points on the circle boundary by setting their colors
	ppm.Set(center.X-(radius-1), center.Y, color)
	ppm.Set(center.X+(radius-1), center.Y, color)
	ppm.Set(center.X, center.Y+(radius-1), color)
	ppm.Set(center.X, center.Y-(radius-1), color)
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Draw the circle outline
	ppm.DrawCircle(center, radius, color)

	for row := 0; row < ppm.height; row++ {
		var coloredPixelPositions []int
		var coloredPixelCount int

		// Iterate through each column of the current row
		for col := 0; col < ppm.width; col++ {
			// Check if the current pixel has the specified fill color
			if ppm.data[row][col] == color {
				coloredPixelCount++
				coloredPixelPositions = append(coloredPixelPositions, col)
			}
		}

		// Fill the space between the leftmost and rightmost colored pixels in the current row
		if coloredPixelCount > 1 {
			leftmostCol := coloredPixelPositions[0] + 1
			rightmostCol := coloredPixelPositions[len(coloredPixelPositions)-1]

			for col := leftmostCol; col < rightmostCol; col++ {
				ppm.data[row][col] = color
			}
		}
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Draw the triangle outline
	ppm.DrawTriangle(p1, p2, p3, color)

	for row := 0; row < ppm.height; row++ {
		var coloredPixelPositions []int
		var coloredPixelCount int

		// Iterate through each column of the current row
		for col := 0; col < ppm.width; col++ {
			// Check if the current pixel has the specified fill color
			if ppm.data[row][col] == color {
				coloredPixelCount++
				coloredPixelPositions = append(coloredPixelPositions, col)
			}
		}

		// Fill the space between the leftmost and rightmost colored pixels in the current row
		if coloredPixelCount > 1 {
			leftmostCol := coloredPixelPositions[0] + 1
			rightmostCol := coloredPixelPositions[len(coloredPixelPositions)-1]

			for col := leftmostCol; col < rightmostCol; col++ {
				ppm.data[row][col] = color
			}
		}
	}
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// Iterate through each point in the polygon except the last one
	for currentPointIndex := 0; currentPointIndex < len(points)-1; currentPointIndex++ {
		// Draw a line between the current point and the next point in the polygon
		ppm.DrawLine(points[currentPointIndex], points[currentPointIndex+1], color)
	}

	// Draw a closing line connecting the last point to the first point in the polygon
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Draw the polygon outline
	ppm.DrawPolygon(points, color)

	for row := 0; row < ppm.height; row++ {
		var coloredPixelPositions []int
		var coloredPixelCount int

		// Iterate through each column of the current row
		for col := 0; col < ppm.width; col++ {
			// Check if the current pixel has the specified fill color
			if ppm.data[row][col] == color {
				coloredPixelCount++
				coloredPixelPositions = append(coloredPixelPositions, col)
			}
		}

		// Fill the space between the leftmost and rightmost colored pixels in the current row
		if coloredPixelCount > 1 {
			leftmostCol := coloredPixelPositions[0] + 1
			rightmostCol := coloredPixelPositions[len(coloredPixelPositions)-1]

			for col := leftmostCol; col < rightmostCol; col++ {
				ppm.data[row][col] = color
			}
		}
	}
}

// DrawSierpinskiTriangle draws a Sierpinski triangle.

/*
	The function uses recursive "divide and conquer" approach to create smaller triangles.

Each recursive call decreases the recursion depth (n) and draws a triangle from three sub-triangles. The base case is when n reaches zero.
*** Flip the image vertically with the function Flop() before saving to have the correct orientation. ***
*/
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	// Base case: stop drawing if n is equal to 0
	if n == 0 {
		return
	}

	// Calculate the three points of the equilateral triangle
	p1 := start
	p2 := Point{X: start.X + width, Y: start.Y}
	height := int(math.Sqrt(3.0) * float64(width) / 2)
	p3 := Point{X: start.X + width/2, Y: start.Y + height}

	// Draw the triangle
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)

	// Calculate midpoints of each side
	mid1 := Point{X: (p1.X + p2.X) / 2, Y: (p1.Y + p2.Y) / 2}
	mid3 := Point{X: (p3.X + p1.X) / 2, Y: (p3.Y + p1.Y) / 2}

	// Draw triangles recursively
	ppm.DrawSierpinskiTriangle(n-1, p1, width/2, color)
	ppm.DrawSierpinskiTriangle(n-1, mid1, width/2, color)
	ppm.DrawSierpinskiTriangle(n-1, mid3, width/2, color)
}
