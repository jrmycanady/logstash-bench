#!/bin/bash

env GOOS=linux go build -o logstash-bench-linux-x64
env GOOS=darwin go build -o logstash-bin-darwin-x64

scp ./logstash-bench-linux-x64 root@jcdesktop.inf.labs/logstashtest