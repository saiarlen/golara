package main

import (
	"ekycapp/database/migrations"
	"flag"
	"fmt"
)

func Commands() error {
	subCommand := flag.String("subcommand", "", "Provide Subcommand to execute")
	flag.Parse()

	switch *subCommand {
	case "migrate":
		migrations.Migrations()
		fmt.Println("Migrations are done!")

	case "b":
		fmt.Println("B Excuted")

	default:
		fmt.Println("Invalid subcommand")
	}

	return nil

}
