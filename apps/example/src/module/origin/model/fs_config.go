package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Enum type for FS
type FS string

// Constants for valid FS values
const (
	S3   FS = "s3"
	FTP  FS = "ftp"
	HTTP FS = "http"
	OS   FS = "os"
)

// Struct to hold the JSON config structure
type Config struct {
	FS         FS          `json:"fs"`
	S3Config   *S3Config   `json:"s3Config"`
	FTPConfig  *FTPConfig  `json:"ftpConfig"`
	HTTPConfig *HTTPConfig `json:"httpConfig"`
	OSConfig   *OSConfig   `json:"osConfig"`
}

// Structs for specific configurations
type S3Config struct {
	Bucket          string `json:"bucket"`
	Path            string `json:"path"`
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	UseSSL          bool   `json:"useSSL"`
	ForcePathStyle  bool
}

type FTPConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Path     string `json:"path"`
}

type HTTPConfig struct {
	BaseURL string `json:"baseURL"`
}

type OSConfig struct {
	BasePath string `json:"basePath"`
}

// func ParseJSONConfig(configJSON string) (config Config, err error) {
// 	err = json.Unmarshal([]byte(configJSON), &config)
// 	if err != nil {
// 		return Config{}, fmt.Errorf("failed to parse JSON config: %v", err)
// 	}
// 	if !isValidFS(FS(config.FS)) {
// 		return Config{}, fmt.Errorf("invalid FS value: %s", config.FS)
// 	}
// 	return config, nil
// }

// Function to validate FS value
func isValidFS(fs FS) bool {
	switch fs {
	case S3, FTP, HTTP, OS:
		return true
	default:
		return false
	}
}

func ParseBase64Config(base64Config string) (config Config, err error) {
	// Decode base64
	configJSON, err := base64.StdEncoding.DecodeString(base64Config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode base64 config: %v", err)
	}

	// Unmarshal JSON
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse JSON config: %v", err)
	}

	if !isValidFS(FS(config.FS)) {
		return Config{}, fmt.Errorf("invalid FS value: %s", config.FS)
	}

	return config, nil
}
