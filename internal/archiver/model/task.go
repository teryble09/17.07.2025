package model

const (
	Waiting        = "waiting"
	Loaded         = "loaded"
	FailedToLoad   = "failed to load"
	NotAllowedType = "not allowed type"
	Archived       = "archived"
)

type TaskID struct {
	Id string
}

type Url struct {
	Address string
	Status  string
}
