package models

type Pages struct {
	Total    int `json:"total"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	Filtered int `json:"filtered"`
}
