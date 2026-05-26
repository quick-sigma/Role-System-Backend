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

type Stat struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CharacterStat struct {
	ID          uint `json:"id"`
	CharacterID uint `json:"character_id"`
	StatID      uint `json:"stat_id"`
	Value       int  `json:"value"`
}

type CharacterProfile struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Age      int             `json:"age"`
	RaceName string          `json:"race_name"`
	Stats    map[string]int  `json:"stats"`
}

type CharacterRepository interface {
	Repository[Character]
	GetProfile(ctx context.Context, id uint) (*CharacterProfile, error)
}

type Race struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
