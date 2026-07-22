package scrapper

import "go-scraper-learning/internal/logging"

func Init(baseUrlToScrap string, filePathName string) {
	historialFile, historialErr := createHistorialBaseFile(baseUrlToScrap, filePathName)
	if historialErr != nil {
		logging.Error(
			"Scrapping might not work correctly as the historial file had an error during its creation or initialization.",
			logging.ErrorType(historialErr),
		)
	}
	defer historialFile.Close()

	// Petition to the URL

}