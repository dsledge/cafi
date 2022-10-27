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

### CAFI NewRandomToken utility function
This is a utility function to generate a new random token hex encoded. Here is an snippet of how to use the function. A full functional example can be found in the examples directory.
```go
byte_length := 16
token, err := cafi.NewRandomToken(byte_length)
if err != nil {
    fmt.Printf("Unable to generate a new random token: %s", err)
}
fmt.Printf("%d byte random token generated: %s", byte_length, token)
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
            "account_name": "<aws_profile_name>",
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

### CAFI AWS Provider Region Support
Here is an example of using the supported regions data from the accounts.json file. Setting the **input.Config.Region** setting before passing the **input.Config** to the new service will set the region to which the service should interact.
```go
for _, region := range input.Account.SupportedRegions {
            input.Config.Region = region
            svc := s3.NewFromConfig(*input.Config)
            // Do work with the service here
}
```

### CAFI AWS Provider Outputs
When creating the named function to execute in each account there are 2 parameters that are expected **func(input *aws_cafi.Input, output aws_cafi.Output)**. Output is a generic interface and can be fullfilled by creating any kind of structure required. This structure gets passed to the **ExecuteOnAccounts()** function and can be used inside the named function running against the account. This is useful when you need the results from all executions to be returned for further processing after gathering data from each account run.
```go
type TestOutput struct {
    Name string
}

output := TestOutput{}
err = aws_cafi.ExecuteOnAccounts(nil, s3listbuckets, &output)
if err != nil {
    fmt.Printf("Error iterating accounts: %s\n", err)
}

fmt.Printf("TESTING OUTPUT: %s\n", output.Name)
```