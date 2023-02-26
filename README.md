[![Main](https://github.com/evolidev/console/actions/workflows/test.yml/badge.svg)](https://github.com/evolidev/console/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/evolidev/console/branch/dev/graph/badge.svg?token=6AGWZTOSKO)](https://codecov.io/github/evolidev/console)

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

## Functions

### AddCommand

The `AddCommand` method allows you to register a new command with the CLI router. The method takes in three arguments:

1. "name": This is the name of the command that will be used to invoke the command in the command line interface.
2. "description": This is a brief description of the command's functionality. This description will be displayed to the user when they use the --help flag with the command.
3. func(c *parse.ParsedCommand): This is the function that will be executed when the command is invoked. The function takes in a pointer to a ParsedCommand struct, which contains information about the command, such as its name, arguments, and flags.

For example, to add a command named "greet" with a description "prints a greeting message" that takes a name argument and prints "Hello, [name]" when invoked, you can use the following code:

```go
command := cli.AddCommand("greet {name}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetArgument("name").String()
    fmt.Println("Hello, ", name)
})
```

The AddCommand method returns the command struct and you can use it to add arguments, flags and more. Keep in mind that once you have added your commands, you need to run the router to listen for commands.

### Coloring

Using the `EnableColoring` method, you can enable colored output for your CLI. This is useful for adding color to your command output, which can make it easier for your users to read and understand.
Otherwise, you can use the `DisableColoring` method to disable colored output.

### Arguments and Flags

You can define arguments within your command using curly braces. For example, if you want to define an argument named "name" that is required, you can use the following code:

```go
command := cli.AddCommand("greet {name}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetArgument("name").String()
    fmt.Println("Hello, ", name)
})
```

You can also define optional arguments by a question mark. For example, if you want to define an argument named "name" that is optional, you can use the following code:

```go
command := cli.AddCommand("greet {name?}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetArgument("name").String()
    fmt.Println("Hello, ", name)
})
```

You can also define multiple arguments by separating them with a space. For example, if you want to define two arguments named "name" and "age", you can use the following code:

```go
command := cli.AddCommand("greet {name} {age}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetArgument("name").String()
    age := c.GetArgument("age").Int()
    fmt.Println("Hello, ", name, "you are", age, "years old")
})
```

You can also define flags for your command. For example, if you want to define a flag named "name" that is required, you can use the following code:

```go
command := cli.AddCommand("greet {--name}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetOption("name").String()
    fmt.Println("Hello, ", name)
})
```

You can also define optional flags by a question mark. For example, if you want to define a flag named "name" that is optional, you can use the following code:

```go
command := cli.AddCommand("greet {--name?}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetOption("name").String()
    fmt.Println("Hello, ", name)
})
```

You can define default values for your arguments and flags by using the `{--name=Name}` syntax. For example, if you want to define a flag named "name" that is optional and has a default value of "Lisa", you can use the following code:

```go
command := cli.AddCommand("greet {--name=Lisa}", "prints a greeting message", func(c *parse.ParsedCommand){
    name := c.GetOption("name").String()
    fmt.Println("Hello, ", name)
})

```