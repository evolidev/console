package console

import (
	"github.com/evolidev/console/parse"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
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

}
