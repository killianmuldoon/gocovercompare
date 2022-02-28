# gocovercompare

gocovercompare is a tool to directly compare the coverage reports of different versions of a Go library using Go's native code coverage reports. 


## Usage

gocovercompare takes one or more go coverage reports - the file with a `.out` extension produced when running `go test -coverprofile=cover.out ./...` on a project. It produces a CSV report by default containing all packages in all versions of the coverage report. The output is to standard out and can be directed to a file for further usage. By passing the `-output table` argument a table can be produces which shows a summary of coverage across versions.

Using the example files contained in `examples/cluster-api` the following command produces a csv file with coverage for each package in a project across four different coverage profiles:

`./gocovercompare -file examples/cluster-api/0.4.coverage -file examples/cluster-api/1-0.coverage -file examples/cluster-api/1-1.coverage -file examples/cluster-api/main.coverage > cluster-api-coverage.csv`

The coverage percentages per package are organized in the order they are supplied to the tool.