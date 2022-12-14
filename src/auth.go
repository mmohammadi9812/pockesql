// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ncruces/zenity"
	"github.com/pkg/browser"
)

type AuthInfo struct {
	ConsumerKey string `json:"consumer_key"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

func rawSendJson(url string, data map[string]string) ([]byte, error) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(data)
	resp, err := http.Post(url, "application/json", payloadBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func sendJson(url string, data map[string]string) (string, error) {
	bodyBytes, err := rawSendJson(url, data)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func RequestToken(consumerKey string) (string, error) {
	data := map[string]string{
		"consumer_key": consumerKey,
		"redirect_uri": "https://getpocket.com/connected_applications",
	}
	body, err := sendJson("https://getpocket.com/v3/oauth/request", data)
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(body)
	if err != nil {
		return "", err
	}
	requestToken := q["code"][0]
	return requestToken, nil
}

func AutherizeUser(consumerKey, requestToken string) (string, string, error) {
	data := map[string]string{
		"consumer_key": consumerKey,
		"code":         requestToken,
	}
	body, err := sendJson("https://getpocket.com/v3/oauth/authorize", data)
	if err != nil {
		return "", "", err
	}

	q, err := url.ParseQuery(body)
	if err != nil {
		return "", "", err
	}

	username := q["username"][0]
	accessToken := q["access_token"][0]
	return username, accessToken, nil
}

func WriteFile(data AuthInfo, filename string) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func OpenConfirmUrl(requestToken string) {
	codesUrl := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		requestToken, "https://getpocket.com/connected_applications")
	browser.OpenURL(codesUrl)
	zenity.Question("Click ok once you have authorized",
		zenity.Title("Authorization"),
		zenity.OKLabel("Completed"))
	// Sleep a second before trying to rush, in case it rushes faster than site acknowledges
	time.Sleep(time.Second)
}

func ReadAuth() (AuthInfo, error) {
	fbytes, err := os.ReadFile("auth.json")
	if err != nil {
		return AuthInfo{}, err
	}
	inf := AuthInfo{}
	err = json.Unmarshal(fbytes, &inf)
	if err != nil {
		return AuthInfo{}, err
	}
	return inf, nil
}
