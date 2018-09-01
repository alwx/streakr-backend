package utils

import (
	"github.com/google/uuid"
	"encoding/base64"
)

type Header struct {
	Name  string
	Value string
}

func GetBasicHeaders(token string) []Header {
	return []Header{
		{Name: "Cache-Control", Value: "no-cache"},
		{Name: "User-Agent", Value: "Streakr App"},
		{Name: "X-Bunq-Client-Authentication", Value: token},
		{Name: "X-Bunq-Client-Request-Id", Value: uuid.New().String()},
		{Name: "X-Bunq-Geolocation", Value: "0 0 0 0 000"},
		{Name: "X-Bunq-Language", Value: "en_US"},
		{Name: "X-Bunq-Region", Value: "nl_NL"},
	}
}

func GetSignature(firstLine string, headers []Header, body string, privateKey string) (string, error) {
	allHeaders := ""
	for _, element := range headers {
		allHeaders += element.Name + ": " + element.Value + "\n"
	}
	signature := firstLine + "\n" + allHeaders + "\n" + body

	signer, err := ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return "", err
	}

	signedData, err := signer.Sign([]byte(signature))
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(signedData)
	return encoded, nil
}
