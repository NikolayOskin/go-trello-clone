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

func CreateBoard(ctx context.Context, b model.Board) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := db.Boards.InsertOne(ctx, b)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateBoard(ctx context.Context, b model.Board) error {
	f := bson.M{"_id": b.ID, "user_id": b.UserId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Boards.UpdateOne(ctx, f, bson.M{"$set": b}); err != nil {
		return err
	}
	return nil
}

func FillBoardWithListsAndCards(ctx context.Context, b *model.Board) error {
	listsRepo := repository.Lists{}
	lists, err := listsRepo.GetByBoardId(ctx, b.ID.Hex())
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		b.Lists = make([]model.List, 0) // for json empty array instead of null
		return nil
	}

	cardsRepo := repository.Cards{}
	cards, err := cardsRepo.GetByBoardId(ctx, b.ID.Hex())
	if err != nil {
		return err
	}

	// combining cards to map by listId
	cardsMap := generateCardsMap(cards)

	for i, list := range lists {
		if _, found := cardsMap[list.ID.Hex()]; !found {
			lists[i].Cards = make([]model.Card, 0) // for json cards empty array instead of null
		} else {
			lists[i].Cards = cardsMap[list.ID.Hex()]
		}
	}
	b.Lists = lists

	return nil
}

func generateCardsMap(cards []model.Card) map[string][]model.Card {
	cardsMap := make(map[string][]model.Card)

	for _, card := range cards {
		if _, found := cardsMap[card.ListId]; found {
			cardsMap[card.ListId] = append(cardsMap[card.ListId], card)
		} else {
			cardsMap[card.ListId] = []model.Card{card}
		}
	}

	return cardsMap
}
