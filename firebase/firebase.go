package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var serverKey = os.Getenv("FIREBASE_API_KEY")

type FCMNotification struct {
	To           string              `json:"to"`
	Notification FCMNotificationData `json:"notification"`
}

type FCMNotificationData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func New(device string, title string, message string) FCMNotification {
	return FCMNotification{
		To: device,
		Notification: FCMNotificationData{
			Title: title,
			Body:  message,
		},
	}
}

func (fcmNotify FCMNotification) Send() {
	url := "https://fcm.googleapis.com/fcm/send"
	payload, err := json.Marshal(fcmNotify)
	if err != nil {
		fmt.Println("Error marshaling notification:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+serverKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}
