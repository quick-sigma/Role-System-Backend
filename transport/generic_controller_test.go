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

func setupGenericTestAPI(t *testing.T) (*gorm.DB, *repository.SQLiteRepository[domain.Character], http.Handler) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Character{}, &domain.Race{}, &domain.Stat{}, &domain.CharacterStat{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	repo := repository.NewSQLiteRepository[domain.Character](db)
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("Motor de Rol API", "1.0.0"))

	RegisterGenericCRUDL[domain.Character](api, repo, "characters", "Characters")

	return db, repo, router
}

func doGenericRequest(t *testing.T, handler http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func TestGenericGetByID_NotFound(t *testing.T) {
	_, _, handler := setupGenericTestAPI(t)

	rec := doGenericRequest(t, handler, http.MethodGet, "/characters/999", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGenericGetByID_Success(t *testing.T) {
	db, _, handler := setupGenericTestAPI(t)
	ctx := context.Background()

	character := &domain.Character{Name: "Aldric", Age: 30, RaceID: 1}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	rec := doGenericRequest(t, handler, http.MethodGet, "/characters/1", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	assert.Equal(t, "Aldric", result.Name)
	assert.Equal(t, 30, result.Age)
}

func TestGenericCreate_Success(t *testing.T) {
	_, _, handler := setupGenericTestAPI(t)

	input := map[string]interface{}{
		"name":    "Aldric",
		"age":     30,
		"race_id": 1,
	}
	rec := doGenericRequest(t, handler, http.MethodPost, "/characters", input)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	assert.Equal(t, "Aldric", result.Name)
	assert.Equal(t, 30, result.Age)
	assert.NotZero(t, result.ID)
}

func TestGenericUpdate_Success(t *testing.T) {
	db, _, handler := setupGenericTestAPI(t)
	ctx := context.Background()

	character := &domain.Character{Name: "Aldric", Age: 30, RaceID: 1}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	input := map[string]interface{}{
		"id":      1,
		"name":    "Aldric the Wise",
		"age":     31,
		"race_id": 2,
	}
	rec := doGenericRequest(t, handler, http.MethodPut, "/characters/1", input)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	assert.Equal(t, "Aldric the Wise", result.Name)
	assert.Equal(t, 31, result.Age)
}

func TestGenericUpdate_NotFound(t *testing.T) {
	_, _, handler := setupGenericTestAPI(t)

	input := map[string]interface{}{
		"id":      999,
		"name":    "Nobody",
		"age":     30,
		"race_id": 1,
	}
	rec := doGenericRequest(t, handler, http.MethodPut, "/characters/999", input)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGenericList_EmptyAndPopulated(t *testing.T) {
	_, repo, handler := setupGenericTestAPI(t)
	ctx := context.Background()

	rec := doGenericRequest(t, handler, http.MethodGet, "/characters", nil)
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

	rec = doGenericRequest(t, handler, http.MethodGet, "/characters", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	var populatedList []domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &populatedList))
	assert.Len(t, populatedList, 2)
}

func TestGenericDelete_Success(t *testing.T) {
	db, _, handler := setupGenericTestAPI(t)
	ctx := context.Background()

	character := &domain.Character{Name: "Aldric", Age: 30, RaceID: 1}
	require.NoError(t, db.WithContext(ctx).Create(character).Error)

	rec := doGenericRequest(t, handler, http.MethodDelete, "/characters/1", nil)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rec = doGenericRequest(t, handler, http.MethodGet, "/characters/1", nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGenericMultipleResources(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Character{}, &domain.Race{}, &domain.Stat{})
	require.NoError(t, err)

	charRepo := repository.NewSQLiteRepository[domain.Character](db)
	raceRepo := repository.NewSQLiteRepository[domain.Race](db)

	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("Motor de Rol API", "1.0.0"))

	RegisterGenericCRUDL[domain.Character](api, charRepo, "characters", "Characters")
	RegisterGenericCRUDL[domain.Race](api, raceRepo, "races", "Races")

	ctx := context.Background()
	require.NoError(t, db.WithContext(ctx).Create(&domain.Character{Name: "Aldric", Age: 30, RaceID: 1}).Error)
	require.NoError(t, db.WithContext(ctx).Create(&domain.Race{Name: "Human"}).Error)

	rec := doGenericRequest(t, router, http.MethodGet, "/characters", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	var chars []domain.Character
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &chars))
	assert.Len(t, chars, 1)

	rec = doGenericRequest(t, router, http.MethodGet, "/races", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	var races []domain.Race
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &races))
	assert.Len(t, races, 1)
}
