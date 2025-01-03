package handlershuma

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type ListTaskOutput struct {
	Body []task.Task
}

func ListTasks(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Get(api, "/api/tasks/list/{date}", func(ctx context.Context, input *struct {
		Date string `path:"date" example:"2024-05-05" doc:"Day to select the list from"`
	}) (*ListTaskOutput, error) {
		resp := &ListTaskOutput{}
		tasks, err := taskRepo.GetTasksFromDate(input.Date)
		if err != nil {
			return nil, err
		}
		resp.Body = tasks
		return resp, nil
	})
}
