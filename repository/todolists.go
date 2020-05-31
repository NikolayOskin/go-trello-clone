package repository

import (
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"fmt"
)

type TodoLists struct {}

func(u *TodoLists) FetchByUser(user model.User) []model.TodoList {
	rows, _ := db.Conn.Query("SELECT id, title FROM todolists")
	var todolists []model.TodoList

	for rows.Next() {
		todoList := model.TodoList{}

		err := rows.Scan(&todoList.Id, &todoList.Title)
		if err != nil {
			fmt.Println(err.Error())
		}

		todolists = append(todolists, todoList)
	}

	return todolists
}
