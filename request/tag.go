package request

type CreateTag struct {
	Name	string	`json:"nameTag" validate:"required"`
}
