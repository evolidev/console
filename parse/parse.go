package parse

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

var parseRegex = `[\-\\\/]{0,2}?((\w+)(?:[=:](\"[^\"]+\"|[^\s\"]+))?|((\/|\.\.?|)[^-\s][^"\s]*))(?:\s+|$)`

type ParsedCommand struct {
	Arguments  map[string]any
	Options    map[string]any
	Command    string
	Name       string
	SubCommand string
	Prefix     string
}

type Value struct {
	fmt.Stringer
	Value any
}

func (o *Value) Bool() bool {
	return cast.ToBool(o.Value)
}

func (o *Value) Integer() int {
	return cast.ToInt(o.Value)
}

func (o *Value) String() string {
	return cast.ToString(o.Value)
}

func (p *ParsedCommand) HasOption(name string) bool {
	if cmd, ok := p.Options[name]; ok {
		return cmd != nil
	}

	return false
}

func (p *ParsedCommand) HasArgument(name string) bool {
	if cmd, ok := p.Arguments[name]; ok {
		return cmd != nil
	}

	return false
}

func (p *ParsedCommand) GetArgument(name string) *Value {
	if !p.HasArgument(name) {
		return nil
	}

	argumentValue := p.Arguments[name]

	return &Value{Value: argumentValue}
}

func (p *ParsedCommand) GetOption(name string) *Value {
	if !p.HasOption(name) {
		return nil
	}
	optionValue := p.Options[name]

	return &Value{Value: optionValue}
}

func (p *ParsedCommand) GetName() string {
	return p.Name
}

func (p *ParsedCommand) GetPrefix() string {
	return p.Prefix
}

func (p *ParsedCommand) GetSubCommand() string {
	return p.SubCommand
}

func (p *ParsedCommand) GetOptionWithDefault(name string, defaultValue any) *Value {
	optionValue := p.Options[name]
	if optionValue == nil || optionValue == "" {
		return &Value{Value: defaultValue}
	}
	return &Value{Value: optionValue}
}

func (p *ParsedCommand) GetArgumentWithDefault(name string, defaultValue any) *Value {
	argumentValue := p.Arguments[name]
	if argumentValue == nil || argumentValue == "" {
		return &Value{Value: defaultValue}
	}
	return &Value{Value: argumentValue}
}

func Parse(definition string, command string) *ParsedCommand {
	arguments := make(map[string]any)
	options := make(map[string]any)

	var argumentsMap []string

	// parse definition
	argumentsMap = parseDefinition(definition, options, arguments, argumentsMap)

	items := parseCommand(command, options, argumentsMap, arguments)

	name := items[0]
	// split name into prefix and subcommand
	nameParts := strings.Split(name, ":")
	prefix := nameParts[0]
	subCommand := ""
	if len(nameParts) > 1 {
		subCommand = nameParts[1]
	}

	return &ParsedCommand{
		Arguments:  arguments,
		Options:    options,
		Command:    command,
		Name:       name,
		SubCommand: subCommand,
		Prefix:     prefix,
	}
}

func parseDefinition(definition string, options map[string]any, arguments map[string]any, argumentsMap []string) []string {
	definitionItems := strings.Split(definition, " ")
	for i := range definitionItems {
		definitionItems[i] = strings.TrimSpace(definitionItems[i])

		// parse definition item
		definitionItem := definitionItems[i]
		// remove curly bracket at the beginning and end of definition item
		definitionItem = strings.Trim(definitionItem, "{}?")

		// split definition item into name and Value
		definitionItemParts := strings.Split(definitionItem, "=")
		definitionItemName := definitionItemParts[0]
		definitionItemValue := ""
		if len(definitionItemParts) > 1 {
			definitionItemValue = definitionItemParts[1]
		}

		if strings.HasPrefix(definitionItemName, "--") {
			optionName := strings.TrimPrefix(definitionItemName, "--")

			optionNameParts := strings.Split(optionName, "|")
			for index := range optionNameParts {
				optionName = strings.TrimSpace(optionNameParts[index])
				options[optionName] = definitionItemValue
			}
		} else {
			arguments[definitionItemName] = definitionItemValue
			argumentsMap = append(argumentsMap, definitionItemName)
		}
	}
	return argumentsMap
}

func parseCommand(command string, options map[string]any, argumentsMap []string, arguments map[string]any) []string {
	// extract all arguments and options
	r, _ := regexp.Compile(parseRegex)
	items := r.FindAllString(command, -1)
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	for index, item := range items {
		if strings.HasPrefix(item, "--") {
			optionName, optionValue := ExtractField(item, "--")
			options[optionName] = optionValue
		} else if strings.HasPrefix(item, "-") {
			optionName, optionValue := ExtractField(item, "-")
			options[optionName] = optionValue
		} else {
			if index > 0 && index < len(argumentsMap) {
				arguments[argumentsMap[index]] = item
			}
		}
	}
	return items
}

func ExtractField(item string, prefix string) (string, any) {
	option := strings.TrimPrefix(item, prefix)
	// extract option name and Value
	parts := strings.Split(option, "=")
	optionName := parts[0]
	var optionValue any
	if len(parts) > 1 && parts[1] != "" {
		optionValue = parts[1]
	} else {
		optionValue = true
	}
	return optionName, optionValue
}
