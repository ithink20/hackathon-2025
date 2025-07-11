package models

// ConfluenceSearchResponse represents the response from Confluence search API
type ConfluenceSearchResponse struct {
	Results []ConfluencePage `json:"results"`
}

// ConfluencePage represents a single page from Confluence
type ConfluencePage struct {
	Content ConfluenceContent `json:"content"`
}

// ConfluenceContent represents the content details of a Confluence page
type ConfluenceContent struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

// PageInfo represents the simplified page information we want to return
type PageInfo struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
}
