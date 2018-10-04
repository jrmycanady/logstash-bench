// Copyright © 2018 Jeremy Canady

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jrmycanady/logstash-bench/bench"
	"github.com/spf13/cobra"
)

var (
	inputCodec     string
	sourceFilePath string
	tempDirPath    string
	numWorkers     int64
	numIterations  int64
	noFileOutput   bool
	logstashPath   string
	filterFilePath string
	details        bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logstash-bench",
	Short: "Performs logstash filter performance benchmarking.",
	Run:   run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputCodec, "input-codec", "c", "json", "The codec logstash should use on input.")
	rootCmd.Flags().StringVarP(&sourceFilePath, "source-file-path", "s", "./input.log", "The path to the input source file.")
	rootCmd.Flags().StringVarP(&tempDirPath, "temp-dir-path", "t", "/dev/shm/", "The path to a readable and writable directory for temporary file storage.")
	rootCmd.Flags().Int64VarP(&numWorkers, "number-of-workers", "w", 1, "The number of workers to start logstash with.")
	rootCmd.Flags().Int64VarP(&numIterations, "number-of-iterations", "i", 1, "The number of time the test should be ran.")
	rootCmd.Flags().StringVarP(&logstashPath, "logstash-executable-path", "l", "/usr/share/logstash/bin/logstash", "The path to the logstash executable.")
	rootCmd.Flags().StringVarP(&filterFilePath, "filter-file-path", "f", "./filter.conf", "The path to the filter to test.")
	rootCmd.Flags().BoolVarP(&details, "details", "d", false, "Shows details of the process on stdout.")
}

func run(cmd *cobra.Command, args []string) {

	cfg := bench.RunCfg{
		InputCodec:     inputCodec,
		SourceFilePath: sourceFilePath,
		TempDirPath:    tempDirPath,
		NumWorkers:     numWorkers,
		NumIterations:  numIterations,
		LogstashPath:   logstashPath,
		FilterFilePath: filterFilePath,
	}

	// Disable detail output if needed.
	if !details {
		log.SetOutput(ioutil.Discard)
	}

	log.Print("starting run")
	r, err := bench.Exec(cfg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r.Screen())

}
