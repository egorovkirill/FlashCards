package entities

type Lists struct {
	Id    int    `json:"id"`
	Title string `json:"title" binding:"required"`
}

type UserLists struct {
	Id     int `json:"id"`
	UserId int `json:"userId"`
	ListId int `json:"listId"`
}

type Cards struct {
	Id        int    `json:"id"`
	Front     string `json:"front" binding:"required"`
	Back      string `json:"back"`
	ImageLink string `json:"imageLink"`
}

type ListCards struct {
	Id     int `json:"id"`
	UserId int `json:"userId"`
	ListId int `json:"listId"`
}
