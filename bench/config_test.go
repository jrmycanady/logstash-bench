package bench

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestFilterFromFile(t *testing.T) {
	s, err := FilterFromFile("./sample_data/filter.conf")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(s)

	es, err := ioutil.ReadFile("./sample_data/expected_output.conf")
	if err != nil {
		t.Fatal(err)
	}

	if string(es) != s {
		t.Fatalf("strings do not match:\n%s\n----\n%s", s, es)
	}
}

func TestBuildInput(t *testing.T) {
	s := BuildInput("json", "/file/path/logs.log", "/file/db/path/bench.db")

	if s != `input {file {codec => "json" mode => "tail" path => "/file/path/logs.log" sincedb_path => "/file/db/path/bench.db" start_position => "beginning"} }` {
		t.Fatal(s)
	}

}

func TestBuildOutput(t *testing.T) {
	s := BuildOutput("/file/path/output.log")
	fmt.Println(s)
	if s != `output {file {path => "/file/path/output.log"} }` {
		t.Fatal()
	}
}

func TestBuildConfig(t *testing.T) {

	s, err := BuildConfig("json", "/file/path/logs.log", "/file/path/output.log", "./sample_data/expected_output.conf", "/file/db/path/bench.db")
	if err != nil {
		t.Fatal()
	}

	es, err := ioutil.ReadFile("./sample_data/expected_config.conf")
	if err != nil {
		t.Fatal("failed to read expected config")
	}

	if s != string(es) {
		t.Fatal(s)
	}
}
