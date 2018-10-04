# logstash-bench
lgostash-bench is a utility that performance tests logstash filters. This is primarily used to aid testing filters before they are pushed to production.

## Usage Basics

1. Obtain a sample of the ingest data. 
    * There are many ways to obtain that data but the simplest is to remove any filters and set the output of logstash to a file and let it run. 
    * The larger the sample file the better the benchmark data will be.
    * If possible locate the file on a memory backed location such as /dev/shm
2. Create a temporary directory to use.
    * The temporary directory will store the output data as well as the sincedb file.
    * Recommend placing it in a location like /dev/shm or another memory backed disk.
3. Select your filter file
    * The filter file should be only the filter {} section.
4. Run! ```logstash-bench -f file.conf -s source.logs```

## Parameters

|long name|short name|description|default|
|---------|----------|-----------|-------|
|input-codec|c|The codec that logstash should use for input.|json|
|source-file-path|s|The path to the source log file.|./input.log|
|temp-dir-path|t|The path to the temporary directory.|/dev/shm/|
|number-of-iterations|i|The number of times the test will run.|1|
|logstash-executable-path|l|The path to the logstash executable.|/usr/share/logstash/bin/logstash|
|filter-file-path|f|The path to the filter file to test.|./filter.conf|
|details|d|Show the details of the process via stdout.|false|

