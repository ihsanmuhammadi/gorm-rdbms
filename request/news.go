package request

import "time"

type News struct {
	Topic		string		`json:"topic"`
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
	PublishedAt	time.Time	`json:"publishedAt"`
	Content		string		`json:"content"`
}

type Source struct {
	Name		string		`json:"name"`
}
