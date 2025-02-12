package tags

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/tag"
)

type GetTagInput struct {
	TagName string `path:"tagname" example:"tagone" doc:"Requested tag"`
}

type GetTagOutput struct {
	Body   []entry.Entry
	Status int
}

func GetTagsHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "GetTag",
		Method:      http.MethodGet,
		Path:        "/api/tags/{tagname}",
		Tags:        []string{"Tag"},
	}
	huma.Register(api, op, func(ctx context.Context, input *GetTagInput) (*GetTagOutput, error) {
		resp := &GetTagOutput{}

		t := tag.Tag{Name: input.TagName}
		if err := db.First(&t).Error; err != nil {
			resp.Status = http.StatusNotFound
			return resp, nil
		}

		et := []tag.EntryTag{}
		err := db.Model(tag.EntryTag{}).Where("tag_id = ?", t.ID).Find(&et).Error
		if err != nil {
			resp.Status = http.StatusInternalServerError
			return resp, err
		}

		var entries []entry.Entry
		for _, ett := range et {
			var entriesFromTag []entry.Entry
			db.Model(entry.Entry{}).Where("ID = ?", ett.EntryID).Find(&entriesFromTag)
			entries = append(entries, entriesFromTag...)
		}

		resp.Body = entries
		resp.Status = http.StatusOK

		return resp, nil
	})
}
