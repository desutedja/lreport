package model

type DataTransaction struct {
	Id           string `json:"id"`
	CategoryName int    `json:"category_name"`
	ReqTransaction
	UserId    string  `json:"user_id"`
	Username  string  `json:"username"`
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

type DataReportTransaction struct {
	Regis        int     `json:"regis"`
	RegisDp      int     `json:"regis_dp"`
	ActivePlayer int     `json:"active_player"`
	TransDp      int     `json:"trans_dp"`
	TransWd      int     `json:"trans_wd"`
	TotalDp      float64 `json:"total_dp"`
	TotalWd      float64 `json:"total_wd"`
	Wl           float64 `json:"wl"`
	ConvDp       float64 `json:"conv_dp"`
	ConvTr       float64 `json:"conv_tr"`
	SubTotal     float64 `json:"sub_total"`
	Ats          float64 `json:"ats"`
	Bonus        float64 `json:"bonus"`
	Total        float64 `json:"total"`
	Day          int     `json:"day"`
	Month        int     `json:"month"`
	Year         int     `json:"year"`
	Period       string  `json:"period"`
}

type RespReportTransaction struct {
	DataReport []DataReportTransaction `json:"data_report"`
	DataKey    []string                `json:"data_key"`
}
