package scrapper

type Book struct {
	Id     int32   `json:"id"`
	Title  string  `json:"title"`
	Rating int8    `json:"rating"`
	Price  uint64 `json:"pricing"`// Cents
}

type BookRegisterDTO struct {
	Title  string  `json:"title"`
	Rating int8    `json:"rating"`
	Price  uint64 `json:"pricing"` // Cents
}

type BooksPage struct {
	WebPage uint16
	Books   []BookRegisterDTO
}
