package main

import (
	"encoding/json"
	"testing"

	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
)

var sqlTestObject = Sql{
	APIVersion: "devopstoolkitseries.com/v1alpha1",
	Kind:       "Sql",
	Metadata: &metav1.ObjectMeta{
		Name:      "my-db",
		Namespace: "production",
	},
	Spec: &SqlSpec{
		ID: "my-db",
		Parameters: SqlSpecParameters{
			Version: "14",
			Size:    "medium",
		},
	},
}

func TestEmptySizeLeadsToApproval(t *testing.T) {
	settings := Settings{}
	sql := Sql{
		APIVersion: "devopstoolkitseries.com/v1alpha1",
		Kind:       "Sql",
		Metadata: &metav1.ObjectMeta{
			Name:      "my-db",
			Namespace: "production",
		},
		Spec: &SqlSpec{
			ID: "my-db",
			Parameters: SqlSpecParameters{
				Version: "14",
				Size:    "medium",
			},
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&sql, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Errorf("Unexpected rejection: msg %s - code %d", *response.Message, *response.Code)
	}
}

func TestApproval(t *testing.T) {
	settings := Settings{
		AllowedSizes: []string{"medium", "large"},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&sqlTestObject, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Error("Unexpected rejection")
	}
}

func TestApproveFixtureDenied(t *testing.T) {
	settings := Settings{
		AllowedSizes: []string{"medium", "large"},
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture(
		"test_data/sql.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted {
		t.Error("Unexpected rejection")
	}
}

func TestRejectionBecauseSizeIsDenied(t *testing.T) {
	settings := Settings{
		AllowedSizes: []string{"medium", "large"},
	}
	sql := sqlTestObject
	sql.Spec.Parameters.Size = "small"
	payload, err := kubewarden_testing.BuildValidationRequest(&sql, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != false {
		t.Error("Unexpected approval")
	}

	expected_message := "The 'my-db' name is on the deny list. The spec.parameters.size cannot be 'small'"

	if response.Message == nil {
		t.Errorf("expected response to have a message")
	}

	if *response.Message != expected_message {
		t.Errorf("Got '%s' instead of '%s'", *response.Message, expected_message)
	}
}
