package model

type DataTransaction struct {
	Id string `json:"id"`
	ReqTransaction
	UserId    string  `json:"user_id"`
	ConvDp    float64 `json:"conv_dp"`
	ConvTr    float64 `json:"conv_tr"`
	SubTotal  float64 `json:"sub_total"`
	Ats       float64 `json:"ats"`
	Total     float64 `json:"total"`
	CreatedOn string  `json:"created_on"`
}

type ReqTransaction struct {
	CategoryId   int     `json:"category_id"`
	Regis        int     `json:"regis"`
	RegisDp      int     `json:"regis_dp"`
	ActivePlayer int     `json:"active_player"`
	TransDp      int     `json:"trans_dp"`
	TransWd      int     `json:"trans_wd"`
	TotalDp      float64 `json:"total_dp"`
	TotalWd      float64 `json:"total_wd"`
	Wl           float64 `json:"wl"`
	TransDate    string  `json:"trans_date"`
}
