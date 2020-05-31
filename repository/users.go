package repository

import (
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"fmt"
)

type Users struct {}

func(u *Users) FetchAll() []model.User {
	rows, _ := db.Conn.Query("SELECT id, email FROM users")
	var users []model.User

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.Id, &user.Email)
		if err != nil {
			fmt.Println(err.Error())
		}
		users = append(users, user)
	}

	return users
}

