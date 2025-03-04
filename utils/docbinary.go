package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func MixedAttachmentsToPdf(attachments []string) ([]byte, int, error) {
	tmpDir, err := os.MkdirTemp("", "mixed-attachments")
	if err != nil {
		return nil, 0, err
	}
	defer os.RemoveAll(tmpDir)

	pdfAttachments := make([]string, 0)
	for _, attachment := range attachments {
		ext := strings.ToLower(filepath.Ext(attachment))
		if ext == ".pdf" {
			pdfAttachments = append(pdfAttachments, attachment)
		} else if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			pdfFile := filepath.Join(tmpDir, strings.TrimSuffix(filepath.Base(attachment), ext)+".pdf")
			cmd := exec.Command("convert", attachment,
				"-resize", "595x842>", // Resize to fit within A4, preserving aspect ratio
				"-gravity", "center", // Center the image on the A4 page
				"-background", "white", // Use white background for padding
				"-extent", "595x842", // Ensure final output is A4 with padding if needed
				"pdf:"+pdfFile) // Output as PDF

			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Println("ImageMagick convert error:", stderr.String())
				continue
			}
			pdfAttachments = append(pdfAttachments, pdfFile)
		} else {
			log.Printf("Unsupported file type: %s\n", attachment)
		}
	}

	outputPDF := filepath.Join(tmpDir, "output.pdf")
	cmd := exec.Command("gs", append([]string{
		"-dBATCH",
		"-dNOPAUSE",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.6",
		"-sPAPERSIZE=a4",
		"-dFIXEDMEDIA",
		"-dPDFFitPage",
		"-sOUTPUTFILE=" + outputPDF,
	}, pdfAttachments...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Println("Ghostscript error:", stderr.String())
		return nil, 0, err
	}

	pdfData, err := os.ReadFile(outputPDF)
	if err != nil {
		return nil, 0, err
	}

	// Count the number of pages in the PDF
	pageCount, err := countPdfPages(outputPDF)
	if err != nil {
		return nil, 0, err
	}

	return pdfData, pageCount, nil
}

// countPdfPages counts the number of pages in a PDF byte slice using Ghostscript.
// func countPdfPages(pdfFilePath string) (int, error) {

// 	gsCommand := "gs"
// 	gsArgs := []string{
// 		"-q",
// 		"-dNODISPLAY",
// 		"-c",
// 		fmt.Sprintf("(%s) (r) file runpdfbegin pdfpagecount = quit", pdfFilePath),
// 	}

// 	cmd := exec.Command(gsCommand, gsArgs...)
// 	var out bytes.Buffer
// 	cmd.Stdout = &out

// 	if err := cmd.Run(); err != nil {
// 		return 0, fmt.Errorf("failed to count pages with Ghostscript: %w", err)
// 	}

// 	var pageCount int
// 	_, err := fmt.Sscanf(out.String(), "%d", &pageCount)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to parse page count: %w", err)
// 	}

// 	return pageCount, nil
// }

// Count pages using pdfinfo light weight binary
func countPdfPages(pdfFilePath string) (int, error) {
	cmd := exec.Command("pdfinfo", pdfFilePath) // Use pdfinfo to get metadata
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to get PDF page count: %w", err)
	}

	// Parse the output for the "Pages" field
	scanner := bufio.NewScanner(&output)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Pages:") {
			pageCountStr := strings.TrimSpace(strings.TrimPrefix(line, "Pages:"))
			pageCount, err := strconv.Atoi(pageCountStr)
			if err != nil {
				return 0, fmt.Errorf("invalid page count format: %w", err)
			}
			return pageCount, nil
		}
	}

	return 0, fmt.Errorf("page count not found in pdfinfo output")

}

