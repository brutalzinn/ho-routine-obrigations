package models

type ObrigationQueuePending struct {
	Id            string
	QrCodeValue   string
	IdObrigation  int
	IdDevice      int
	DeviceName    string
	TokenFirebase string
}
