package data

type User struct {
	ID			string	`json:"id" sql:"id"`
	Email		string  `json:"email" validate:"required" sql:"email"`
	Password	string  `json:"password" validate:"required" sql:"password"`
	// CreatedOn 	string  `json:"-"`
	// UpdatedOn 	string	`json:"-"`	
}