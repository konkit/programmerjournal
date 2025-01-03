package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type ListTaskOutput struct {
	Body []task.Task
}

func ListTasks(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "ListTasks",
		Method:      http.MethodGet,
		Path:        "/api/tasks/list/{date}",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *struct {
		Date string `path:"date" example:"2024-05-05" doc:"Day to select the list from"`
	}) (*ListTaskOutput, error) {
		resp := &ListTaskOutput{}
		tasks, err := taskRepo.ListTasks(input.Date)
		if err != nil {
			return nil, err
		}
		resp.Body = tasks
		return resp, nil
	})
}
