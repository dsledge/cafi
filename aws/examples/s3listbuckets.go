package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	cafi "github.com/dsledge/cafi"
	aws_cafi "github.com/dsledge/cafi/aws"
	scribble "github.com/dsledge/scribble"
)

var (
	logfile  = flag.String("logfile", "console", "The default log file will log to system console")
	loglevel = flag.Int("loglevel", 0, "Sets the default log level to INFO messages and higher")
)

type TestOutput struct {
	Name string
}

func main() {
	flag.Parse()

	// Configure the CAFI SDK
	cafi.Configure(logfile, loglevel)

	// Configure the AWS Provider
	err := aws_cafi.Configure("accounts.json")
	if err != nil {
		scribble.Fatal("Error: %s", err)
	}

	// Function to run through the account iterator
	s3listbuckets := func(input *aws_cafi.Input, output aws_cafi.Output) {
		// Iterate over the region list
		for _, region := range input.Account.SupportedRegions {
			input.Config.Region = region
			svc := s3.NewFromConfig(*input.Config)

			// Call the list buckets api operation
			listBucketsResult, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
			if err != nil {
				panic("Couldn't list buckets")
			}

			// Iterate over the results and print out the buckets
			for _, bucket := range listBucketsResult.Buckets {
				fmt.Printf("Account: %s \tRegion: %s \tBucket: %s\n", input.Account.AccountName, region, *bucket.Name)
			}

			// Assert the object type from the interface to the known object type
			testoutput := output.(*TestOutput)
			testoutput.Name = "Test"
		}
	}

	output := TestOutput{}
	err = aws_cafi.ExecuteOnAccounts(nil, s3listbuckets, &output)
	if err != nil {
		fmt.Printf("Error iterating accounts: %s\n", err)
	}

	fmt.Printf("TESTING OUTPUT: %s\n", output.Name)
}
