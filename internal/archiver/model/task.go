package model

const (
	Waiting = iota
	Loaded
	Processing
	FailedLoad
	NotAllowedType
)

type TaskID struct {
	Id string `json:"id"`
}

type Task struct {
	Urls    []string
	Status  []int
	Archive []byte
}
