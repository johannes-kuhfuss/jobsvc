package dto

import "time"

type JobResponse struct {
	Id            string    `json:"id"`
	CorrelationId string    `json:"correlationId"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	ModifiedAt    time.Time `json:"modifiedAt"`
	ModifiedBy    string    `json:"modifiedBy"`
	Status        string    `json:"status"`
	Source        string    `json:"source"`
	Destination   string    `json:"destination"`
	Type          string    `json:"type"`
	SubType       string    `json:"subType"`
	Action        string    `json:"action"`
	ActionDetails string    `json:"actionDetails"`
	Progress      int32     `json:"progress"`
	History       string    `json:"history"`
	ExtraData     string    `json:"extraData"`
	Priority      string    `json:"priority"`
	Rank          int32     `json:"rank"`
}
