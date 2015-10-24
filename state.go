package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/terraform"
)

// Function fetchState takes a directory name (a string) and returns
// an instance of *terraform.State based on what it finds in that
// directory and an error.
func fetchState(dir string) (*terraform.State, error) {
	src, err := openState(dir)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	state, err := terraform.ReadState(src)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read state from src"))
	}
	return state, nil
}

// Function openState takes a directory string and returns an open
// file for the terraform.tfstate file in that directory and an error.
func openState(dir string) (io.ReadCloser, error) {
	file := filepath.Join(dir, "terraform.tfstate")
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open '%s'", file))
	}
	return f, nil
}
