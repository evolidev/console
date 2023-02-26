package main

import (
	"github.com/evolidev/console"
	"github.com/evolidev/console/color"
	"github.com/evolidev/console/parse"

	"fmt"
)

func main() {

	cli := console.New()

	cli.AddCommand("simple", "Simple command that outputs the current time", func(c *parse.ParsedCommand) {
		fmt.Println("Hello World:")
	})

	cli.SetTitle(fmt.Sprintf(
		"Simple Console %s", color.Text(169, "0.0.1"),
	))

	cli.Run()
}
