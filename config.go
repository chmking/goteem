package horde

type Config struct {
	Tasks   []*Task
	WaitMin int64
	WaitMax int64
}
