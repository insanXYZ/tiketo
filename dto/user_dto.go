package dto

type User struct {
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
