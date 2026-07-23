package scrapper

import (
	"database/sql"
	"go-scraper-learning/internal/logging"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var pageCounter int

func Init(baseUrlToScrap string, filePathName string, db *sql.DB, engine string) {
	historialFile, historialErr := createHistorialBaseFile(baseUrlToScrap, filePathName)
	if historialErr != nil {
		logging.Error(
			"Scrapping might not work correctly as the historial file had an error during its creation or initialization.",
			logging.ErrorType(historialErr),
		)
	}
	defer historialFile.Close()

	pageCounter = 0

	//booksPages := make([]BooksPage, 0, 50) // 50 because I know there are 50 html pages
	// this is the best to make slices, even if i dont know how much im going to need
	// if it needs more space, it will ask for it, but for the preassigned space wont

	for {
		books, isNext := scrapping(baseUrlToScrap)

		booksPage := BooksPage {
			WebPage: uint16(pageCounter),
			Books: books,
		}

		// Petition to the URL
		/*booksPages = append(
			booksPages,
			booksPage,
		)*/

		InsertBooks(db, &booksPage, engine)

		if !isNext {
			break
		}
	}

	

	// ======= DEBUG ======= //
	/*for _, bookpage := range booksPages {
		fmt.Printf(
			"\nWeb page: %d",
			bookpage.WebPage,
		)

		for _, book := range bookpage.Books {
			fmt.Printf(
				"\n\tBook:\n\t\tTitle: %s\n\t\tRating: %d\n\t\tPrice: %d",
				book.Title,
				book.Rating,
				book.Price,
			)
		}
	}*/
}

func scrapping(url string) ([]BookRegisterDTO, bool) {
	pageCounter++
	url = adaptPageUrl(url, pageCounter)
	res, err := http.Get(url)
	if err != nil {
		logging.Fatal(
			"Petition error.",
			logging.ErrorType(err),
		)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logging.Fatal(
			"Wrong status code received.",
			logging.IntType("status_code", res.StatusCode),
		)
	}

	// Charge the html
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logging.Fatal(
			"Parsing error.",
			logging.ErrorType(err),
		)
	}

	booksRegisters := make([]BookRegisterDTO, 0, 20) // 20 because I know how many books appears in every page, if not, an operation should be done

	doc.Find(".product_pod").Each(func(_ int, sel *goquery.Selection) {
		title, _ := sel.Find("h3 > a").Attr("title")
		
		ratingClassAtr, _ := sel.Find("p").Attr("class")
		rat := strings.Split(ratingClassAtr, " ")[1]
		var rating int8
		switch rat {
		case "One":
			rating = 1
		case "Two":
			rating = 2
		case "Three":
			rating = 3
		case "Four":
			rating = 4
		case "Five":
			rating = 5
		}

		priceStr := sel.Find(".price_color").Text()
		price, err := strconv.ParseFloat(strings.ReplaceAll(priceStr, "£", ""), 64)
		if err != nil {
			logging.Error(
				"Parse string to float32 error.",
				logging.ErrorType(err),
			)
			return
		}
		booksRegisters = append(
			booksRegisters,
			BookRegisterDTO{
				Title: title,
				Rating: rating,
				Price: uint64(math.Round(price * 100)), // Store in cents
			},
		)
	})

	isMorePages := false
	_, exists := doc.Find(".next a").Attr("href")
	if exists {
		isMorePages = true
	}

	if len(booksRegisters) > 0 {
		return booksRegisters, isMorePages
	}
	return nil, isMorePages
}

func adaptPageUrl(url string, number int) (string) {
	return strings.Replace(url, "X", strconv.Itoa(number), 1)
}