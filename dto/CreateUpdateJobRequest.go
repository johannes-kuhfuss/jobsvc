package dto

type CreateUpdateJobRequest struct {
	CorrelationId string `json:"correlation_id"`
	Name          string `json:"name"`
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	Type          string `json:"type"`
	SubType       string `json:"sub_type"`
	Action        string `json:"action"`
	ActionDetails string `json:"action_details"`
	ExtraData     string `json:"extra_data"`
	Priority      string `json:"priority"`
	Rank          int32  `json:"rank"`
}
