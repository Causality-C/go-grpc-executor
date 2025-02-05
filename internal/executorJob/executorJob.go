package executorjob

// Dataclass to hold the job data
type ExecutorJob struct {
	TaskID     int32
	OrderSize  int32
	Product    string
	Filename   string
	OutputPath string
	UserID     int32
	Text       string
}
