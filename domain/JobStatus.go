package domain

type JobStatus string

const (
	StatusCreated  JobStatus = "created"
	StatusQueued   JobStatus = "queued"
	StatusRunning  JobStatus = "running"
	StatusPaused   JobStatus = "paused"
	StatusFinished JobStatus = "finished"
	StatusFailed   JobStatus = "failed"
)
