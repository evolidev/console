package parse

import (
	"reflect"
	"testing"
)

func TestGetOptionWithDefault(t *testing.T) {
	tests := []struct {
		name          string
		parsedCommand ParsedCommand
		optionName    string
		defaultValue  interface{}
		expectedValue interface{}
	}{
		{
			name: "Option exists",
			parsedCommand: ParsedCommand{
				options: map[string]interface{}{
					"option1": "value1",
				},
			},
			optionName:    "option1",
			defaultValue:  "default1",
			expectedValue: "value1",
		},
		{
			name: "Option does not exist",
			parsedCommand: ParsedCommand{
				options: map[string]interface{}{},
			},
			optionName:    "option2",
			defaultValue:  "default2",
			expectedValue: "default2",
		},
		{
			name: "Option exists but is empty",
			parsedCommand: ParsedCommand{
				options: map[string]interface{}{
					"option3": "",
				},
			},
			optionName:    "option3",
			defaultValue:  "default3",
			expectedValue: "default3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.parsedCommand.GetOptionWithDefault(test.optionName, test.defaultValue)
			if value.Value != test.expectedValue {
				t.Errorf("Expected value %v, but got %v", test.expectedValue, value.Value)
			}
		})
	}
}

func TestExtractField(t *testing.T) {
	tests := []struct {
		item    string
		prefix  string
		want    string
		wantAny interface{}
	}{
		{"field=value", "--", "field", "value"},
		{"field", "-", "field", true},
		{"field=", "--", "field", true},
		{"field=", "--", "field", true},
		{"", "", "", true},
	}

	for _, test := range tests {
		t.Run(test.item, func(t *testing.T) {
			got, gotAny := extractField(test.item, test.prefix)
			if got != test.want {
				t.Errorf("extractField(%q, %q) = %q, want %q", test.item, test.prefix, got, test.want)
			}
			if !reflect.DeepEqual(gotAny, test.wantAny) {
				t.Errorf("extractField(%q, %q) = %v, want %v", test.item, test.prefix, gotAny, test.wantAny)
			}
		})
	}
}
