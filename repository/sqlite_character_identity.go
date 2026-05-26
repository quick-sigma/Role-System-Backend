package repository

import (
	"context"
	"fmt"

	"motor-de-rol/backend/domain"

	"gorm.io/gorm"
)

type SQLiteCharacterIdentityRepo struct {
	db *gorm.DB
}

func NewSQLiteCharacterIdentityRepo(db *gorm.DB) *SQLiteCharacterIdentityRepo {
	return &SQLiteCharacterIdentityRepo{db: db}
}

func (r *SQLiteCharacterIdentityRepo) CreateFullCharacter(ctx context.Context, char *domain.CharacterFull, avatar *domain.CharacterAvatar, attrs []domain.CharacterAttribute, stats []domain.CharacterStat, skills []domain.CharacterSkill) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(char).Error; err != nil {
			return fmt.Errorf("failed to create character: %w", err)
		}

		if avatar != nil {
			avatar.CharacterID = char.ID
			if err := tx.Create(avatar).Error; err != nil {
				return fmt.Errorf("failed to create character avatar: %w", err)
			}
		}

		for i := range attrs {
			attrs[i].CharacterID = char.ID
		}
		if len(attrs) > 0 {
			if err := tx.Create(&attrs).Error; err != nil {
				return fmt.Errorf("failed to create character attributes: %w", err)
			}
		}

		for i := range stats {
			stats[i].CharacterID = char.ID
		}
		if len(stats) > 0 {
			if err := tx.Create(&stats).Error; err != nil {
				return fmt.Errorf("failed to create character stats: %w", err)
			}
		}

		for i := range skills {
			skills[i].CharacterID = char.ID
		}
		if len(skills) > 0 {
			if err := tx.Create(&skills).Error; err != nil {
				return fmt.Errorf("failed to create character skills: %w", err)
			}
		}

		return nil
	})
}

func (r *SQLiteCharacterIdentityRepo) GetFullCharacterSheet(ctx context.Context, charID uint) (*domain.CharacterFullSheetDTO, error) {
	var character domain.CharacterFull
	if err := r.db.WithContext(ctx).First(&character, charID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("character not found with id %d", charID)
		}
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	dto := &domain.CharacterFullSheetDTO{
		Character: &character,
	}

	var avatar domain.CharacterAvatar
	if err := r.db.WithContext(ctx).Where("character_id = ?", charID).First(&avatar).Error; err == nil {
		dto.Avatar = &avatar
	}

	var attributes []domain.CharacterAttribute
	if err := r.db.WithContext(ctx).Where("character_id = ?", charID).Find(&attributes).Error; err != nil {
		return nil, fmt.Errorf("failed to get character attributes: %w", err)
	}
	dto.Attributes = attributes

	var stats []domain.CharacterStat
	if err := r.db.WithContext(ctx).Where("character_id = ?", charID).Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get character stats: %w", err)
	}
	dto.Stats = stats

	var skills []domain.CharacterSkill
	if err := r.db.WithContext(ctx).Where("character_id = ?", charID).Find(&skills).Error; err != nil {
		return nil, fmt.Errorf("failed to get character skills: %w", err)
	}
	dto.Skills = skills

	var psychology []domain.CharacterPsychology
	if err := r.db.WithContext(ctx).Where("character_id = ?", charID).Find(&psychology).Error; err != nil {
		return nil, fmt.Errorf("failed to get character psychology: %w", err)
	}
	dto.Psychology = psychology

	var organizations []domain.CharacterOrganization
	if err := r.db.WithContext(ctx).Where("character_id = ?", charID).Find(&organizations).Error; err != nil {
		return nil, fmt.Errorf("failed to get character organizations: %w", err)
	}
	dto.Organizations = organizations

	return dto, nil
}

func (r *SQLiteCharacterIdentityRepo) DeleteCharacter(ctx context.Context, charID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM character_psychologies WHERE character_id = ?", charID).Error; err != nil {
			return fmt.Errorf("failed to delete character psychologies: %w", err)
		}

		if err := tx.Exec("DELETE FROM character_organizations WHERE character_id = ?", charID).Error; err != nil {
			return fmt.Errorf("failed to delete character organizations: %w", err)
		}

		if err := tx.Exec("DELETE FROM character_attributes WHERE character_id = ?", charID).Error; err != nil {
			return fmt.Errorf("failed to delete character attributes: %w", err)
		}

		if err := tx.Exec("DELETE FROM character_stats WHERE character_id = ?", charID).Error; err != nil {
			return fmt.Errorf("failed to delete character stats: %w", err)
		}

		if err := tx.Exec("DELETE FROM character_skills WHERE character_id = ?", charID).Error; err != nil {
			return fmt.Errorf("failed to delete character skills: %w", err)
		}

		if err := tx.Exec("DELETE FROM character_avatars WHERE character_id = ?", charID).Error; err != nil {
			return fmt.Errorf("failed to delete character avatar: %w", err)
		}

		result := tx.Delete(&domain.CharacterFull{}, charID)
		if result.Error != nil {
			return fmt.Errorf("failed to delete character: %w", result.Error)
		}

		return nil
	})
}
