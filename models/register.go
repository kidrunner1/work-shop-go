package models

// json คือชื่อ key ที่รับมา ส่วน validate คือกฎตรวจสอบข้อมูล
type RegisterRequest struct {
	Username        string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Email           string `json:"email" validate:"required,email,max=254"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
