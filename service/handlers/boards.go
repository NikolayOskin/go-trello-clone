package handlers

import (
	"context"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct{}

func (h Board) Create(ctx context.Context, b model.Board) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := db.Boards.InsertOne(ctx, b)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (h Board) Update(ctx context.Context, b model.Board) error {
	f := bson.M{"_id": b.ID, "user_id": b.UserId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Boards.UpdateOne(ctx, f, bson.M{"$set": b}); err != nil {
		return err
	}
	return nil
}

func (h Board) FillBoardWithListsAndCards(ctx context.Context, board *model.Board) error {
	lists, err := repository.Lists.GetByBoardId(ctx, board.ID.Hex())
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		board.Lists = make([]model.List, 0) // for json empty array instead of null
		return nil
	}

	cards, err := repository.Cards.GetByBoardId(ctx, board.ID.Hex())
	if err != nil {
		return err
	}

	// combining cards to map by listId
	cardsMap := mapCardsByListId(cards, len(lists))

	for i, list := range lists {
		if _, found := cardsMap[list.ID.Hex()]; !found {
			lists[i].Cards = make([]model.Card, 0) // for json cards empty array instead of null
		} else {
			lists[i].Cards = cardsMap[list.ID.Hex()]
		}
	}
	board.Lists = lists

	return nil
}

func mapCardsByListId(cards []model.Card, listsCount int) map[string][]model.Card {
	cardsMap := make(map[string][]model.Card)
	avgListLength := len(cards)/listsCount + 1

	for _, c := range cards {
		if _, found := cardsMap[c.ListId]; found {
			cardsMap[c.ListId] = append(cardsMap[c.ListId], c)
		} else {
			s := make([]model.Card, 0, avgListLength) // preallocate slice with average cards per list capacity
			cardsMap[c.ListId] = append(s, c)
		}
	}

	return cardsMap
}
