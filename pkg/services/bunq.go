package services

import (
	"github.com/spf13/viper"
	"bytes"
	"net/http"
	"io/ioutil"
	"streakr-backend/pkg/utils"
	"github.com/tidwall/gjson"
	"fmt"
	"github.com/google/uuid"
)

func execute(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()
	return resp, nil
}

func BunqInstallation(user *NewUser) (string, error) {
	url := viper.GetString("bunq.api") + "/v1/installation"

	json := "{\"client_public_key\": " + fmt.Sprintf("%q", user.PublicKey) + "}"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return "", err
	}

	headers := []utils.Header{
		{Name: "Cache-Control", Value: "no-cache"},
		{Name: "User-Agent", Value: "Streakr App"},
		{Name: "X-Bunq-Client-Request-Id", Value: uuid.New().String()},
		{Name: "X-Bunq-Geolocation", Value: "0 0 0 0 000"},
		{Name: "X-Bunq-Language", Value: "en_US"},
		{Name: "X-Bunq-Region", Value: "nl_NL"},
	}
	for _, element := range headers {
		req.Header.Set(element.Name, element.Value)
	}

	response, err := execute(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return gjson.Get(string(body), "Response.1.Token.token").String(), nil
}

func BunqDeviceServer(user *NewUser) (int64, error) {
	url := viper.GetString("bunq.api") + "/v1/device-server"

	json := "{\"description\": \"Streakr\", \"secret\": \"" + user.APIKey + "\"}"

	headers := utils.GetBasicHeaders(user.Token)
	signedSignature, err := utils.GetSignature("POST /v1/device-server", headers, json, user.PrivateKey)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return 0, err
	}
	for _, element := range headers {
		req.Header.Set(element.Name, element.Value)
	}
	req.Header.Set("X-Bunq-Client-Signature", signedSignature)

	response, err := execute(req)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	return gjson.Get(string(body), "Response.0.Id.id").Int(), nil
}

type SessionServer struct {
	UserPersonId int64
	DisplayName  string
	Token string
}

func BunqSessionServer(user *NewUser) (SessionServer, error) {
	url := viper.GetString("bunq.api") + "/v1/session-server"

	json := "{\"secret\": \"" + user.APIKey + "\"}"

	headers := utils.GetBasicHeaders(user.Token)
	signedSignature, err := utils.GetSignature("POST /v1/session-server", headers, json, user.PrivateKey)
	if err != nil {
		return SessionServer{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return SessionServer{}, err
	}
	for _, element := range headers {
		req.Header.Set(element.Name, element.Value)
	}
	req.Header.Set("X-Bunq-Client-Signature", signedSignature)

	response, err := execute(req)
	if err != nil {
		return SessionServer{}, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return SessionServer{}, err
	}

	userPerson := gjson.Get(string(body), "Response.2.UserPerson")

	return SessionServer{
		UserPersonId: userPerson.Get("id").Int(),
		DisplayName: userPerson.Get("display_name").String(),
		Token: gjson.Get(string(body), "Response.1.Token.token").String(),
	}, nil
}

func BunqGetUser(user User) (User, string, error) {
	endpoint := fmt.Sprintf("/v1/user-person/%d", user.UserPersonId)
	url := viper.GetString("bunq.api") + endpoint

	json := ""

	headers := utils.GetBasicHeaders(user.Token)
	signedSignature, err := utils.GetSignature("GET " + endpoint, headers, json, user.PrivateKey)
	if err != nil {
		return user, "", err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return user, "", err
	}
	for _, element := range headers {
		req.Header.Set(element.Name, element.Value)
	}
	req.Header.Set("X-Bunq-Client-Signature", signedSignature)

	response, err := execute(req)
	if err != nil {
		return user, "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, "", err
	}

	return user, gjson.Get(string(body), "Response.0.UserPerson").String(), nil
}

func BunqSetNotificationFilters(user User, bunqUser string) (string, error) {
	endpoint := fmt.Sprintf("/v1/user-person/%d", user.UserPersonId)
	url := viper.GetString("bunq.api") + endpoint

	json := "{\"notification_filters\": [{\"notification_delivery_method\": \"URL\", \"notification_target\": \"https://streakr.alwx.me/push\", \"category\": \"MUTATION\"}]}"

	headers := utils.GetBasicHeaders(user.Token)
	signedSignature, err := utils.GetSignature("PUT " + endpoint, headers, json, user.PrivateKey)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return "string", err
	}
	for _, element := range headers {
		req.Header.Set(element.Name, element.Value)
	}
	req.Header.Set("X-Bunq-Client-Signature", signedSignature)

	response, err := execute(req)
	if err != nil {
		return "string", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "string", err
	}

	println(string(body))

	return "ok", nil
}