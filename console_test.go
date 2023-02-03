package console

import (
	"fmt"
	"github.com/evolidev/console/color"
	"github.com/evolidev/console/parse"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseSimpleCommand(t *testing.T) {
	t.Parallel()
	t.Run("Parse simple command with required parameter", func(t *testing.T) {
		command := "mail:send foo"
		definition := "mail:send {user}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "foo", cmd.GetArgument("user").String())
	})

	t.Run("Parse simple command with optional parameter", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user?}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "", cmd.GetArgument("user").String())
	})

	t.Run("Parse simple command with optional parameter", func(t *testing.T) {
		command := "mail:send foo"
		definition := "mail:send {user?}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "foo", cmd.GetArgument("user").String())
	})

	t.Run("Parse simple command with optional parameter and default value", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user=foo}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "foo", cmd.GetArgument("user").String())
	})

	t.Run("Parse command and pass options", func(t *testing.T) {
		command := "mail:send foo --queue"
		definition := "mail:send {user} {--queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, true, cmd.GetOption("queue").Bool())
	})

	t.Run("Parse command and pass options", func(t *testing.T) {
		command := "serve"
		definition := "serve {--port}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, 1010, cmd.GetOptionWithDefault("port", 1010).Integer())
	})

	t.Run("Parse command and check argument with incorrect default value", func(t *testing.T) {
		command := "serve https"
		definition := "serve {secure} {--port}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "https", cmd.GetArgumentWithDefault("secure", "not-correct").String())
	})

	t.Run("Parse command and check argument with default value", func(t *testing.T) {
		command := "serve"
		definition := "serve {secure} {--port}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "correct", cmd.GetArgumentWithDefault("secure", "correct").String())
	})

	t.Run("Parse command and pass required option", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--queue=}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "", cmd.GetOption("queue").String())
	})

	t.Run("Parse command and pass option and alias", func(t *testing.T) {
		command := "mail:send -Q"
		definition := "mail:send {user} {--Q|queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, true, cmd.GetOption("Q").Bool())
	})

	t.Run("Get name of command", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--Q|queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "mail:send", cmd.GetName())
	})

	t.Run("Get prefix of command", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--Q|queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "mail", cmd.GetPrefix())
	})

	t.Run("Get subcommand of command", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--Q|queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "send", cmd.GetSubCommand())
	})

	t.Run("Get empty subcommand of command", func(t *testing.T) {
		command := "mail"
		definition := "mail:send {user} {--Q|queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "", cmd.GetSubCommand())
	})

	t.Run("Get default value of option", func(t *testing.T) {
		command := "mail"
		definition := "mail:send {user} {--Q|queue}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "default", cmd.GetOptionWithDefault("Q", "default").String())
	})

	t.Run("Get empty option", func(t *testing.T) {
		command := "mail"
		definition := "mail:send"

		cmd := parse.Parse(definition, command)

		assert.Nil(t, cmd.GetOption("Q"), "Option should be nil")
	})

	t.Run("Get empty argument", func(t *testing.T) {
		command := "mail"
		definition := "mail:send"

		cmd := parse.Parse(definition, command)

		assert.Nil(t, cmd.GetArgument("Q"), "Argument should be nil")
	})

	t.Run("Get argument with default argument", func(t *testing.T) {
		command := "mail"
		definition := "mail:send {Q=test}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "test", cmd.GetArgumentWithDefault("Q", "default").String())
	})

	t.Run("Get argument with default argument", func(t *testing.T) {
		command := "mail"
		definition := "mail:send {Q=test}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "test", cmd.GetArgumentWithDefault("Q", "").String())
	})

	t.Run("Get argument with default option", func(t *testing.T) {
		command := "mail"
		definition := "mail:send {--Q=test}"

		cmd := parse.Parse(definition, command)

		assert.Equal(t, "test", cmd.GetOptionWithDefault("Q", "default").String())
	})

	t.Run("Create simple command", func(t *testing.T) {
		cli := New()
		cli.AddCommand("mail:send {user}", "Send email", func(cmd *parse.ParsedCommand) {
			assert.Equal(t, "foo", cmd.GetArgument("user").String())
		})

		cli.Call([]string{"mail:send", "foo"})
	})

	t.Run("Create simple command with default value", func(t *testing.T) {
		cli := New()
		command := &Command{
			Definition:  "mail:send {user=test}",
			Description: "Send email",
			Execution: func(cmd *parse.ParsedCommand) {
				assert.Equal(t, "test", cmd.GetArgument("user").String())
			},
		}
		cli.Add(command)

		cli.Call([]string{"mail:send"})
	})

	t.Run("Run non existing command", func(t *testing.T) {
		cli := New()

		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.Call([]string{"mail:send", "foo"})

		err := w.Close()
		assert.Nil(t, err)
		out, _ := io.ReadAll(r)

		assert.True(t, strings.Contains(string(out), "Command not found"), "Expected 'Command not found' but got '%s'", string(out))
	})

	t.Run("Make sure command is rendered in stdout", func(t *testing.T) {
		cli := New()

		cli.AddCommand("mail:send {user}", "Send email", func(cmd *parse.ParsedCommand) {})

		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.Run()

		err := w.Close()
		assert.Nil(t, err)

		out, _ := io.ReadAll(r)

		assert.Contains(t, string(out), "mail:", "Expected 'mail:send' but got '%s'", string(out))
		assert.Contains(t, string(out), "Send email", "Expected 'Send email' but got '%s'", string(out))
	})

	t.Run("Create simple command with default value", func(t *testing.T) {
		command := &Command{
			Definition:  "mail:send {user=test}",
			Description: "Send email",
			Execution: func(cmd *parse.ParsedCommand) {
				assert.Equal(t, "test", cmd.GetArgument("user").String())
			},
		}

		assert.Equal(t, "send", command.GetCommand())
		assert.Equal(t, "Send email", command.GetDescription())
	})

	t.Run("Make sure command is rendered in stdout", func(t *testing.T) {
		cli := New()

		cli.AddCommand("mail:send {user}", "Send email", func(cmd *parse.ParsedCommand) {})
		cli.AddCommand("queue:run", "Run queue", func(cmd *parse.ParsedCommand) {})

		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.Run()

		err := w.Close()
		assert.Nil(t, err)

		out, _ := io.ReadAll(r)

		assert.Contains(t, string(out), "queue:", "Expected 'mail:send' but got '%s'", string(out))
	})

}

