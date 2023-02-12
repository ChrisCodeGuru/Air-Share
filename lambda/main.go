package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func main() {
	lambda.Start(KeyRotate)
}

func KeyRotate() {
	// Create new session with kms and configure request params
	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	svc := kms.New(session)
	createRSA4096KeyParams := &kms.CreateKeyInput{
		KeySpec:  aws.String("RSA_4096"),
		KeyUsage: aws.String("ENCRYPT_DECRYPT"),
	}

	createKeyResult, err := svc.CreateKey(createRSA4096KeyParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(kms.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case kms.ErrCodeDependencyTimeoutException:
				fmt.Println(kms.ErrCodeDependencyTimeoutException, aerr.Error())
			case kms.ErrCodeInvalidArnException:
				fmt.Println(kms.ErrCodeInvalidArnException, aerr.Error())
			case kms.ErrCodeUnsupportedOperationException:
				fmt.Println(kms.ErrCodeUnsupportedOperationException, aerr.Error())
			case kms.ErrCodeInternalException:
				fmt.Println(kms.ErrCodeInternalException, aerr.Error())
			case kms.ErrCodeLimitExceededException:
				fmt.Println(kms.ErrCodeLimitExceededException, aerr.Error())
			case kms.ErrCodeTagException:
				fmt.Println(kms.ErrCodeTagException, aerr.Error())
			case kms.ErrCodeCustomKeyStoreNotFoundException:
				fmt.Println(kms.ErrCodeCustomKeyStoreNotFoundException, aerr.Error())
			case kms.ErrCodeCustomKeyStoreInvalidStateException:
				fmt.Println(kms.ErrCodeCustomKeyStoreInvalidStateException, aerr.Error())
			case kms.ErrCodeCloudHsmClusterInvalidConfigurationException:
				fmt.Println(kms.ErrCodeCloudHsmClusterInvalidConfigurationException, aerr.Error())
			case kms.ErrCodeXksKeyInvalidConfigurationException:
				fmt.Println(kms.ErrCodeXksKeyInvalidConfigurationException, aerr.Error())
			case kms.ErrCodeXksKeyAlreadyInUseException:
				fmt.Println(kms.ErrCodeXksKeyAlreadyInUseException, aerr.Error())
			case kms.ErrCodeXksKeyNotFoundException:
				fmt.Println(kms.ErrCodeXksKeyNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	fmt.Println(createKeyResult)

	// Update key aliases
	updateKeyAliasParams := &kms.UpdateAliasInput{
		AliasName:   aws.String("alias/ISPJ"),
		TargetKeyId: aws.String(*createKeyResult.KeyMetadata.Arn),
	}

	updateKeyAliasResult, err := svc.UpdateAlias(updateKeyAliasParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeDependencyTimeoutException:
				fmt.Println(kms.ErrCodeDependencyTimeoutException, aerr.Error())
			case kms.ErrCodeNotFoundException:
				fmt.Println(kms.ErrCodeNotFoundException, aerr.Error())
			case kms.ErrCodeInternalException:
				fmt.Println(kms.ErrCodeInternalException, aerr.Error())
			case kms.ErrCodeLimitExceededException:
				fmt.Println(kms.ErrCodeLimitExceededException, aerr.Error())
			case kms.ErrCodeInvalidStateException:
				fmt.Println(kms.ErrCodeInvalidStateException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(updateKeyAliasResult)
}
