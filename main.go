package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/cover"
)

func main() {
	var files fileList
	var output string
	//compare := *flag.Bool("compare", false, "whether to show a comparison table")

	flag.Var(&files, "file", "Coverage file to compare.")
	flag.StringVar(&output, "output", "csv", "kind of output required from tool")

	flag.Parse()
	switch output {
	case "table":
		fmt.Print(coverageTable(files...))
	case "csv":
		fmt.Print(coverageCSV(files...))

	default:
		log.Fatal("output must be set to \"csv\" or \"table\"")
	}

}

func coverageCSV(profileFiles ...string) string {
	var buf strings.Builder

	profiles := []profile{}
	for _, filename := range profileFiles {
		profiles = append(profiles, directoryCoverageFromProfile(filename))
	}
	dirs := directoryList(profiles...)
	for _, name := range dirs {
		var s []string
		for _, prof := range profiles {
			if _, ok := prof[name]; !ok {
				prof[name] = &directory{coveredStatements: -1, totalStatements: 1}
			}
			s = append(s, fmt.Sprintf("%d,%d", prof[name].coveredStatements, prof[name].totalStatements))
		}

		buf.WriteString(fmt.Sprintf("%s,%s\n", name, strings.Join(s, ",")))
	}
	return buf.String()
}

type fileList []string

func (f *fileList) String() string {
	var s string
	for _, item := range *f {
		s = fmt.Sprintf(s, item)
	}
	return s
}

func (f *fileList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func coverageTable(profileFiles ...string) string {

	profiles := []profile{}
	for _, filename := range profileFiles {
		profiles = append(profiles, directoryCoverageFromProfile(filename))
	}
	tableWriter, rowFormat := generateHeader(profileFiles...)
	dirs := directoryList(profiles...)
	for _, name := range dirs {
		var rowValues []interface{}
		rowValues = append(rowValues, name)

		for _, prof := range profiles {

			if _, ok := prof[name]; !ok {
				prof[name] = &directory{coveredStatements: -1, totalStatements: 1}
			}
			rowValues = append(rowValues, prof[name].coverageString())
		}
		tableWriter.WriteString(fmt.Sprintf(rowFormat, rowValues...))
	}
	return tableWriter.String()
}

func generateHeader(profileFiles ...string) (*strings.Builder, string) {
	var buf strings.Builder
	rowFormat := "%-80s"
	line1 := []interface{}{"package"}
	line2 := []interface{}{"--------"}
	for i, _ := range profileFiles {
		rowFormat = rowFormat + " %8s"
		line1 = append(line1, strconv.Itoa(i+1))
		line2 = append(line2, "------")
	}
	rowFormat = rowFormat + "\n"
	buf.WriteString(fmt.Sprintf(rowFormat, line1...))
	buf.WriteString(fmt.Sprintf(rowFormat, line2...))
	return &buf, rowFormat
}

func description(first, second float64) string {
	if first < 0 {
		return "new"
	}
	if second < 0 {
		return "removed"
	}

	return fmt.Sprintf("%+6.2f%%", second-first)
}

func directoryList(profiles ...profile) []string {
	set := map[string]struct{}{}
	for _, p := range profiles {
		for name, _ := range p {
			set[name] = struct{}{}
		}
	}
	var res []string
	for name := range set {
		res = append(res, name)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}

func directoryCoverageFromProfile(filename string) profile {
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: coverage file %v not valid", err, filename))
	}
	dirs := map[string]*directory{}
	for _, p := range profiles {
		pkgName := filepath.Dir(p.FileName)
		if _, ok := dirs[pkgName]; !ok {
			dirs[pkgName] = &directory{0, 0}
		}
		for _, block := range p.Blocks {
			if block.Count > 0 {
				dirs[pkgName].coveredStatements += block.NumStmt
			}
			dirs[pkgName].totalStatements += block.NumStmt
		}
	}
	return dirs
}

type profile map[string]*directory

func (p profile) coverageString() string {
	return fmt.Sprintf("%4.2f%%", p.coverage())
}

func (p profile) coverage() float64 {
	var coveredStatements, totalSatements int
	for _, d := range p {
		coveredStatements += d.coveredStatements
		totalSatements += d.totalStatements
	}
	return 100 * float64(coveredStatements) / float64(totalSatements)
}

type directory struct {
	coveredStatements int
	totalStatements   int
}

func (d directory) coverageString() string {
	if d.coverage() < 0 {
		return "-"
	}
	return fmt.Sprintf("%4.2f%%", d.coverage())
}

func (d directory) coverage() float64 {
	coverage := 100 * float64(d.coveredStatements) / float64(d.totalStatements)
	return coverage
}
