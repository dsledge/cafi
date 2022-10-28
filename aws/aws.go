/*
Module: CAFI
Package: AWS
Description: AWS Provider for the CAFI Module. This package provides support for the AWS cloud platform and is used to assume roles into AWS accounts and provide a mechanism
to execute functions against those accounts.
*/
package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/dsledge/scribble"
)

var (
	sts_clients map[string]*STSClient
	err         error
)

type STSClient struct {
	Service *sts.Client
	Config  *aws.Config
	Session string
}

// Input for the function passed to the execute methods
type Input struct {
	Config  *aws.Config
	Account *Account
}

// Generic interface to allow for any struct to be passed into the supported functions
type Output interface{}

// Initialize once
func init() {
	sts_clients = make(map[string]*STSClient)
}

// Create STS Client using shared credentials profile.
func getSTSService(profile_name string, profile_region string) (*sts.Client, error) {
	if sts_client, ok := sts_clients[profile_name]; ok {
		scribble.Trace("STS Client already exist for profile_name: %s, returning it to the caller", profile_name)
		return sts_client.Service, nil
	} else {
		scribble.Trace("STS Client does not exist for profile_name: %s, creating a new one", profile_name)
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile_name), config.WithRegion(profile_region))
		if err != nil {
			return nil, err
		}

		sts_clients[profile_name] = &STSClient{}
		sts_clients[profile_name].Service = sts.NewFromConfig(cfg)
		sts_clients[profile_name].Config = &cfg
		printCallerIdentity(sts_clients[profile_name])

		return sts_clients[profile_name].Service, nil
	}
}

// Print out the AWS caller identity
func printCallerIdentity(client *STSClient) {
	result, err := client.Service.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		scribble.Warn("Error getting CallerIdentity: %s", err)
	}

	// Set the session_name used via the shared credentials
	if len(client.Session) < 1 {
		arn_split := strings.Split(*result.Arn, "/")
		client.Session = arn_split[len(arn_split)-1]
	}

	scribble.Debug("STS Identity {Account: \"%s\", ARN: \"%s\", UserId: \"%s\"}", *result.Account, *result.Arn, *result.UserId)
}

// AssumeRole and get sts credentials for the remote account
func assumeRole(account *Account) error {
	service, err := getSTSService(account.ProfileName, account.ProfileRegion)
	if err != nil {
		return err
	}

	scribble.Trace("Acquiring credentials for the assumed role")
	creds := stscreds.NewAssumeRoleProvider(service, fmt.Sprintf(account.STSRoleArn, account.AccountNumber), func(options *stscreds.AssumeRoleOptions) {
		options.RoleSessionName = sts_clients[account.ProfileName].Session
		options.ExternalID = aws.String(*account.STSExternalId)
	})
	sts_clients[account.ProfileName].Config.Credentials = aws.NewCredentialsCache(creds)

	svc := sts.NewFromConfig(*sts_clients[account.ProfileName].Config)
	result, err := svc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	scribble.Debug("STS Identity {Account: \"%s\", ARN: \"%s\", UserId: \"%s\"}", *result.Account, *result.Arn, *result.UserId)
	scribble.Info("Switching to AWS Account: %s (%s)", account.AccountNumber, account.AccountName)
	return nil
}

// Configure the CASTS SDK
func Configure(filepath string) error {
	// Load the accounts json file
	err := loadAccounts(filepath)
	if err != nil {
		return err
	}

	return nil
}

// Execute function logic on multiple accounts
func ExecuteOnAccounts(accounts *[]string, f func(input *Input, out Output), out Output) {
	if accounts == nil {
		scribble.Trace("Iterating all accounts found in the %s", accounts_file)
		for _, account := range aws_accounts.Accounts {
			// AssumeRole and get sts credentials for the remote account
			err = assumeRole(&account)
			if err != nil {
				scribble.Error("Could not assume role for account: %s (%s)\n\nError Message: %s\n", account.AccountNumber, account.AccountName, err)
			}

			// Execute the function used for iteration
			input := Input{sts_clients[account.ProfileName].Config, &account}
			f(&input, out)
		}
	} else {
		scribble.Trace("Iterating accounts provided to ExecuteOnAccounts function: %s", *accounts)
		for _, account_number := range *accounts {
			// AssumeRole and get sts credentials for the remote account
			account := getAccount(account_number)
			err = assumeRole(account)
			if err != nil {
				scribble.Error("Could not assume role for account: %s (%s)\n\nError Message: %s\n", account.AccountNumber, account.AccountName, err)
			}

			// Execute the function used for iteration
			input := Input{sts_clients[account.ProfileName].Config, account}
			f(&input, out)
		}
	}
}
