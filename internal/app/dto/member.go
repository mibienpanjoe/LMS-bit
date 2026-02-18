package dto

type RegisterMemberInput struct {
	ID    string
	Name  string
	Email string
	Phone string
}

type UpdateMemberInput struct {
	ID    string
	Name  string
	Email string
	Phone string
}
