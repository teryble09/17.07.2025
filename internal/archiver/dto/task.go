package dto

import "github.com/teryble09/17.07.2025/internal/archiver/model"

type AddURLRequest struct {
	TaskId string
	Adress string `json:"url"`
}

type AddURLResponse struct{}

type CreateTaskRequest struct{}

type CreateTaskResponse struct {
	Id string `json:"task_id"`
}
type GetStatusRequest struct {
	TaskId string
}

type GetStatusResponse struct {
	Urls []model.Url `json:"urls"`
}
