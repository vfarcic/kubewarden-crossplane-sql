package main

import (
	// "encoding/json"
	// "fmt"
	// "strings"

	// onelog "github.com/francoispqt/onelog"
	// corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	"encoding/json"
	"fmt"
	"strings"

	onelog "github.com/francoispqt/onelog"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	// kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

type Sql struct {
	APIVersion string             `json:"apiVersion,omitempty"`
	Kind       string             `json:"kind,omitempty"`
	Metadata   *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec       *SqlSpec           `json:"spec,omitempty"`
}

type SqlSpec struct {
	ID         string            `json:"id"`
	Parameters SqlSpecParameters `json:"parameters"`
}

type SqlSpecParameters struct {
	Version string `json:"version"`
	Size    string `json:"size"`
}

func validate(payload []byte) ([]byte, error) {
	// Create a ValidationRequest instance from the incoming payload
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Access the **raw** JSON that describes the object
	sqlJson := validationRequest.Request.Object

	// Try to create a SQL instance using the RAW JSON we got from the
	// ValidationRequest.
	sql := &Sql{}
	if err := json.Unmarshal([]byte(sqlJson), sql); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("Cannot decode SQL object: %s", err.Error())),
			kubewarden.Code(400))
	}

	logger.DebugWithFields("validating SQL object", func(e onelog.Entry) {
		e.String("name", sql.Metadata.Name)
		e.String("namespace", sql.Metadata.Namespace)
	})

	if !settings.IsSizeAllowed(sql.Spec.Parameters.Size) {
		logger.InfoWithFields("rejecting SQL object", func(e onelog.Entry) {
			e.String("name", sql.Metadata.Name)
			e.String("denied_names", strings.Join(settings.AllowedSizes, ","))
		})

		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("The '%s' name is on the deny list. The spec.parameters.size cannot be '%s'", sql.Metadata.Name, sql.Spec.Parameters.Size)),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
