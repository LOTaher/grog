package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
    Use: "grog",
    Short: "grog is a lightweight node package manager written in go.",
}

func Execute() {
    err := root.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    // root.AddCommand(Command Name)
	root.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
