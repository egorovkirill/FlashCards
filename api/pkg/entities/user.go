package entities

type User struct {
	Id       int    `json:"id" postgresql:"id"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
