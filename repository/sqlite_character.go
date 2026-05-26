package repository

import (
	"context"
	"database/sql"
	"fmt"

	"motor-de-rol/backend/domain"

	"gorm.io/gorm"
)

type SQLiteCharacterRepo struct {
	*SQLiteRepository[domain.Character]
	db *gorm.DB
}

func NewSQLiteCharacterRepo(db *gorm.DB) *SQLiteCharacterRepo {
	return &SQLiteCharacterRepo{
		SQLiteRepository: NewSQLiteRepository[domain.Character](db),
		db:               db,
	}
}

func (r *SQLiteCharacterRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM character_stats WHERE character_id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete character stats: %w", err)
		}

		if err := tx.Delete(&domain.Character{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete character: %w", err)
		}

		return nil
	})
}

func (r *SQLiteCharacterRepo) GetProfile(ctx context.Context, id uint) (*domain.CharacterProfile, error) {
	query := `
		SELECT 
			c.id, c.name, c.age,
			r.name as race_name,
			s.strength, s.dexterity, s.intelligence
		FROM characters c
		LEFT JOIN character_stats s ON c.id = s.character_id
		LEFT JOIN races r ON c.race_id = r.id
		WHERE c.id = ?
	`

	var profile domain.CharacterProfile
	row := r.db.WithContext(ctx).Raw(query, id).Row()

	var raceName *string
	var strength, dexterity, intelligence *int
	err := row.Scan(
		&profile.ID,
		&profile.Name,
		&profile.Age,
		&raceName,
		&strength,
		&dexterity,
		&intelligence,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get character profile: %w", err)
	}

	if raceName != nil {
		profile.RaceName = *raceName
	}
	if strength != nil {
		profile.Strength = *strength
	}
	if dexterity != nil {
		profile.Dexterity = *dexterity
	}
	if intelligence != nil {
		profile.Intelligence = *intelligence
	}

	return &profile, nil
}
