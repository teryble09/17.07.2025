package model

const (
	Waiting         = "waiting"
	FailedToLoad    = "failed to load"
	NotAllowedType  = "not allowed type"
	FailedToArchive = "failed to archive"
	Archived        = "archived"
)

type TaskID struct {
	Id string
}

type Url struct {
	Address string `json:"address"`
	Status  string `json:"status"`
}
