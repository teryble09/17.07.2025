package dto

type AddURLRequest struct {
	TaskId string
	Adress string `json:"url"`
}

type AddURLResponse struct{}

type CreateTaskRequest struct{}

type CreateTaskResponse struct {
	Id string `json:"task_id"`
}
