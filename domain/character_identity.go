package domain

import "context"

type CharacterFull struct {
	ID             uint   `json:"id"`
	RulesetID      uint   `json:"ruleset_id"`
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	Pronouns       string `json:"pronouns"`
	Age            int    `json:"age"`
	AlignmentID    uint   `json:"alignment_id"`
	Faith          string `json:"faith"`
	BackgroundID   uint   `json:"background_id"`
	CityOfOriginID uint   `json:"city_of_origin_id"`
	Biography      string `json:"biography"`
}

type CharacterAvatar struct {
	CharacterID     uint    `json:"character_id"`
	Height          float64 `json:"height"`
	Weight          float64 `json:"weight"`
	EyeColor        string  `json:"eye_color"`
	SkinColor       string  `json:"skin_color"`
	HairColor       string  `json:"hair_color"`
	DistinctiveTraits string `json:"distinctive_traits"`
}

type CharacterAttribute struct {
	ID          uint `json:"id"`
	CharacterID uint `json:"character_id"`
	AttributeID uint `json:"attribute_id"`
	Value       int  `json:"value"`
}

type CharacterSkill struct {
	ID               uint `json:"id"`
	CharacterID      uint `json:"character_id"`
	SkillID          uint `json:"skill_id"`
	ProficiencyLevel int  `json:"proficiency_level"`
}

type CharacterPsychology struct {
	ID                uint   `json:"id"`
	CharacterID       uint   `json:"character_id"`
	PsychologyTraitID uint   `json:"psychology_trait_id"`
	CustomDescription string `json:"custom_description"`
}

type CharacterOrganization struct {
	ID             uint   `json:"id"`
	CharacterID    uint   `json:"character_id"`
	OrganizationID uint   `json:"organization_id"`
	RankOrRole     string `json:"rank_or_role"`
}

type CharacterFullSheetDTO struct {
	Character    *CharacterFull             `json:"character"`
	Avatar       *CharacterAvatar           `json:"avatar,omitempty"`
	Attributes   []CharacterAttribute       `json:"attributes"`
	Stats        []CharacterStat            `json:"stats"`
	Skills       []CharacterSkill           `json:"skills"`
	Psychology   []CharacterPsychology      `json:"psychology"`
	Organizations []CharacterOrganization   `json:"organizations"`
}

type CharacterIdentityRepository interface {
	CreateFullCharacter(ctx context.Context, char *CharacterFull, avatar *CharacterAvatar, attrs []CharacterAttribute, stats []CharacterStat, skills []CharacterSkill) error
	GetFullCharacterSheet(ctx context.Context, charID uint) (*CharacterFullSheetDTO, error)
	DeleteCharacter(ctx context.Context, charID uint) error
}
