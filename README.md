# CAFI
Cross Account Function Iterator to run a go function across multiple cloud accounts.

The CAFI package provides a simple way to work various cloud providers allowing for a function to be executed across multiple accounts / projects.

### CAFI Module Initialization
Configure the CAFI Module for logging settings
```go
    // Import the CAFI module
	import cafi "github.com/dsledge/cafi"

	// Configure the CAFI SDK
	cafi.Configure(logfile, loglevel)
```    

## AWS Provider
The CAFI AWS package provides a simple way to work with the Amazon AWS API across multiple accounts using STS authentication. It starts with simple functions to work with your shared credentials to begin the process of STS access. It also provides a structure to manage a list of accounts with which to iteract with.

### Example of accounts.json file
This is an example of an account.json file and the fields it requires, This file is used to know how to connect to AWS accounts using AWS STS access.
```json
{
    "accounts": [
        {
            "account_number": "<aws_account_number>",
            "account_name": "<aws_profile_name>",,
            "supported_regions": ["<aws_region_1>","<aws_region_2>","<aws_region_3>"],
            "profile_name": "<aws_profile_name>",
            "profile_region": "<aws_region>",
            "sts_external_id": "<aws_external_id>",
            "sts_role_arn": "arn:aws:iam::%s:role/<sts_role_to_assume>"
        }
    ]
}
```

### CAFI AWS Provider Initialization
Configure the CAFI Module for logging settings
```go
    // Import the CAFI module
	import aws_cafi "github.com/dsledge/cafi/aws"

    // Configure the AWS Provider
	err := aws_cafi.Configure("accounts.json")
	if err != nil {
		scribble.Fatal("Error: %s", err)
	}
```
