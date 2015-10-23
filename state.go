package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/terraform"
)

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

func openState(dir string) (io.ReadCloser, error) {
	file := filepath.Join(dir, "terraform.tfstate")
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open '%s'", file))
	}
	return f, nil
}
