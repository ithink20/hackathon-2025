package models

type FilterResponse struct {
	IsProblematic   bool   `json:"isProblematic"`
	HelpText        string `json:"helpText"`
	ContentCategory string `json:"contentCategory"`
	EnglishContent  string `json:"englishContent"`
}

type FilterRequest struct {
	UserContent string `json:"user_content"`
}

const (
	CategoryQuestion     = "Question"
	CategoryPost         = "Post"
	CategoryOthers       = "Others"
	CategoryAppreciation = "Appreciation"
)
