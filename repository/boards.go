package repository

import (
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"fmt"
)

type Boards struct {}

func(u *Boards) FetchByUser(user model.User) []model.Board {
	rows, _ := db.Conn.Query("SELECT title, user_id FROM boards")
	var boards []model.Board

	for rows.Next() {
		board := model.Board{}

		err := rows.Scan(&board.Title, &board.UserId)
		if err != nil {
			fmt.Println(err.Error())
		}

		boards = append(boards, board)
	}

	return boards
}