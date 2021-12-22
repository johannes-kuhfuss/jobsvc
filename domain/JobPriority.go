package domain

type JobPriority string

const (
	PriorityRealtime JobPriority = "realtime"
	PriorityHigh     JobPriority = "high"
	PriorityMedium   JobPriority = "medium"
	PriorityLow      JobPriority = "low"
	PriorityIdle     JobPriority = "idle"
)
