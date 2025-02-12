package tags

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/tag"
)

type ListTagsInput struct {
}

type ListTagsOutput struct {
	Body   []tag.Tag
	Status int
}

func ListTagsHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "ListTags",
		Method:      http.MethodGet,
		Path:        "/api/tags/list",
		Tags:        []string{"Tag"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ListTagsInput) (*ListTagsOutput, error) {
		resp := &ListTagsOutput{}

		tags := []tag.Tag{}
		err := db.Model(tag.Tag{}).Find(&tags).Error

		if err != nil {
			return nil, err
		}

		resp.Body = tags
		resp.Status = http.StatusOK

		return resp, nil
	})
}
