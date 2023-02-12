package detection

import (
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

// Finds sensitive data from rekognition output
func ImageSensitiveData(imageText []types.TextDetection) []types.TextDetection {

	// create empty array of sensitive data
	var sensitiveData []types.TextDetection

	// iterate through array of text detections
	for _, item := range imageText {
		if SensitiveData(*item.DetectedText, true) {
			sensitiveData = append(sensitiveData, item)
		}
	}

	return sensitiveData
}
