package domain

import "context"

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]T, error)
}

type Character struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	RaceID uint   `json:"race_id"`
}

type CharacterProfile struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Age             int    `json:"age"`
	RaceName        string `json:"race_name"`
	Strength        int    `json:"strength"`
	Dexterity       int    `json:"dexterity"`
	Intelligence    int    `json:"intelligence"`
}

type CharacterRepository interface {
	Repository[Character]
	GetProfile(ctx context.Context, id uint) (*CharacterProfile, error)
}

type Race struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CharacterStat struct {
	ID           uint `json:"id"`
	CharacterID  uint `json:"character_id"`
	Strength     int  `json:"strength"`
	Dexterity    int  `json:"dexterity"`
	Intelligence int  `json:"intelligence"`
}
