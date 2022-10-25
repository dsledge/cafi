/*
Module: CAFI
Package: AWS
Description: Accounts structure and functions to read and parse account configuration json file.
*/
package aws

import (
	"encoding/json"
	"io/ioutil"
)

var (
	aws_accounts  Accounts
	accounts_file string
)

type Account struct {
	AccountNumber    string   `json:"account_number"`
	AccountName      string   `json:"account_name"`
	SupportedRegions []string `json:"supported_regions"`
	ProfileName      string   `json:"profile_name"`
	ProfileRegion    string   `json:"profile_region"`
	STSExternalId    *string  `json:"sts_external_id"`
	STSRoleArn       string   `json:"sts_role_arn"`
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

func loadAccounts(filePath string) error {
	accounts_file = filePath
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &aws_accounts)
	if err != nil {
		return err
	}

	return nil
}

func getAccount(acct string) *Account {
	for _, account := range aws_accounts.Accounts {
		if acct == account.AccountNumber {
			return &account
		}
	}

	return nil
}
