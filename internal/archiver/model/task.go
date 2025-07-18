package model

const (
	Loaded = iota
	Processing
	FailedLoad
	NotAllowedType
)

type TaskID struct {
	Id string `json:"id"`
}

type Task struct {
	urls    []string
	status  []int
	archive []byte
}
