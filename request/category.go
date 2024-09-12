package request

type CreateCategory struct {
	Name 		string		`json:"nameCategory" validate:"required"`
	Description	string		`json:"description" validate:"required"`
}
