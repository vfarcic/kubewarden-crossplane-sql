package main

import (
	"encoding/json"
	"fmt"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// Settings is the structure that describes the policy settings.
type Settings struct {
	AllowedSizes []string `json:"allowed_sizes"`
}

// No special checks have to be done
func (s *Settings) Valid() (bool, error) {
	return true, nil
}

func (s *Settings) IsSizeAllowed(size string) bool {
	if len(s.AllowedSizes) == 0 {
		return true
	}
	for _, allowedSize := range s.AllowedSizes {
		if allowedSize == size {
			return true
		}
	}
	return false
}

func NewSettingsFromValidationReq(validationReq *kubewarden_protocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	return settings, err
}

func validateSettings(payload []byte) ([]byte, error) {
	logger.Info("validating settings")

	settings := Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}

	valid, err := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	logger.Warn("rejecting settings")
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}
