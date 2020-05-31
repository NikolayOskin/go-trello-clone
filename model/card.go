package model

type Card struct {
	Text string `json:"text"`
	TodoListId uint `json:"todolist_id"`
	BoardId uint `json:"board_id"`
	Order uint `json:"order"`
}