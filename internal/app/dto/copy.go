package dto

type CreateCopyInput struct {
	ID            string
	BookID        string
	Barcode       string
	ConditionNote string
}

type UpdateCopyInput struct {
	ID            string
	Barcode       string
	Status        string
	ConditionNote string
}
