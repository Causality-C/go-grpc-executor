package executorJob

// Dataclass to hold the job data
type ExecutorJob struct {
	TaskId     int32
	OrderSize  int32
	Product    string
	Filename   string
	OutputPath string
	UserId     int32
	Text       string
}
