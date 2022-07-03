/*
* Copyright 2022 Mohammad Mohamamdi. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
*/

package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use: "auth",
	Short: "Authenticate your account to allow fetching entries",
	Run: func(_ *cobra.Command, _ []string) {
		AuthenticateCmd()
	},
}

func RequestToken(consumerKey string) (string, error) {
	authUrl := "https://getpocket.com/v3/oauth/request"
	data := map[string]string {
		"consumer_key": consumerKey,
		"redirect_uri": "https://getpocket.com/connected_applications",
	}
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(data)
	resp, err := http.Post(authUrl, "application/json", payloadBuffer)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body := string(bodyBytes)

	q, err := url.ParseQuery(body)
	if err != nil {
		return "", err
	}
	requestToken := q["code"][0]
	return requestToken, nil
}

func AutherizeUser(consumerKey, requestToken string) (username, accessToken string, err error) {
	oathUrl := "https://getpocket.com/v3/oauth/authorize"
	data := map[string]string {
		"consumer_key": consumerKey,
		"code": requestToken,
	}
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(data)
	resp, err := http.Post(oathUrl, "application/json", payloadBuffer)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	body := string(bodyBytes)

	q, err := url.ParseQuery(body)
	if err != nil {
		return "", "", err
	}

	username = q["username"][0]
	accessToken = q["access_token"][0]
	return username, accessToken, nil
}

func WriteFile(data AuthInfo, filename string) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func AskUser(requestToken string) {
	codesUrl := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
							requestToken, "https://getpocket.com/connected_applications")
	fmt.Printf("Please visit this page and sign-in to your pocket account: %s\n", codesUrl)
	fmt.Println("Once you have signed in there, hit <enter> to continue")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func AuthenticateCmd() {
	consumerKey, ok := os.LookupEnv("POCKET_CONSUMER_KEY")
	if !ok {
		log.Fatal("Consumer key was not found")
	}

	requestToken, err := RequestToken(consumerKey)
	if err != nil {
		log.Fatalf("An error occured while trying to fetch request token: %v", err)
	}

	AskUser(requestToken)

	username, accessToken, err := AutherizeUser(consumerKey, requestToken)
	if err != nil {
		log.Fatalf("An error occured while trying to authorize user: %v", err)
	}

	authInfo := AuthInfo{
		ConsumerKey: consumerKey,
		Username: username,
		AccessToken: accessToken,
	}
	err = WriteFile(authInfo, "auth.json")
	if err != nil {
		log.Fatalf("An error occured while trying to write json file: %v", err)
	}
}
