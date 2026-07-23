package scrapper

import (
	"fmt"
	"go-scraper-learning/internal/logging"
	"os"
	"strings"
)

func createHistorialBaseFile(baseUrl string, filePathName string) (*os.File, error) {
	var sb strings.Builder
	fileName := ""

	if filePathName != "" {
		parts := strings.SplitSeq(filePathName, "/")
		for part := range parts {
			if strings.Contains(part, ".") {
				fileName = part
				continue
			}
			fmt.Fprint(&sb, part)
		}
	}

	filePath := sb.String()

	var err error

	if filePath != "" {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			logging.Error(
				"Historial directory path error.",
				logging.ErrorType(err),
			)
			return nil, err
		}
	}
	
	var file *os.File
	flags := os.O_CREATE | // Create if no exists
			 os.O_RDWR   | // Open to read AND write
			 os.O_APPEND   // Add text to the end, do not delete anything

	if fileName != "" {
		file, err = os.OpenFile(filePathName, flags, 0644) // 
		if err != nil {
			logging.Error(
				"File creation error",
				logging.ErrorType(err),
			)
			return nil, err
		}
	}

	/*if baseUrl != "" {
		_, err = file.WriteString(baseUrl)
		if err != nil {
			logging.Error(
				"Base url has not been able to be written in the file.",
				logging.StringType("base_url", baseUrl),
				logging.StringType("real_file_name", file.Name()),
				logging.ErrorType(err),
			)
			return nil, err
		}
	}*/

	return file, nil
}