func TestText(t *testing.T) {
	tests := []struct {
		code    int
		value   interface{}
		want    string
		hasAnsi bool
	}{
		{0, "hello", "\u001b[38;5;0mhello\u001b[0m", true},
		{1, "world", "\u001b[38;5;1mworld\u001b[0m", true},
		{255, "!", "\u001b[38;5;255m!\u001b[0m", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d-%v", test.code, test.value), func(t *testing.T) {
			got := color.Text(test.code, test.value)
			if got != test.want {
				t.Errorf("Text(%d, %v) = %q, want %q", test.code, test.value, got, test.want)
			}
			if (len(got) > 4) != test.hasAnsi {
				t.Errorf("Text(%d, %v) has ANSI = %t, want %t", test.code, test.value, len(got) > 4, test.hasAnsi)
			}
		})
	}
}

func TestBg(t *testing.T) {
	tests := []struct {
		code    int
		value   interface{}
		want    string
		hasAnsi bool
	}{
		{0, "hello", "\u001b[48;5;0mhello\u001b[0m", true},
		{1, "world", "\u001b[48;5;1mworld\u001b[0m", true},
		{255, "!", "\u001b[48;5;255m!\u001b[0m", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d-%v", test.code, test.value), func(t *testing.T) {
			got := color.Bg(test.code, test.value)
			if got != test.want {
				t.Errorf("Bg(%d, %v) = %q, want %q", test.code, test.value, got, test.want)
			}
			if (len(got) > 4) != test.hasAnsi {
				t.Errorf("Bg(%d, %v) has ANSI = %t, want %t", test.code, test.value, len(got) > 4, test.hasAnsi)
			}
		})
	}
}

func TestGetOptionWithDefault(t *testing.T) {
	tests := []struct {
		name          string
		parsedCommand parse.ParsedCommand
		optionName    string
		defaultValue  interface{}
		expectedValue interface{}
	}{
		{
			name: "Option exists",
			parsedCommand: parse.ParsedCommand{
				Options: map[string]interface{}{
					"option1": "value1",
				},
			},
			optionName:    "option1",
			defaultValue:  "default1",
			expectedValue: "value1",
		},
		{
			name: "Option does not exist",
			parsedCommand: parse.ParsedCommand{
				Options: map[string]interface{}{},
			},
			optionName:    "option2",
			defaultValue:  "default2",
			expectedValue: "default2",
		},
		{
			name: "Option exists but is empty",
			parsedCommand: parse.ParsedCommand{
				Options: map[string]interface{}{
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
			got, gotAny := parse.ExtractField(test.item, test.prefix)
			if got != test.want {
				t.Errorf("ExtractField(%q, %q) = %q, want %q", test.item, test.prefix, got, test.want)
			}
			if !reflect.DeepEqual(gotAny, test.wantAny) {
				t.Errorf("ExtractField(%q, %q) = %v, want %v", test.item, test.prefix, gotAny, test.wantAny)
			}
		})
	}
}
