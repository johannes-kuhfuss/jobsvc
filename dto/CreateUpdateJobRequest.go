package dto

type CreateUpdateJobRequest struct {
	CorrelationId string `json:"correlationId"`
	Name          string `json:"name"`
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	Type          string `json:"type"`
	SubType       string `json:"subType"`
	Action        string `json:"action"`
	ActionDetails string `json:"actionDetails"`
	ExtraData     string `json:"extraData"`
	Priority      string `json:"priority"`
	Rank          int32  `json:"rank"`
}
