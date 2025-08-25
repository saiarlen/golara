package cmd

import (
	"flag"
	"fmt"
	"os"
)

func RunCommands() error {
	if len(os.Args) < 3 {
		ShowHelp()
		return fmt.Errorf("insufficient arguments")
	}

	subCommand := flag.String("subcommand", "", "Provide Subcommand to execute")
	flag.Parse()

	if subCommand == nil {
		return fmt.Errorf("subcommand flag is nil")
	}

	switch *subCommand {
	case "init":
		return InitProject()

	case "help":
		ShowHelp()

	default:
		fmt.Printf("❌ Command '%s' should be run with: go run main.go -subcommand %s\n", *subCommand, *subCommand)
		return fmt.Errorf("use main.go for framework commands")
	}

	return nil
}

func ShowHelp() {
	fmt.Println(`
🔥 Golara Framework CLI

Available commands:
  init                       Initialize a new Golara project
  help                       Show this help message

For other commands, use: go run main.go -subcommand <command>
  migrate                    Run database migrations
  migrate:rollback [steps]   Rollback migrations (default: 1 step)
  migrate:status             Show migration status
  make:controller <name>     Generate a new controller
  make:model <name>          Generate a new model
  make:middleware <name>     Generate a new middleware
  make:job <name>            Generate a new job
  make:view <name>           Generate a new view template
  make:migration <name>      Generate a new migration

Usage:
  ./golara -subcommand init
  go run main.go -subcommand <command> [arguments]

Examples:
  ./golara -subcommand init
  go run main.go -subcommand migrate
  go run main.go -subcommand make:controller User
  go run main.go -subcommand make:model Product`)
}