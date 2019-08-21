package models

// PayloadImages payload image
type PayloadImages struct {
	Image string `json:"image"`
	Port  int    `json:"port"`
}

// Payload primary payload
type Payload struct {
	StateID string          `json:"stateID"`
	SvcPort int             `json:"svcPort"`
	Images  []PayloadImages `json:"images"`
}
