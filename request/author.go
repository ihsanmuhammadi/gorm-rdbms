package request

type CreateAuthor struct {
	Name		string		`json:"nameAuthor" validate:"required"`
	Email		string		`json:"emailAuthor" validate:"required,email"`
}
