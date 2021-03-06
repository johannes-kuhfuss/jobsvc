package dto

type CreateUpdateJobRequest struct {
	CorrelationId string `json:"correlationId" san:"trim,xss"`
	Name          string `json:"name" san:"trim,xss"`
	Source        string `json:"source" san:"trim,xss"`
	Destination   string `json:"destination" san:"trim,xss"`
	Type          string `json:"type" san:"trim,xss"`
	SubType       string `json:"sub_type" san:"trim,xss"`
	Action        string `json:"action" san:"trim,xss"`
	ActionDetails string `json:"action_details" san:"trim,xss"`
	ExtraData     string `json:"extra_data" san:"trim,xss"`
	Priority      string `json:"priority" san:"trim,xss,lower"`
	Rank          int32  `json:"rank" san:"def=0,min=0,max=2147483647"`
}
