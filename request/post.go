package request

type CreatePost struct {
	Title		string			`json:"title" validate:"required"`
	Content		string			`json:"content" validate:"required"`
	Author		CreateAuthor
	Category	[]CreateCategory	`validate:"required,dive"`
	Tags		[]CreateTag		`validate:"required,dive"`
}
