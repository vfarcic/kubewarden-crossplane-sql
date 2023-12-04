package main

import (
	"encoding/json"
	"testing"
)

// TODO: Rewrite
func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := json.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if len(settings.AllowedSizes) != 0 {
		t.Errorf("Expected DeniedNames to be empty")
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

func TestIsSizeDenied(t *testing.T) {
	settings := Settings{
		AllowedSizes: []string{"medium", "large"},
	}
	if settings.IsSizeAllowed("small") {
		t.Errorf("size should be denied")
	}
	if !settings.IsSizeAllowed("medium") {
		t.Errorf("size should not be allowed")
	}
	if !settings.IsSizeAllowed("large") {
		t.Errorf("size should not be allowed")
	}
}
