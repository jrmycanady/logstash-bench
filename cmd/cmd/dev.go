// Copyright © 2018 Jeremy Canady

package cmd

import (
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development and debugging related commands.",
}

func init() {
	rootCmd.AddCommand(devCmd)

}
