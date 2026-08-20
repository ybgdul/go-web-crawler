package models

type Result struct{ 
	URL string `json:"url"`
	Title string `json:"title"`
	Err error `json:"-"`
}

type CrawlRequest struct {
	Urls []string `json:"urls"`
}