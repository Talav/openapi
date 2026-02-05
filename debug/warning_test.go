package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWarning(t *testing.T) {
	warning := NewWarning(WarnDegradationWebhooks, "#/webhooks", "webhooks are 3.1-only; dropped")

	assert.Equal(t, WarnDegradationWebhooks, warning.Code())
	assert.Equal(t, "#/webhooks", warning.Path())
	assert.Equal(t, "webhooks are 3.1-only; dropped", warning.Message())
	assert.Contains(t, warning.String(), string(WarnDegradationWebhooks))
	assert.Contains(t, warning.String(), "webhooks are 3.1-only; dropped")
}

func TestWarningString(t *testing.T) {
	warning := NewWarning(WarnDegradationInfoSummary, "#/info/summary", "info.summary is 3.1-only")

	str := warning.String()
	assert.Contains(t, str, "[DEGRADATION_INFO_SUMMARY]")
	assert.Contains(t, str, "info.summary is 3.1-only")
}

func TestWarningsHas(t *testing.T) {
	warnings := Warnings{
		NewWarning(WarnDegradationWebhooks, "#/webhooks", "test"),
		NewWarning(WarnDegradationInfoSummary, "#/info/summary", "test"),
	}

	assert.True(t, warnings.Has(WarnDegradationWebhooks))
	assert.True(t, warnings.Has(WarnDegradationInfoSummary))
	assert.False(t, warnings.Has(WarnDegradationMutualTLS))
}

func TestWarningsHas_EmptyList(t *testing.T) {
	var warnings Warnings

	assert.False(t, warnings.Has(WarnDegradationWebhooks))
}

func TestWarningsHas_NilList(t *testing.T) {
	var warnings Warnings = nil

	assert.False(t, warnings.Has(WarnDegradationWebhooks))
}

func TestWarningsAppend(t *testing.T) {
	var warnings Warnings

	warnings.Append(NewWarning(WarnDegradationWebhooks, "#/webhooks", "test1"))
	assert.Len(t, warnings, 1)
	assert.True(t, warnings.Has(WarnDegradationWebhooks))

	warnings.Append(NewWarning(WarnDegradationInfoSummary, "#/info/summary", "test2"))
	assert.Len(t, warnings, 2)
	assert.True(t, warnings.Has(WarnDegradationInfoSummary))
}

func TestWarningsAppend_Multiple(t *testing.T) {
	var warnings Warnings

	warnings.Append(NewWarning(WarnDegradationWebhooks, "#/webhooks", "msg1"))
	warnings.Append(NewWarning(WarnDegradationInfoSummary, "#/info", "msg2"))
	warnings.Append(NewWarning(WarnDegradationWebhooks, "#/webhooks2", "msg3"))

	assert.Len(t, warnings, 3)
	assert.True(t, warnings.Has(WarnDegradationWebhooks))
	assert.True(t, warnings.Has(WarnDegradationInfoSummary))
}

func TestWarningCodes(t *testing.T) {
	codes := []WarningCode{
		WarnDegradationWebhooks,
		WarnDegradationInfoSummary,
		WarnDegradationLicenseIdentifier,
		WarnDegradationMutualTLS,
		WarnDegradationConstToEnum,
		WarnDegradationConstToEnumConflict,
		WarnDegradationPathItems,
		WarnDegradationPatternProperties,
		WarnDegradationUnevaluatedProperties,
		WarnDegradationContentEncoding,
		WarnDegradationContentMediaType,
		WarnDegradationMultipleExamples,
		WarnInvalidExampleMutualExclusivity,
	}

	for _, code := range codes {
		t.Run(string(code), func(t *testing.T) {
			assert.NotEmpty(t, code.String())
			assert.Equal(t, string(code), code.String())
		})
	}
}

func TestWarningCodeString(t *testing.T) {
	code := WarnDegradationWebhooks
	assert.Equal(t, "DEGRADATION_WEBHOOKS", code.String())
}

func TestWarningInterface(t *testing.T) {
	_ = NewWarning(WarnDegradationWebhooks, "#/test", "test message")
}

func TestWarningsCollection(t *testing.T) {
	warnings := make(Warnings, 0)

	warnings.Append(NewWarning(WarnDegradationWebhooks, "#/webhooks", "msg1"))
	warnings.Append(NewWarning(WarnDegradationInfoSummary, "#/info", "msg2"))

	assert.Len(t, warnings, 2)

	// Check individual warnings
	assert.Equal(t, WarnDegradationWebhooks, warnings[0].Code())
	assert.Equal(t, "#/webhooks", warnings[0].Path())
	assert.Equal(t, "msg1", warnings[0].Message())

	assert.Equal(t, WarnDegradationInfoSummary, warnings[1].Code())
	assert.Equal(t, "#/info", warnings[1].Path())
	assert.Equal(t, "msg2", warnings[1].Message())
}
