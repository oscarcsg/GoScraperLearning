package scrapper

type Book struct {
	Id     int32   `json:"id"`
	Title  string  `json:"title"`
	Rating int8    `json:"rating"`
	Price  float32 `json:"pricing"`
}

type BookRegisterDTO struct {
	Title  string  `json:"title"`
	Rating int8    `json:"rating"`
	Price  float32 `json:"pricing"`
}

type BooksPage struct {
	WebPage uint16
	Books   []BookRegisterDTO
}
