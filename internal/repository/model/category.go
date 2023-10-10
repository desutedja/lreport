package model

type ReqCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
