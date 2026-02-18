package dto

type CreateBookInput struct {
	ID        string
	Title     string
	Authors   []string
	ISBN      string
	Category  string
	Publisher string
	Year      int
}

type UpdateBookInput struct {
	ID        string
	Title     string
	Authors   []string
	ISBN      string
	Category  string
	Publisher string
	Year      int
}
