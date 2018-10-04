// Copyright © 2018 Jeremy Canady

package cmd

import (
	"fmt"

	"github.com/jrmycanady/logstash-bench/bench"
	"github.com/spf13/cobra"
)

// genconfigCmd represents the genconfig command
var devGenconfigCmd = &cobra.Command{
	Use:   "genconfig",
	Short: "Generates the configuration file.",

	Run: func(cmd *cobra.Command, args []string) {
		s, err := bench.BuildConfig(inputCodec, sourceFilePath, fmt.Sprintf("%s/output.log", tempDirPath), filterFilePath, fmt.Sprintf("%s/sincedb.log", tempDirPath))
		if err != nil {
			fmt.Printf("failed to build config: %s\n", err)
			return
		}
		fmt.Println(s)
	},
}

func init() {
	devCmd.AddCommand(devGenconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
