package models

// DoseResponse is model for dose/response data.
type DoseResponse struct {
	Dose     float64 `json:"dose"`
	Response float64 `json:"response"`
}

// DoseResponses is a collection of DoseResponse.
type DoseResponses []DoseResponse
