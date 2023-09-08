package webmodels

type MobileConfirmRequest struct {
	Value         string `json:"value"`
	DeviceName    string `json:"device_name"`
	FirebaseToken string `json:"firebase_token"`
}
