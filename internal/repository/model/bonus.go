package model

type DataBonus struct {
	Id       int    `json:"id"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	ReqBonus
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
	CreatedOn    string  `json:"created_on"`
}

type ReqBonus struct {
	CategoryId int     `json:"category_id"`
	NewMember  float64 `json:"new_member"`
	CbSl       float64 `json:"cb_sl"`
	RbSl       float64 `json:"rb_sl"`
	CbCa       float64 `json:"cb_ca"`
	RollCa     float64 `json:"roll_ca"`
	CbSp       float64 `json:"cb_sp"`
	RbSp       float64 `json:"rb_sp"`
	Refferal   float64 `json:"refferal"`
	Promo      float64 `json:"promo"`
	TransDate  string  `json:"trans_date"`
}
