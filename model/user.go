package model

type User struct {
	ID string `json:"id" xml:"id"`
	*UserInput
}

type UserInput struct {
	FirstName string `json:"firstName" xml:"firstName" validate:"required"`
	LastName  string `json:"lastName" xml:"lastName"  validate:"required"`
	Email     string `json:"email" xml:"email" validate:"required,email"`
}
