package transport

import (
	"context"
	"net/http"

	"motor-de-rol/backend/domain"

	"github.com/danielgtaylor/huma/v2"
)

type CharacterIdentityController struct {
	repo domain.CharacterIdentityRepository
}

func NewCharacterIdentityController(repo domain.CharacterIdentityRepository) *CharacterIdentityController {
	return &CharacterIdentityController{repo: repo}
}

func (c *CharacterIdentityController) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "CreateFullCharacter",
		Method:      http.MethodPost,
		Path:        "/characters/full",
		Summary:     "Create a full character with avatar, attributes, stats, and skills",
	}, c.CreateFullCharacter)

	huma.Register(api, huma.Operation{
		OperationID: "GetFullCharacterSheet",
		Method:      http.MethodGet,
		Path:        "/characters/{id}/sheet",
		Summary:     "Get the full character sheet consolidated",
	}, c.GetFullCharacterSheet)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteCharacterIdentity",
		Method:      http.MethodDelete,
		Path:        "/characters/{id}",
		Summary:     "Delete a character and all its dependencies",
	}, c.DeleteCharacter)
}

type AttrInput struct {
	AttributeID uint `json:"attribute_id"`
	Value       int  `json:"value" minimum:"0" maximum:"100"`
}

type StatInput struct {
	StatID uint `json:"stat_id"`
	Value  int  `json:"value" minimum:"0" maximum:"100"`
}

type SkillInput struct {
	SkillID          uint `json:"skill_id"`
	ProficiencyLevel int  `json:"proficiency_level" minimum:"0" maximum:"10"`
}

type AvatarInput struct {
	Height          float64 `json:"height" minimum:"0" maximum:"500"`
	Weight          float64 `json:"weight" minimum:"0" maximum:"1000"`
	EyeColor        string  `json:"eye_color" maxLength:"50"`
	SkinColor       string  `json:"skin_color" maxLength:"50"`
	HairColor       string  `json:"hair_color" maxLength:"50"`
	DistinctiveTraits string `json:"distinctive_traits" maxLength:"500"`
}

type CharacterInput struct {
	RulesetID      uint   `json:"ruleset_id"`
	UserID         uint   `json:"user_id"`
	Name           string `json:"name" minLength:"1" maxLength:"100"`
	Gender         string `json:"gender" maxLength:"50"`
	Pronouns       string `json:"pronouns" maxLength:"50"`
	Age            int    `json:"age" minimum:"1" maximum:"9999"`
	AlignmentID    uint   `json:"alignment_id"`
	Faith          string `json:"faith" maxLength:"100"`
	BackgroundID   uint   `json:"background_id"`
	CityOfOriginID uint   `json:"city_of_origin_id"`
	Biography      string `json:"biography" maxLength:"5000"`
}

type CreateFullCharacterInput struct {
	Body struct {
		Character  CharacterInput `json:"character"`
		Avatar     *AvatarInput   `json:"avatar,omitempty"`
		Attributes []AttrInput    `json:"attributes,omitempty"`
		Stats      []StatInput    `json:"stats,omitempty"`
		Skills     []SkillInput   `json:"skills,omitempty"`
	}
}

func (c *CharacterIdentityController) CreateFullCharacter(ctx context.Context, input *CreateFullCharacterInput) (*struct {
	Body *domain.CharacterFullSheetDTO
}, error) {
	char := &domain.CharacterFull{
		RulesetID:      input.Body.Character.RulesetID,
		UserID:         input.Body.Character.UserID,
		Name:           input.Body.Character.Name,
		Gender:         input.Body.Character.Gender,
		Pronouns:       input.Body.Character.Pronouns,
		Age:            input.Body.Character.Age,
		AlignmentID:    input.Body.Character.AlignmentID,
		Faith:          input.Body.Character.Faith,
		BackgroundID:   input.Body.Character.BackgroundID,
		CityOfOriginID: input.Body.Character.CityOfOriginID,
		Biography:      input.Body.Character.Biography,
	}

	var avatar *domain.CharacterAvatar
	if input.Body.Avatar != nil {
		avatar = &domain.CharacterAvatar{
			Height:          input.Body.Avatar.Height,
			Weight:          input.Body.Avatar.Weight,
			EyeColor:        input.Body.Avatar.EyeColor,
			SkinColor:       input.Body.Avatar.SkinColor,
			HairColor:       input.Body.Avatar.HairColor,
			DistinctiveTraits: input.Body.Avatar.DistinctiveTraits,
		}
	}

	attrs := make([]domain.CharacterAttribute, len(input.Body.Attributes))
	for i, a := range input.Body.Attributes {
		attrs[i] = domain.CharacterAttribute{
			AttributeID: a.AttributeID,
			Value:       a.Value,
		}
	}

	stats := make([]domain.CharacterStat, len(input.Body.Stats))
	for i, s := range input.Body.Stats {
		stats[i] = domain.CharacterStat{
			StatID: s.StatID,
			Value:  s.Value,
		}
	}

	skills := make([]domain.CharacterSkill, len(input.Body.Skills))
	for i, sk := range input.Body.Skills {
		skills[i] = domain.CharacterSkill{
			SkillID:          sk.SkillID,
			ProficiencyLevel: sk.ProficiencyLevel,
		}
	}

	if err := c.repo.CreateFullCharacter(ctx, char, avatar, attrs, stats, skills); err != nil {
		return nil, huma.Error500InternalServerError("failed to create full character", err)
	}

	sheet, err := c.repo.GetFullCharacterSheet(ctx, char.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to retrieve created character", err)
	}

	return &struct {
		Body *domain.CharacterFullSheetDTO
	}{Body: sheet}, nil
}

type CharacterSheetIDInput struct {
	ID uint `path:"id" minimum:"1"`
}

func (c *CharacterIdentityController) GetFullCharacterSheet(ctx context.Context, input *CharacterSheetIDInput) (*struct {
	Body *domain.CharacterFullSheetDTO
}, error) {
	sheet, err := c.repo.GetFullCharacterSheet(ctx, input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("character sheet not found", err)
	}

	return &struct {
		Body *domain.CharacterFullSheetDTO
	}{Body: sheet}, nil
}

func (c *CharacterIdentityController) DeleteCharacter(ctx context.Context, input *CharacterSheetIDInput) (*struct{}, error) {
	if err := c.repo.DeleteCharacter(ctx, input.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete character", err)
	}

	return nil, nil
}
