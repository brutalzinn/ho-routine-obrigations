package webmodels

type ObrigationConfirmRequest struct {
	Value  string `json:"value"`
	Device string `json:"firebase_token"`
}
