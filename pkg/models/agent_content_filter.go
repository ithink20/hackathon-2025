package models

type FilterResponse struct {
	IsProblematic bool   `json:"isProblematic"`
	HelpText      string `json:"helpText"`
}

type FilterRequest struct {
	UserContent string `json:"user_content"`
}
