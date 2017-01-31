package main

import (
	"fmt"
	"os"

	"github.com/bfirsh/whalebrew/cmd"
)

func main() {

	if len(os.Args) > 1 {
		// Check if not command exists
		if _, _, err := cmd.RootCmd.Find(os.Args); err != nil {
			// Check if file exists
			if _, err := os.Stat(os.Args[1]); err == nil {
				cmd.RootCmd.SetArgs(append([]string{"run"}, os.Args[1:]...))
			}
		}
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
