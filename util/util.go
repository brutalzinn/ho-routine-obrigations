package util

import (
	"encoding/base64"
	"log"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeHtmlImageTag(content string, size int) string {
	qrCodeImageData, taskError := qrcode.Encode(content, qrcode.High, size)
	if taskError != nil {
		log.Fatalln("Error generating QR code. ", taskError)
	}
	encodedData := base64.StdEncoding.EncodeToString(qrCodeImageData)
	return "<img src=\"data:image/png;base64, " + encodedData + "\">"
}