// Merge HTML and pdfs in single file. If no HTML file, even used as a pdf merger
func MergeHtmlAndPdfs(files []interface{}, headerText string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "mixed-files")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	processedPdfs := []string{}

	for i, file := range files {
		switch v := file.(type) {
		case string: // File path
			ext := strings.ToLower(filepath.Ext(v))
			if ext == ".html" {
				// Convert HTML to PDF
				outputPdf := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pdf", i))
				cmd := exec.Command("./xbin/wkhtmltopdf-amd64",
					"--margin-top", "6mm",
					"--margin-bottom", "6mm",
					"--margin-left", "1mm",
					"--margin-right", "1mm",
					"--header-right", headerText,
					"--header-font-size", "7",
					//"--image-quality", "150",
					"--footer-center", "[page]", // Add page number to the right side of the footer
					"--footer-font-size", "6", // Set font size for footer text

					v, outputPdf)

				// fmt.Printf("Running command: %v\n", cmd.Args)
				var stderr bytes.Buffer
				cmd.Stderr = &stderr
				err := cmd.Run()
				if err != nil {
					log.Printf("Failed to convert HTML to PDF for %s: %s\n", file, stderr.String())
					return nil, fmt.Errorf("failed to convert HTML file '%s' to PDF: %w", file, err)
				}

				processedPdfs = append(processedPdfs, outputPdf)
			} else if ext == ".pdf" {
				// Add PDF directly
				processedPdfs = append(processedPdfs, v)
			} else {
				// Treat as HTML string
				htmlPath := filepath.Join(tmpDir, fmt.Sprintf("file_%d.html", i))
				err := os.WriteFile(htmlPath, []byte(v), 0644)
				if err != nil {
					return nil, fmt.Errorf("failed to write HTML string to file: %w", err)
				}

				outputPdf := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pdf", i))
				cmd := exec.Command("./xbin/wkhtmltopdf-amd64",
					"--margin-top", "6mm",
					"--margin-bottom", "6mm",
					"--margin-left", "1mm",
					"--margin-right", "1mm",
					"--header-right", headerText,
					"--header-font-size", "7",
					"--footer-center", "[page]",
					"--footer-font-size", "6",
					htmlPath, outputPdf)

				var stderr bytes.Buffer
				cmd.Stderr = &stderr
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
					log.Printf("Failed to convert HTML string to PDF or invalid file: %s\n", stderr.String())
					return nil, fmt.Errorf("failed to convert HTML string to PDF or invalid file: %w", err)
				}

				processedPdfs = append(processedPdfs, outputPdf)

				// log.Printf("Unsupported file type: %s\n", v)
				// return nil, fmt.Errorf("unsupported file type: %s", v)
			}
		case []byte: // Raw PDF data
			// Save raw PDF data to a temporary file
			rawPdfPath := filepath.Join(tmpDir, fmt.Sprintf("raw_file_%d.pdf", i))
			err := os.WriteFile(rawPdfPath, v, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to write raw PDF data to file: %w", err)
			}
			processedPdfs = append(processedPdfs, rawPdfPath)

		default:
			return nil, fmt.Errorf("unsupported input type at index %d", i)
		}
	}

	// Merge all processed PDFs using ghostscript
	// outputPdf := filepath.Join(tmpDir, "merged_output.pdf")
	// cmd := exec.Command("gs", append([]string{
	// 	"-dBATCH",
	// 	"-dNOPAUSE",
	// 	"-sDEVICE=pdfwrite",
	// 	"-dCompatibilityLevel=1.6",
	// 	"-sPAPERSIZE=a4",
	// 	"-dFIXEDMEDIA",
	// 	"-dPDFFitPage",
	// 	"-sOUTPUTFILE=" + outputPdf,
	// }, processedPdfs...)...)
	// var stderr bytes.Buffer
	// cmd.Stderr = &stderr
	// err = cmd.Run()
	// if err != nil {
	// 	log.Println("Ghostscript error:", stderr.String())
	// 	return nil, fmt.Errorf("failed to merge PDFs: %w", err)
	// }
	//End of GS merge

	// Use pdfcpu to merge PDFs
	outputPdf := filepath.Join(tmpDir, "merged_output.pdf")
	err = api.MergeCreateFile(processedPdfs, outputPdf, false, nil) // true for overwriting the output file
	if err != nil {
		return nil, fmt.Errorf("failed to merge PDFs with pdfcpu: %w", err)
	}

	// Path for the output file with updated properties
	finalPdfPath := filepath.Join(tmpDir, "final_output_with_properties.pdf")

	// Set properties using api.AddPropertiesFile
	properties := map[string]string{
		"Title":   "eKYC Agreement",
		"Author":  "mb",
		"Creator": "mb",
	}

	err = api.AddPropertiesFile(outputPdf, finalPdfPath, properties, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to add properties: %w", err)
	}
	//End of PDFcpu merge

	// Read the final merged PDF
	finalPdfData, err := os.ReadFile(finalPdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read merged PDF: %w", err)
	}

	return finalPdfData, nil
}

