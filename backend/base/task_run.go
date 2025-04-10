package base

// TaskRunStatus is the status of a task run.
type TaskRunStatus string

const (
	// TaskRunUnknown is the task run status for UNKNOWN.
	TaskRunUnknown TaskRunStatus = "UNKNOWN"
	// TaskRunPending is the task run status of PENDING.
	TaskRunPending TaskRunStatus = "PENDING"
	// TaskRunRunning is the task run status for RUNNING.
	TaskRunRunning TaskRunStatus = "RUNNING"
	// TaskRunDone is the task run status for DONE.
	TaskRunDone TaskRunStatus = "DONE"
	// TaskRunFailed is the task run status for FAILED.
	TaskRunFailed TaskRunStatus = "FAILED"
	// TaskRunCanceled is the task run status for CANCELED.
	TaskRunCanceled TaskRunStatus = "CANCELED"
	// TaskRunNotStarted is the task run status for NOT_STARTED.
	TaskRunNotStarted TaskRunStatus = "NOT_STARTED"
	// TaskRunSkipped is the task run status for SKIPPED.
	TaskRunSkipped TaskRunStatus = "SKIPPED"
)

func (s TaskRunStatus) String() string {
	return string(s)
}
