package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	inputDir := flag.String("input", ".", "Directory containing the input images")
	outputFile := flag.String("output", "output.png", "Output image file")
	maxWidth := flag.Int("maxwidth", 820, "Maximum width of the output image")
	thumbnailHeight := flag.Int("height", 200, "Height of the thumbnails")
	flag.Parse()

	thumbnailsDir := "thumbnails"

	// Check if the thumbnails directory exists and prompt for removal if it does
	if _, err := os.Stat(thumbnailsDir); !os.IsNotExist(err) {
		fmt.Printf("The directory %s already exists. Do you want to remove it and proceed? (y/n): ", thumbnailsDir)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Operation cancelled.")
			return
		}
		os.RemoveAll(thumbnailsDir)
	}

	// Create thumbnails directory
	os.Mkdir(thumbnailsDir, 0755)

	// Resize images and save thumbnails
	files, err := filepath.Glob(filepath.Join(*inputDir, "*.jpg"))
	if err != nil {
		log.Fatal(err)
	}

	files = append(files, findJPEGFiles(*inputDir)...)

	thumbnails := []image.Image{}
	for _, file := range files {
		if filepath.Base(file) == filepath.Base(*outputFile) {
			continue // Ignore the output file
		}

		img, err := loadImage(file)
		if err != nil {
			log.Printf("failed to load image %s: %v", file, err)
			continue
		}
		thumbnail := resize.Resize(0, uint(*thumbnailHeight), img, resize.Lanczos3) // Resize to fixed height
		saveThumbnail(filepath.Join(thumbnailsDir, filepath.Base(file)), thumbnail)
		thumbnails = append(thumbnails, thumbnail)
	}

	// Sort thumbnails by width in descending order
	sort.Slice(thumbnails, func(i, j int) bool {
		return thumbnails[i].Bounds().Dx() > thumbnails[j].Bounds().Dx()
	})

	// Create the masonry layout
	createMasonryLayout(thumbnails, *outputFile, *maxWidth)
}

func findJPEGFiles(inputDir string) []string {
	files, err := filepath.Glob(filepath.Join(inputDir, "*.jpeg"))
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func saveThumbnail(filename string, img image.Image) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("failed to save thumbnail %s: %v", filename, err)
		return
	}
	defer file.Close()

	jpeg.Encode(file, img, nil)
}

func createMasonryLayout(images []image.Image, outputFilename string, maxWidth int) {
	var rows [][]image.Image
	var currentRow []image.Image
	currentRowWidth := 0
	totalWidth := maxWidth

	for _, img := range images {
		if currentRowWidth+img.Bounds().Dx() > totalWidth {
			if len(rows) == 0 {
				totalWidth = currentRowWidth // Set the totalWidth based on the first row
			}
			rows = append(rows, currentRow)
			currentRow = []image.Image{}
			currentRowWidth = 0
		}
		currentRow = append(currentRow, img)
		currentRowWidth += img.Bounds().Dx()
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
		if len(rows) == 1 {
			totalWidth = currentRowWidth // Adjust totalWidth if only one row
		}
	}

	totalHeight := 0
	for _, row := range rows {
		maxHeight := 0
		for _, img := range row {
			if img.Bounds().Dy() > maxHeight {
				maxHeight = img.Bounds().Dy()
			}
		}
		totalHeight += maxHeight
	}

	outputImage := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	yOffset := 0
	for _, row := range rows {
		xOffset := 0
		maxHeight := 0
		for _, img := range row {
			draw.Draw(outputImage, img.Bounds().Add(image.Pt(xOffset, yOffset)), img, image.Point{}, draw.Over)
			xOffset += img.Bounds().Dx()
			if img.Bounds().Dy() > maxHeight {
				maxHeight = img.Bounds().Dy()
			}
		}
		yOffset += maxHeight
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	png.Encode(outputFile, outputImage)
	fmt.Printf("Masonry layout created and saved to %s\n", outputFilename)
}

func max(vals ...int) int {
	m := vals[0]
	for _, v := range vals {
		if v > m {
			m = v
		}
	}
	return m
}

func sum(vals []int) int {
	total := 0
	for _, v := range vals {
		total += v
	}
	return total
}
