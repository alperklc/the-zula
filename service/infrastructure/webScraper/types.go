package webScraper

type PageContent struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Length      int    `json:"int"`
	Excerpt     string `json:"excerpt"`
	SiteName    string `json:"siteName"`
	Image       string `json:"image"`
	Favicon     string `json:"favicon"`
	HTMLContent string `json:"htmlContent"`
	MDContent   string `json:"mdContent"`
}
