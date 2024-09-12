package request

type GetNews struct {
	Status		string		`json:"status"`
	TotalResults	int		`json:"totalResults"`
	Articles	[]Articles
}

type Articles struct {
	Source		Source
	Author		string		`json:"author"`
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	Url		string		`json:"url"`
	UrlToImage	string		`json:"urlToImage"`
	Content		string		`json:"content"`
}

type Source struct {
	Name		string		`json:"name"`
}
