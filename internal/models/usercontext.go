package models

type UserContext struct {
	Query    string `json:"query"`
	Count    int    `json:"count"`
	Position int    `json:"position"`
}
