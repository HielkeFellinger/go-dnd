package helpers

import (
	"encoding/base64"
	"errors"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func SaveImageToCampaign(image FileUpload, campaignId uint, rawFileName string) (string, error) {
	baseLocation := os.Getenv("CAMPAIGN_DATA_DIR") + "/" + strconv.Itoa(int(campaignId))

	// Ensure Campaign dir exists
	if _, err := os.Stat(baseLocation + "/images"); os.IsNotExist(err) {
		if err := os.MkdirAll(baseLocation+"/images", os.ModePerm); err != nil {
			return "", err
		}
	}

	// Clean filename
	reg, _ := regexp.Compile("\\s+")
	strippedFileName := reg.ReplaceAllString(rawFileName, "")

	// Decode base64 content
	newImageContent, err := base64.StdEncoding.DecodeString(image.FileBase64)
	if err != nil {
		return "", err
	}

	// Check mimetype and get extension
	mimeType := http.DetectContentType(newImageContent)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", errors.New("not an image")
	}
	fileExtensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}
	if len(fileExtensions) == 0 {
		return "", errors.New("not an image")
	}

	// Attempt to write to storage; do not override
	newImageFileName := baseLocation + "/images/" + strippedFileName + fileExtensions[0]
	if _, err := os.Stat(newImageFileName); os.IsNotExist(err) {
		if err := os.WriteFile(newImageFileName, newImageContent, 0644); err != nil {
			return "", err
		}
		return newImageFileName, nil
	} else {
		return "", errors.New("file already exists")
	}
}

type FileUpload struct {
	Name       string `json:"Name"`
	Size       int32  `json:"Size"`
	Type       string `json:"Type"`
	FileBase64 string `json:"FileBase64"`
}
