// Package bench provides benchmarking tools for logstash.
package bench

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RunCfg struct {
	InputCodec     string // The input codes to use.
	SourceFilePath string // The path to the source file.
	TempDirPath    string // The path to the temporary directory.
	NumWorkers     int64  // The number of logstash workers to use.
	NumIterations  int64  // The number of times the test should be ran.
	LogstashPath   string // The past to the logstash binary.
	FilterFilePath string // The path to the filter file to test.
}

// Exec executes the run in the configuration.
func Exec(cfg RunCfg) (Result, error) {
	// Generating unique ID for this run.
	id, err := uuid.NewRandom()
	if err != nil {
		return Result{}, fmt.Errorf("generating run unqiue id: %s", err)
	}
	log.Printf("run id is %s\n", id.String())

	// Building temporary file paths and cleanup from previous runs.
	runTempDirPath := fmt.Sprintf("%s/%s", cfg.TempDirPath, id.String())
	if err := os.MkdirAll(runTempDirPath, 0600); err != nil {
		return Result{}, fmt.Errorf("creating temporary directory: %s", err)
	}
	defer cleanupTempDir(runTempDirPath)
	log.Printf("temp directory is %s\n", runTempDirPath)

	outputFilePath := fmt.Sprintf("%s/output.json", runTempDirPath)
	sinceFileDBPath := fmt.Sprintf("%s/sincedb.db", runTempDirPath)
	configFilePath := fmt.Sprintf("%s/filter.conf", runTempDirPath)

	// Finding absolute paths for all files. Logstash will not work with relative paths.
	sourceFilePathAbs, err := filepath.Abs(cfg.SourceFilePath)
	if err != nil {
		return Result{}, fmt.Errorf("finding absolute path for sourceFilePath: %s", err)
	}
	log.Printf("source file path is %s\n", sourceFilePathAbs)

	outputFilePathAbs, err := filepath.Abs(outputFilePath)
	if err != nil {
		return Result{}, fmt.Errorf("finding absolute file path for outputFilePath: %s", err)
	}
	log.Printf("output file path is %s\n", outputFilePathAbs)

	sinceFileDBPathAbs, err := filepath.Abs(sinceFileDBPath)
	if err != nil {
		return Result{}, fmt.Errorf("finding absolute path for sinceFileDBPath: %s", err)
	}
	log.Printf("sincedb file path is %s\n", sinceFileDBPathAbs)

	// Generate configuration file data.
	cfgStr, err := BuildConfig(cfg.InputCodec, sourceFilePathAbs, outputFilePathAbs, cfg.FilterFilePath, sinceFileDBPathAbs)
	if err != nil {
		return Result{}, fmt.Errorf("building configuration: %s", err)
	}
	if err := ioutil.WriteFile(configFilePath, []byte(cfgStr), 0600); err != nil {
		return Result{}, fmt.Errorf("writing config file to disk: %s", err)
	}
	log.Printf("config file generated at %s", configFilePath)

	// Start the command.
	log.Printf("starting logstash at %s", cfg.LogstashPath)
	cmd := exec.Command(cfg.LogstashPath, "-f", configFilePath, "-w", strconv.FormatInt(cfg.NumWorkers, 10))
	if err := cmd.Start(); err != nil {
		return Result{}, fmt.Errorf("starting command call: %s", err)
	}

	// Monitor for completion
	var completeTime time.Time
	go func() {
		log.Println("monitoring for sincedb file creation")

		// Monitor sync file
		for {
			fi, err := os.Stat(sinceFileDBPathAbs)
			// if err != nil {
			// 	log.Println("waiting for sync db file to be created")
			// }
			if err == nil {

				if fi.Size() > 0 {
					log.Printf("sincedb file size is %d, stopping logstash", fi.Size())
					cmd.Process.Kill()
				}
			}

			time.Sleep(time.Duration(1) * time.Second)

		}

	}()

	// Waiting for logstash to finish.
	log.Printf("waiting for logstash to stop")
	cmd.Wait()
	completeTime = time.Now()
	log.Printf("logstash has stopped processing")

	// Process results.
	log.Printf("beginning result processing")
	fi, err := os.Stat(outputFilePathAbs)
	if err != nil {
		return Result{}, fmt.Errorf("failed to stat output file: %s", err)
	}

	f, err := os.Open(outputFilePathAbs)
	if err != nil {
		return Result{}, fmt.Errorf("opening output file for processing: %s", err)
	}
	defer f.Close()

	// Read the first line
	reader := bufio.NewReader(f)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		return Result{}, fmt.Errorf("reading output file: %s", err)
	}

	if line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}

	l := LogLine{}
	err = json.Unmarshal(line, &l)
	if err != nil {
		return Result{}, fmt.Errorf("decoding first line: %s", err)
	}

	sfi, err := os.Stat(sourceFilePathAbs)
	if err != nil {
		return Result{}, fmt.Errorf("stating source file for size: %s", err)
	}

	return Result{
		FilterFile:            cfg.FilterFilePath,
		FirstProcessedAt:      l.ProcessedAt,
		LastProcessedAt:       completeTime,
		Duration:              completeTime.Sub(l.ProcessedAt),
		InputFileSize:         sfi.Size(),
		OutputFileSize:        fi.Size(),
		FileSizePercentChange: float64(fi.Size()-sfi.Size()) / float64(sfi.Size()) * 100.0,
	}, nil
}

type LogLine struct {
	ProcessedAt time.Time `json:"processed_at"`
}

type Result struct {
	FilterFile            string
	FirstProcessedAt      time.Time
	LastProcessedAt       time.Time
	Duration              time.Duration
	InputFileSize         int64
	OutputFileSize        int64
	FileSizePercentChange float64
}

func (r Result) Screen() string {
	return fmt.Sprintf("Filter File: %s\nProcesssing Started: %s\nProcessing Ended: %s\nDuration: %0.4f\nInput File Size: %d\nOutput File Size: %d\nFile Size Change Percentage: %0.2f", r.FilterFile, r.FirstProcessedAt.Format(time.RFC3339), r.LastProcessedAt.Format(time.RFC3339), r.Duration.Seconds(), r.InputFileSize, r.OutputFileSize, r.FileSizePercentChange)
}

// cleanupTempDir cleans up the directory specified by path.
func cleanupTempDir(path string) {
	os.RemoveAll(path)
}
