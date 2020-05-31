package model

type TodoList struct {
	Title string `json:"title"`
	BoardId uint `json:"board_id"`
	Order uint `json:"order"`
}