package console

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/evolidev/console/color"
	"github.com/evolidev/console/parse"
	"github.com/olekukonko/tablewriter"
)

type Command struct {
	Definition  string
	Description string
	Execution   func(c *parse.ParsedCommand)
}

func (cmd *Command) GetName() string {
	parts := strings.Split(cmd.Definition, " ")
	return parts[0]
}

func (cmd *Command) GetCommand() string {
	name := cmd.GetName()

	parts := strings.Split(name, ":")
	if len(parts) > 1 {
		name = parts[len(parts)-1]
	}

	return name
}

func (cmd *Command) GetDescription() string {
	return cmd.Description
}

type CommandGroup struct {
	Name        string
	Description string
	Prefix      string
	Commands    []*Command
}

type Console struct {
	Commands map[string]*Command
	Coloring bool
	Output   io.Writer
	Title    string
}

func (c *Console) Run() {
	args := os.Args[1:]

	c.Call(args)
}

func (c *Console) Call(args []string) {
	if len(args) > 0 {
		command := args[0]
		if cmd, ok := c.Commands[command]; ok {
			parsed := parse.Parse(cmd.Definition, strings.Join(args, " "))
			cmd.Execution(parsed)
			return
		} else {
			fmt.Println()
			fmt.Println(c.Bg(210, fmt.Sprintf("%46s", " ")))
			fmt.Println(
				c.Bg(210, c.Text(255, fmt.Sprintf("%s", "    Sorry, but the command does not exist.    "))),
			)
			fmt.Println(c.Bg(210, fmt.Sprintf("%46s", " ")))
			fmt.Println()
		}
	}

	c.Render()
}

func (c *Console) Add(command *Command) {
	c.Commands[command.GetName()] = command
}

func (c *Console) AddCommand(name string, description string, execution func(c *parse.ParsedCommand)) *Command {
	command := &Command{name, description, execution}
	c.Add(command)

	return command
}

func New() *Console {
	return &Console{
		Commands: make(map[string]*Command),
		Coloring: true,
		Output:   os.Stdout,
	}
}

func groupCommands(commands map[string]*Command) []CommandGroup {
	groups := make(map[string][]*Command)

	var groupKeys []string
	for _, cmd := range commands {
		commandParts := strings.Split(cmd.GetName(), ":")
		prefix := ""
		if len(commandParts) > 1 {
			prefix = commandParts[0]
		}

		if _, ok := groups[prefix]; !ok {
			groupKeys = append(groupKeys, prefix)
		}

		groups[prefix] = append(groups[prefix], cmd)
	}

	sort.Strings(groupKeys)

	var groupedCommands []CommandGroup
	for _, key := range groupKeys {
		groupItems := groups[key]

		sort.Slice(groupItems, func(i, j int) bool {
			return groupItems[i].GetCommand() < groupItems[j].GetCommand()
		})

		groupedCommands = append(groupedCommands, CommandGroup{
			Name:        key,
			Description: "",
			Prefix:      key,
			Commands:    groupItems,
		})
	}

	return groupedCommands
}

func (c *Console) Render() {
	table := c.SetupTable()

	c.AddCommandsToTable(table)

	if c.Title != "" {
		fmt.Println()
		fmt.Println(c.Title)
		fmt.Println()
	}

	fmt.Println(c.Text(249, "USAGE:"))
	fmt.Printf("   command [options] [arguments]")
	fmt.Println()
	fmt.Println()

	table.Render()
}

func (c *Console) AddCommandsToTable(table *tablewriter.Table) {
	groupedCommands := groupCommands(c.Commands)
	for _, group := range groupedCommands {
		prefix := ""
		if group.Name != "" {
			table.Rich([]string{group.Name, group.Description}, []tablewriter.Colors{
				{tablewriter.FgHiGreenColor},
				{},
			})

			prefix = c.Text(140, group.Prefix+":")
		}

		for _, cmd := range group.Commands {
			table.Append([]string{
				prefix + c.Text(169, "  "+cmd.GetCommand()),
				c.Text(245, cmd.Description),
			})
		}

		table.Append([]string{""})
	}
}

func (c *Console) SetupTable() *tablewriter.Table {
	table := tablewriter.NewWriter(c.Output)
	table.SetHeader([]string{"AVAILABLE COMMANDS", ""})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	if c.Coloring {
		table.SetColumnColor(
			tablewriter.Colors{tablewriter.FgHiMagentaColor},
			tablewriter.Colors{tablewriter.FgHiBlackColor},
		)

		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
			tablewriter.Colors{tablewriter.FgHiBlackColor},
		)
	}

	return table
}

func (c *Console) DisableColors() {
	c.Coloring = false
}

func (c *Console) EnableColors() {
	c.Coloring = true
}

func (c *Console) SetOutput(output io.Writer) {
	c.Output = output
}

func (c *Console) Text(code int, value interface{}) string {
	if c.Coloring {
		return color.Text(code, value)
	}

	return fmt.Sprintf("%v", value)
}

func (c *Console) Bg(code int, value interface{}) string {
	if c.Coloring {
		return color.Bg(code, value)
	}

	return fmt.Sprintf("%v", value)
}

func (c *Console) SetTitle(title string) {
	c.Title = title
}
