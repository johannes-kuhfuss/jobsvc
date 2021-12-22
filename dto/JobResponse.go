package dto

import "time"

type JobResponse struct {
	Id            string    `json:"id"`
	CorrelationId string    `json:"correlation_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	ModifiedAt    time.Time `json:"modified_at"`
	ModifiedBy    string    `json:"modified_by"`
	Status        string    `json:"status"`
	Source        string    `json:"source"`
	Destination   string    `json:"destination"`
	Type          string    `json:"type"`
	SubType       string    `json:"sub_type"`
	Action        string    `json:"action"`
	ActionDetails string    `json:"action_details"`
	Progress      int32     `json:"progress"`
	History       string    `json:"history"`
	ExtraData     string    `json:"extra_data"`
	Priority      string    `json:"priority"`
	Rank          int32     `json:"rank"`
}
