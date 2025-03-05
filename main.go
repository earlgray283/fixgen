package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/earlgray283/fixgen/internal/cmd"
	"github.com/earlgray283/fixgen/internal/cmdutil"
)

func main() {
	rootCmd := cmd.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		var exitCodeErr *cmdutil.ExitCodeError
		if errors.As(err, &exitCodeErr) {
			os.Exit(exitCodeErr.ExitCode)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
