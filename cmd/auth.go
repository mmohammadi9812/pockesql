/*
* Copyright 2022 Mohammad Mohamamdi. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
 */

package cmd

import (
	"log"
	"os"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate your account to allow fetching entries",
	Run: func(_ *cobra.Command, _ []string) {
		AuthenticateCmd()
	},
}

func AuthenticateCmd() {
	consumerKey, ok := os.LookupEnv("POCKET_CONSUMER_KEY")
	if !ok {
		log.Fatal("Consumer key was not found")
	}

	requestToken, err := src.RequestToken(consumerKey)
	if err != nil {
		log.Fatalf("An error occured while trying to fetch request token: %v", err)
	}

	src.OpenConfirmUrl(requestToken)

	username, accessToken, err := src.AutherizeUser(consumerKey, requestToken)
	if err != nil {
		log.Fatalf("An error occured while trying to authorize user: %v", err)
	}

	authInfo := src.AuthInfo{
		ConsumerKey: consumerKey,
		Username:    username,
		AccessToken: accessToken,
	}
	err = src.WriteFile(authInfo, "auth.json")
	if err != nil {
		log.Fatalf("An error occured while trying to write json file: %v", err)
	}
}
