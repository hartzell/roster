package main

import (
	"flag"

	"github.com/mitchellh/cli"
)

type DefaultCommand struct {
	Dir string
	Ui  cli.Ui
	FS  *flag.FlagSet
}

func (dc *DefaultCommand) Help() string {
	// stub this out.  It never seems to get called, so that fact that
	// it's shared amongst all of the Commands isn't a problem.  I'm not
	// sure *why* I'm getting lucky and need to walk through it, but for
	// now, just take it.
	return ""
}

// InitFlagSet intializes the DefaultCommand's FS element and adds
// default flags.
func (dc *DefaultCommand) InitFlagSet() {
	dc.FS = flag.NewFlagSet("inventory", flag.ContinueOnError)
	dc.FS.StringVar(&dc.Dir, "dir", ".", "The path to the terraform directory")
}
