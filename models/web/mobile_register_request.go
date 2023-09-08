package webmodels

type MobileRegisterDevice struct {
	Name          string `json:"name"`
	FirebaseToken string `json:"firebase_token"`
}
