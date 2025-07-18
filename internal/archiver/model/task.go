package model

const (
	Waiting        = "waiting"
	Loaded         = "loaded"
	FailedToLoad   = "failed to load"
	NotAllowedType = "not allowed type"
	Archived       = "archived"
)

type TaskID struct {
	Id string `json:"id"`
}

type Url struct {
	Address string `json:"address"`
	Status  string `json:"status"`
}
