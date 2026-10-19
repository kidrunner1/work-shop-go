package models

// json คือชื่อ key ที่รับมา ส่วน validate คือกฎตรวจสอบข้อมูล
type RegisterRequest struct {
	Username        string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Email           string `json:"email" validate:"required,email,max=254"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	LineID          string `json:"line_id" ` // ไม่บังคับกรอกตามภาพ ถ้าอยากตรวจสอบให้ใช้ validate:"required,min=3,max=30"
	Phone           string `json:"phone" validate:"required,number"`
	BusinessType    string `json:"business_type" validate:"required"`
	WebsiteName     string `json:"website_name" validate:"required,min=2,max=30,hostname_rfc1123,lowercase,excludesall=."` // เฉพาะชื่อหน้า .sogoodweb.com
}
