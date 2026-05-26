package transport

import (
	"context"
	"net/http"

	"motor-de-rol/backend/domain"

	"github.com/danielgtaylor/huma/v2"
)

type GenericIDInput struct {
	ID uint `path:"id" minimum:"1" doc:"ID único del recurso"`
}

type GenericCreateInput[T any] struct {
	Body T
}

type GenericUpdateInput[T any] struct {
	GenericIDInput
	Body T
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
	}, func(ctx context.Context, input *GenericCreateInput[T]) (*struct {
		Body T
	}, error) {
		if err := repo.Create(ctx, &input.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to create record", err)
		}
		return &struct {
			Body T
		}{Body: input.Body}, nil
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
	}, func(ctx context.Context, input *GenericUpdateInput[T]) (*struct {
		Body T
	}, error) {
		if _, err := repo.GetByID(ctx, input.ID); err != nil {
			return nil, huma.Error404NotFound("record not found", err)
		}
		if err := repo.Update(ctx, &input.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to update record", err)
		}
		return &struct {
			Body T
		}{Body: input.Body}, nil
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
