package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"motor-de-rol/backend/domain"

	"github.com/danielgtaylor/huma/v2"
)

type GenericIDInput struct {
	ID uint `path:"id" minimum:"1" doc:"ID único del recurso"`
}

type GenericCreateBodyInput struct {
	Body map[string]any `doc:"Datos del recurso a crear"`
}

type GenericUpdateBodyInput struct {
	GenericIDInput
	Body map[string]any `doc:"Datos del recurso a actualizar"`
}

func RegisterGenericCRUDL[T any](api huma.API, repo domain.Repository[T], basePath string, groupTag string) {
	resourceName := groupTag

	huma.Register(api, huma.Operation{
		OperationID: "List" + resourceName,
		Method:      http.MethodGet,
		Path:        "/" + basePath,
		Summary:     "List all " + resourceName,
		Tags:        []string{groupTag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body []T
	}, error) {
		entities, err := repo.List(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list records", err)
		}
		return &struct {
			Body []T
		}{Body: entities}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "Create" + resourceName,
		Method:      http.MethodPost,
		Path:        "/" + basePath,
		Summary:     "Create a new " + resourceName,
		Tags:        []string{groupTag},
	}, func(ctx context.Context, input *GenericCreateBodyInput) (*struct {
		Body T
	}, error) {
		var entity T
		data, _ := json.Marshal(input.Body)
		if err := json.Unmarshal(data, &entity); err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid request body", err)
		}
		if err := repo.Create(ctx, &entity); err != nil {
			return nil, huma.Error500InternalServerError("failed to create record", err)
		}
		return &struct {
			Body T
		}{Body: entity}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "Get" + resourceName,
		Method:      http.MethodGet,
		Path:        "/" + basePath + "/{id}",
		Summary:     "Get " + resourceName + " by ID",
		Tags:        []string{groupTag},
	}, func(ctx context.Context, input *GenericIDInput) (*struct {
		Body T
	}, error) {
		entity, err := repo.GetByID(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound("record not found", err)
		}
		return &struct {
			Body T
		}{Body: *entity}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "Update" + resourceName,
		Method:      http.MethodPut,
		Path:        "/" + basePath + "/{id}",
		Summary:     "Update " + resourceName,
		Tags:        []string{groupTag},
	}, func(ctx context.Context, input *GenericUpdateBodyInput) (*struct {
		Body T
	}, error) {
		if _, err := repo.GetByID(ctx, input.ID); err != nil {
			return nil, huma.Error404NotFound("record not found", err)
		}
		var entity T
		data, _ := json.Marshal(input.Body)
		if err := json.Unmarshal(data, &entity); err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid request body", err)
		}
		if err := repo.Update(ctx, &entity); err != nil {
			return nil, huma.Error500InternalServerError("failed to update record", err)
		}
		return &struct {
			Body T
		}{Body: entity}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "Delete" + resourceName,
		Method:      http.MethodDelete,
		Path:        "/" + basePath + "/{id}",
		Summary:     "Delete " + resourceName,
		Tags:        []string{groupTag},
	}, func(ctx context.Context, input *GenericIDInput) (*struct{}, error) {
		if err := repo.Delete(ctx, input.ID); err != nil {
			return nil, huma.Error500InternalServerError("failed to delete record", err)
		}
		return nil, nil
	})
}