func ByteExtractPdfPages(pdfData []byte, startPage int, endPage int) ([]byte, error) {
	// Create a temporary input file with the PDF data
	tmpInputFile, err := os.CreateTemp("", "input-pdf-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp input file: %w", err)
	}
	defer os.Remove(tmpInputFile.Name())

	// Write the PDF data to the temporary input file
	_, err = tmpInputFile.Write(pdfData)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF data to temp file: %w", err)
	}

	// Create a temporary output file
	tmpOutputFile, err := os.CreateTemp("", "extracted-pages-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer os.Remove(tmpOutputFile.Name())

	// Construct Ghostscript command
	outputPdf := tmpOutputFile.Name()
	cmd := exec.Command("gs",
		"-q",                   // Quiet mode
		"-dNOPAUSE", "-dBATCH", // Non-interactive mode
		"-sDEVICE=pdfwrite",                      // Output as PDF
		"-dCompatibilityLevel=1.6",               // Set PDF compatibility level
		fmt.Sprintf("-dFirstPage=%d", startPage), // Start page
		fmt.Sprintf("-dLastPage=%d", endPage),    // End page
		"-sOUTPUTFILE="+outputPdf,                // Output file
		tmpInputFile.Name(),                      // Input PDF file
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ghostscript error: %s", stderr.String())
	}

	// Read the output PDF file
	extractedPdf, err := os.ReadFile(outputPdf)
	if err != nil {
		return nil, fmt.Errorf("failed to read output PDF: %w", err)
	}

	return extractedPdf, nil
}

// Cut Pdfs by the given range and returns output pdf
func ExtractPdfPages(inputPdf string, startPage int, endPage int) ([]byte, error) {
	// Create a temporary output file
	tmpFile, err := os.CreateTemp("", "extracted-pages-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Construct Ghostscript command
	outputPdf := tmpFile.Name()
	cmd := exec.Command("gs",
		"-q",                   // Quiet mode
		"-dNOPAUSE", "-dBATCH", // Non-interactive mode
		"-sDEVICE=pdfwrite",                      // Output as PDF
		"-dCompatibilityLevel=1.6",               // Set PDF compatibility level
		fmt.Sprintf("-dFirstPage=%d", startPage), // Start page
		fmt.Sprintf("-dLastPage=%d", endPage),    // End page
		"-sOUTPUTFILE="+outputPdf,                // Output file
		inputPdf,                                 // Input PDF file
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ghostscript error: %s", stderr.String())
	}

	// Read the output PDF file
	extractedPdf, err := os.ReadFile(outputPdf)
	if err != nil {
		return nil, fmt.Errorf("failed to read output PDF: %w", err)
	}

	return extractedPdf, nil
}

// ExtractFrame extracts an image from the given video file at the specified time
func ExtractFrame(videoFile string, timestamp string) ([]byte, error) {
	// Set the output format to "jpg"
	outputFormat := "jpg"

	// Create a temporary file for the extracted image
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("frame-*.%s", outputFormat))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Construct the FFmpeg command
	cmdArgs := []string{"-i", videoFile}

	// Add the timestamp argument if provided
	if timestamp != "" {
		cmdArgs = append(cmdArgs, "-ss", timestamp)
	}

	// Add the arguments to extract a single frame and output as a JPG
	cmdArgs = append(cmdArgs,
		"-frames:v", "1", // Extract only one frame
		"-q:v", "2", // Set quality (lower is better, for JPEG)
		"-y",           // Force overwrite of existing file
		tmpFile.Name(), // Output image file
	)

	// Execute the FFmpeg command
	cmd := exec.Command("ffmpeg", cmdArgs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %s", stderr.String())
	}

	// Read the extracted image
	imageData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	return imageData, nil
}

// IsPdfPasswordProtected checks if a PDF is password-protected using pdfcpu.
func IsPdfPasswordProtected(inputPdf string, pdfData []byte, inputType string) (bool, error) {
	switch inputType {
	case "file":
		// Validate using the file path
		if _, err := os.Stat(inputPdf); err != nil {
			if os.IsNotExist(err) {
				return false, errors.New("docbinary.go: file does not exist :func passwordcheck")
			}
			return false, err
		}
		err := api.ValidateFile(inputPdf, nil)
		return err != nil, nil

	case "byte":
		// Validate using the PDF data
		if len(pdfData) == 0 {
			return false, errors.New("docbinary.go: pdf data is empty :func passwordcheck")
		}
		reader := bytes.NewReader(pdfData)
		err := api.Validate(reader, nil)
		return err != nil, nil

	default:
		return false, errors.New("docbinary.go: invalid input type; use 'file' or 'byte'")
	}
}
