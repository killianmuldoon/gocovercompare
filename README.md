# gocovercompare

gocovercompare is a tool to directly compare the coverage reports of different versions of a Go library using Go's native code coverage reports. 


## Usage

gocovercompare takes one or more go coverage reports - produced when running `go test -coverprofile=cover.out ./...` on a project. It produces a CSV report by default containing a test coverage comparison for all packages in all versions of the coverage report. The coverage percentages per package are organized in the order they are supplied to the tool. 

The output is to standard out and can be directed to a file for further usage. By passing the `-output table` argument a table can be produced which shows a summary of coverage across versions.

Using the example files contained in `examples/cluster-api` the following command produces a csv file with coverage for each package in a project across four different coverage profiles:

```bash
./gocovercompare -file examples/cluster-api/0.4.coverage -file examples/cluster-api/1-0.coverage -file examples/cluster-api/1-1.coverage -file examples/cluster-api/main.coverage > cluster-api-coverage.csv
```

Using the `-output=table argument`, as below, produces a formatted table to standard out.
```bash
./gocovercompare -file examples/cluster-api/0.4.coverage -file examples/cluster-api/1-0.coverage -file examples/cluster-api/1-1.coverage -file examples/cluster-api/main.coverage -output table`

package                                                                                 1        2        3        4
--------                                                                           ------   ------   ------   ------
sigs.k8s.io/cluster-api/api/v1alpha3                                               36.07%   35.97%   35.97%   36.12%
sigs.k8s.io/cluster-api/api/v1alpha4                                               48.88%   32.26%   37.22%   37.24%
sigs.k8s.io/cluster-api/api/v1alpha4/index                                         55.32%        -        -        -
sigs.k8s.io/cluster-api/api/v1beta1                                                     -   48.97%   40.60%   40.65%
sigs.k8s.io/cluster-api/api/v1beta1/index                                               -   55.32%   55.32%   55.32%

```

