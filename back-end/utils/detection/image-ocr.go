package detection

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

func ImageOCR(file []byte) (*rekognition.DetectTextOutput, error) {
	// load aws configuations
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// configure rekognition
	rekognitionClient := rekognition.NewFromConfig(cfg)

	// call rekognition to detect text
	response, err := rekognitionClient.DetectText(context.TODO(), &rekognition.DetectTextInput{
		Image: &types.Image{
			Bytes: file,
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
