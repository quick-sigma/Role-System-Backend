package transport

import (
	"context"
	"net/http"

	"motor-de-rol/backend/domain"

	"github.com/danielgtaylor/huma/v2"
)

type CharacterController struct {
	repo domain.CharacterRepository
}

func NewCharacterController(repo domain.CharacterRepository) *CharacterController {
	return &CharacterController{repo: repo}
}

func (c *CharacterController) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "CreateCharacter",
		Method:      http.MethodPost,
		Path:        "/characters",
		Summary:     "Create a new character",
	}, c.Create)

	huma.Register(api, huma.Operation{
		OperationID: "GetCharacter",
		Method:      http.MethodGet,
		Path:        "/characters/{id}",
		Summary:     "Get character by ID",
	}, c.GetByID)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateCharacter",
		Method:      http.MethodPut,
		Path:        "/characters/{id}",
		Summary:     "Update character",
	}, c.Update)

	huma.Register(api, huma.Operation{
		OperationID: "ListCharacters",
		Method:      http.MethodGet,
		Path:        "/characters",
		Summary:     "List all characters",
	}, c.List)

	huma.Register(api, huma.Operation{
		OperationID: "GetCharacterProfile",
		Method:      http.MethodGet,
		Path:        "/characters/{id}/profile",
		Summary:     "Get character full profile with stats and race",
	}, c.GetProfile)
}

type CreateCharacterInput struct {
	Body struct {
		Name   string `json:"name" minLength:"1" maxLength:"100"`
		Age    int    `json:"age" minimum:"1" maximum:"9999"`
		RaceID uint   `json:"race_id"`
	}
}

func (c *CharacterController) Create(ctx context.Context, input *CreateCharacterInput) (*struct {
	Body *domain.Character
}, error) {
	character := &domain.Character{
		Name:   input.Body.Name,
		Age:    input.Body.Age,
		RaceID: input.Body.RaceID,
	}

	if err := c.repo.Create(ctx, character); err != nil {
		return nil, huma.Error500InternalServerError("failed to create character", err)
	}

	return &struct {
		Body *domain.Character
	}{Body: character}, nil
}

type CharacterIDInput struct {
	ID uint `path:"id" minimum:"1"`
}

func (c *CharacterController) GetByID(ctx context.Context, input *CharacterIDInput) (*struct {
	Body *domain.Character
}, error) {
	character, err := c.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("character not found", err)
	}

	return &struct {
		Body *domain.Character
	}{Body: character}, nil
}

type UpdateCharacterInput struct {
	ID uint `path:"id" minimum:"1"`
	Body struct {
		Name   string `json:"name" minLength:"1" maxLength:"100"`
		Age    int    `json:"age" minimum:"1" maximum:"9999"`
		RaceID uint   `json:"race_id"`
	}
}

func (c *CharacterController) Update(ctx context.Context, input *UpdateCharacterInput) (*struct {
	Body *domain.Character
}, error) {
	character, err := c.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("character not found", err)
	}

	character.Name = input.Body.Name
	character.Age = input.Body.Age
	character.RaceID = input.Body.RaceID

	if err := c.repo.Update(ctx, character); err != nil {
		return nil, huma.Error500InternalServerError("failed to update character", err)
	}

	return &struct {
		Body *domain.Character
	}{Body: character}, nil
}

func (c *CharacterController) List(ctx context.Context, input *struct{}) (*struct {
	Body []domain.Character
}, error) {
	characters, err := c.repo.List(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list characters", err)
	}

	return &struct {
		Body []domain.Character
	}{Body: characters}, nil
}

func (c *CharacterController) GetProfile(ctx context.Context, input *CharacterIDInput) (*struct {
	Body *domain.CharacterProfile
}, error) {
	profile, err := c.repo.GetProfile(ctx, input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("character profile not found", err)
	}

	return &struct {
		Body *domain.CharacterProfile
	}{Body: profile}, nil
}
