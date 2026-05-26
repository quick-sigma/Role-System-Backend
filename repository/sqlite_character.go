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

		result := tx.Delete(&domain.Character{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete character: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return nil
	})
}

func (r *SQLiteCharacterRepo) GetProfile(ctx context.Context, id uint) (*domain.CharacterProfile, error) {
	charQuery := `
		SELECT c.id, c.name, c.age, COALESCE(r.name, '')
		FROM characters c
		LEFT JOIN races r ON c.race_id = r.id
		WHERE c.id = ?
	`

	var profile domain.CharacterProfile
	var raceName string
	row := r.db.WithContext(ctx).Raw(charQuery, id).Row()
	err := row.Scan(&profile.ID, &profile.Name, &profile.Age, &raceName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get character profile: %w", err)
	}
	profile.RaceName = raceName
	profile.Stats = make(map[string]int)

	statsQuery := `
		SELECT s.name, cs.value
		FROM character_stats cs
		JOIN stats s ON cs.stat_id = s.id
		WHERE cs.character_id = ?
	`

	rows, err := r.db.WithContext(ctx).Raw(statsQuery, id).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get character stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var statName string
		var value int
		if err := rows.Scan(&statName, &value); err != nil {
			return nil, fmt.Errorf("failed to scan stat row: %w", err)
		}
		profile.Stats[statName] = value
	}

	return &profile, nil
}
