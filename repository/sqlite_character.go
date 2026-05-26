package repository

import (
	"context"
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
		if err == gorm.ErrRecordNotFound {
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
