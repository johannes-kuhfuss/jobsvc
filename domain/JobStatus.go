package domain

import "strings"

type JobStatus string

const (
	StatusCreated  JobStatus = "created"
	StatusQueued   JobStatus = "queued"
	StatusRunning  JobStatus = "running"
	StatusPaused   JobStatus = "paused"
	StatusFinished JobStatus = "finished"
	StatusFailed   JobStatus = "failed"
)

func IsValidJobStatus(statusVal string) bool {
	val := strings.TrimSpace(strings.ToLower(statusVal))
	if (val == string(StatusCreated)) ||
		(val == string(StatusQueued)) ||
		(val == string(StatusRunning)) ||
		(val == string(StatusPaused)) ||
		(val == string(StatusFinished)) ||
		(val == string(StatusFailed)) {
		return true
	} else {
		return false
	}
}
