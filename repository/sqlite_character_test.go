package repository

import (
	"context"
	"testing"

	"motor-de-rol/backend/domain"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Character{}, &domain.Race{}, &domain.CharacterStat{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)
	return db
}

func TestCreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	character := &domain.Character{
		Name:   "Aldric",
		Age:    30,
		RaceID: 1,
	}

	err := repo.Create(ctx, character)
	require.NoError(t, err)
	assert.NotZero(t, character.ID)

	retrieved, err := repo.GetByID(ctx, character.ID)
	require.NoError(t, err)
	assert.Equal(t, "Aldric", retrieved.Name)
	assert.Equal(t, 30, retrieved.Age)
	assert.Equal(t, uint(1), retrieved.RaceID)
}

func TestGetProfile_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	race := &domain.Race{Name: "Human"}
	require.NoError(t, db.Create(race).Error)

	character := &domain.Character{
		Name:   "Aldric",
		Age:    30,
		RaceID: race.ID,
	}
	require.NoError(t, repo.Create(ctx, character))

	stats := &domain.CharacterStat{
		CharacterID:  character.ID,
		Strength:     15,
		Dexterity:    12,
		Intelligence: 18,
	}
	require.NoError(t, db.Create(stats).Error)

	profile, err := repo.GetProfile(ctx, character.ID)
	require.NoError(t, err)

	assert.Equal(t, character.ID, profile.ID)
	assert.Equal(t, "Aldric", profile.Name)
	assert.Equal(t, 30, profile.Age)
	assert.Equal(t, "Human", profile.RaceName)
	assert.Equal(t, 15, profile.Strength)
	assert.Equal(t, 12, profile.Dexterity)
	assert.Equal(t, 18, profile.Intelligence)
}

func TestGetProfile_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	profile, err := repo.GetProfile(ctx, 999)

	assert.Nil(t, profile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetProfile_WithNullStats(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	race := &domain.Race{Name: "Elf"}
	require.NoError(t, db.Create(race).Error)

	character := &domain.Character{
		Name:   "Legolas",
		Age:    500,
		RaceID: race.ID,
	}
	require.NoError(t, repo.Create(ctx, character))

	profile, err := repo.GetProfile(ctx, character.ID)
	require.NoError(t, err)

	assert.Equal(t, character.ID, profile.ID)
	assert.Equal(t, "Legolas", profile.Name)
	assert.Equal(t, "Elf", profile.RaceName)
	assert.Equal(t, 0, profile.Strength)
	assert.Equal(t, 0, profile.Dexterity)
	assert.Equal(t, 0, profile.Intelligence)
}

func TestDelete_Cascade(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	character := &domain.Character{
		Name:   "Gandalf",
		Age:    2000,
		RaceID: 1,
	}
	require.NoError(t, repo.Create(ctx, character))

	stats := &domain.CharacterStat{
		CharacterID:  character.ID,
		Strength:     10,
		Dexterity:    8,
		Intelligence: 20,
	}
	require.NoError(t, db.Create(stats).Error)

	err := repo.Delete(ctx, character.ID)
	require.NoError(t, err)

	var charCount int64
	db.Raw("SELECT COUNT(*) FROM characters WHERE id = ?", character.ID).Scan(&charCount)
	assert.Equal(t, int64(0), charCount)

	var statsCount int64
	db.Raw("SELECT COUNT(*) FROM character_stats WHERE character_id = ?", character.ID).Scan(&statsCount)
	assert.Equal(t, int64(0), statsCount)
}

func TestDelete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	characters := []domain.Character{
		{Name: "Aldric", Age: 30, RaceID: 1},
		{Name: "Beren", Age: 25, RaceID: 1},
		{Name: "Celeborn", Age: 7000, RaceID: 2},
	}
	for i := range characters {
		require.NoError(t, repo.Create(ctx, &characters[i]))
	}

	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteCharacterRepo(db)
	ctx := context.Background()

	character := &domain.Character{
		Name:   "Aldric",
		Age:    30,
		RaceID: 1,
	}
	require.NoError(t, repo.Create(ctx, character))

	character.Name = "Aldric the Wise"
	character.Age = 31
	require.NoError(t, repo.Update(ctx, character))

	retrieved, err := repo.GetByID(ctx, character.ID)
	require.NoError(t, err)
	assert.Equal(t, "Aldric the Wise", retrieved.Name)
	assert.Equal(t, 31, retrieved.Age)
}
