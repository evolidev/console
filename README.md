[![Main](https://github.com/evolidev/console/actions/workflows/test.yml/badge.svg)](https://github.com/evolidev/console/actions/workflows/test.yml)

# evolidev/console

Welcome to the "evolidev/console" repository!

Are you ready to take your command line interface game to the next level? Look no further than evolidev/console. With its intuitive and easy-to-use design, you can quickly and easily create custom CLI commands that will blow your users away.

Getting started is a breeze. Simply import the evolidev/console package and use the provided structs and functions to define your command's functionality. You can then register your command with the evolidev/console router and watch as it seamlessly integrates into your application's command line interface.

But evolidev/console isn't just about creating new commands. We've also included advanced features such as command auto-completion and grouping, allowing for a more organized and efficient command line experience for your users. And with built-in support for command flags and arguments, you can easily customize and configure your commands to suit your specific needs.

Don't settle for a bland and boring command line interface. Choose evolidev/console and elevate your CLI game to new heights!

## Examples

A very simple exmaple that demostrates how the console plugin can be used.
```go
package main

import (
    "github.com/evolidev/console"
    "github.com/evolidev/console/parse"

    "fmt"
)

func main() {

    cli := console.New()

    cli.AddCommand("simple", "Simple command that outputs the current time", func(c *parse.ParsedCommand) {
        fmt.Println("Hello World:")
    })

    cli.Run()
}

```