package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/terraform/terraform"
)

func fetchState(dir string) (*terraform.State, error) {
	file := "terraform.tfstate"
	if flag.Arg(0) != "" {
		file = flag.Arg(0)
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open '%s'", file))
	}
	defer f.Close()

	state, err := terraform.ReadState(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read state from '%s'", file))
	}
	return state, nil
}
