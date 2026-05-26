package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"motor-de-rol/backend/domain"
	"motor-de-rol/backend/repository"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestAPI(t *testing.T) (*gorm.DB, *repository.SQLiteCharacterRepo, http.Handler) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Character{}, &domain.Race{}, &domain.CharacterStat{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	repo := repository.NewSQLiteCharacterRepo(db)
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("Motor de Rol API", "1.0.0"))

	controller := NewCharacterController(repo)
	controller.Register(api)

	return db, repo, router
}

func doRequest(t *testing.T, handler http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(data)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestGetByID_NotFound(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	rec := doRequest(t, handler, http.MethodGet, "/characters/999", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetByID_Success(t *testing.T) {
	db, _, handler := setupTestAPI(t)
	ctx := context.Background()

	character := &domain.Character{Name: "Aldric", Age: 30, RaceID: 1}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	rec := doRequest(t, handler, http.MethodGet, "/characters/1", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	assert.Equal(t, "Aldric", result.Name)
	assert.Equal(t, 30, result.Age)
}

func TestCreate_Success(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	input := map[string]interface{}{
		"name":    "Aldric",
		"age":     30,
		"race_id": 1,
	}
	rec := doRequest(t, handler, http.MethodPost, "/characters", input)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	assert.Equal(t, "Aldric", result.Name)
	assert.Equal(t, 30, result.Age)
	assert.NotZero(t, result.ID)
}

func TestCreate_Validation_NameEmpty(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	input := map[string]interface{}{
		"name":    "",
		"age":     30,
		"race_id": 1,
	}
	rec := doRequest(t, handler, http.MethodPost, "/characters", input)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreate_Validation_AgeOutOfRange(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	input := map[string]interface{}{
		"name":    "Aldric",
		"age":     0,
		"race_id": 1,
	}
	rec := doRequest(t, handler, http.MethodPost, "/characters", input)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreate_Validation_NameTooLong(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	input := map[string]interface{}{
		"name":    string(make([]byte, 101)),
		"age":     30,
		"race_id": 1,
	}
	rec := doRequest(t, handler, http.MethodPost, "/characters", input)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdate_Success(t *testing.T) {
	db, _, handler := setupTestAPI(t)
	ctx := context.Background()

	character := &domain.Character{Name: "Aldric", Age: 30, RaceID: 1}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	input := map[string]interface{}{
		"name":    "Aldric the Wise",
		"age":     31,
		"race_id": 2,
	}
	rec := doRequest(t, handler, http.MethodPut, "/characters/1", input)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	assert.Equal(t, "Aldric the Wise", result.Name)
	assert.Equal(t, 31, result.Age)
}

func TestUpdate_NotFound(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	input := map[string]interface{}{
		"name":    "Nobody",
		"age":     30,
		"race_id": 1,
	}
	rec := doRequest(t, handler, http.MethodPut, "/characters/999", input)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDelete_Success(t *testing.T) {
	db, _, handler := setupTestAPI(t)
	ctx := context.Background()

	character := &domain.Character{Name: "Gandalf", Age: 2000, RaceID: 1}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	rec := doRequest(t, handler, http.MethodDelete, "/characters/1", nil)

	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Raw("SELECT COUNT(*) FROM characters WHERE id = ?", character.ID).Scan(&count)
	assert.Equal(t, int64(0), count)
}

func TestList_EmptyAndPopulated(t *testing.T) {
	_, repo, handler := setupTestAPI(t)
	ctx := context.Background()

	rec := doRequest(t, handler, http.MethodGet, "/characters", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	var emptyList []domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &emptyList))
	assert.Len(t, emptyList, 0)

	characters := []domain.Character{
		{Name: "Aldric", Age: 30, RaceID: 1},
		{Name: "Beren", Age: 25, RaceID: 1},
	}
	for i := range characters {
		require.NoError(t, repo.Create(ctx, &characters[i]))
	}

	rec = doRequest(t, handler, http.MethodGet, "/characters", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	var populatedList []domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &populatedList))
	assert.Len(t, populatedList, 2)
}

func TestGetProfile_Success(t *testing.T) {
	db, _, handler := setupTestAPI(t)
	ctx := context.Background()

	race := &domain.Race{Name: "Human"}
	require.NoError(t, db.WithContext(ctx).Create(race).Error)

	character := &domain.Character{
		Name:   "Aldric",
		Age:    30,
		RaceID: race.ID,
	}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	stats := &domain.CharacterStat{
		CharacterID:  character.ID,
		Strength:     15,
		Dexterity:    12,
		Intelligence: 18,
	}
	require.NoError(t, db.WithContext(ctx).Create(stats).Error)

	rec := doRequest(t, handler, http.MethodGet, "/characters/1/profile", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	var profile domain.CharacterProfile
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &profile))
	assert.Equal(t, "Aldric", profile.Name)
	assert.Equal(t, "Human", profile.RaceName)
	assert.Equal(t, 15, profile.Strength)
	assert.Equal(t, 12, profile.Dexterity)
	assert.Equal(t, 18, profile.Intelligence)
}

func TestGetProfile_NotFound(t *testing.T) {
	_, _, handler := setupTestAPI(t)

	rec := doRequest(t, handler, http.MethodGet, "/characters/999/profile", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
