package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/terraform"
)

func main() {

	exitStatus, err := doIt()
	if exitStatus != 0 {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func doIt() (exitStatus int, errorMessage string) {
	var list bool
	var host string

	flag.BoolVar(&list, "list", false, "Run the list command")
	flag.StringVar(&host, "host", "", "specify a host")
	flag.Parse()

	if list && host != "" {
		return 1,
			"Provide either \"--list\" or \"--host\", not both (\"--help\" for help)."
	}

	if host != "" {
		return doHost()
	}

	return doList()
}

func doHost() (exitStatus int, errorMessage string) {
	fmt.Println("{}")
	return 0, ""
}

func doList() (exitStatus int, errorMessage string) {
	f, err := os.Open("terraform.tfstate")
	if err != nil {
		return 1, "Unable to open state file "
	}
	state, err := terraform.ReadState(f)
	if err != nil {
		return 1, "Unable to read state file"
	}
	spew.Dump(state)

	return 0, ""
}
