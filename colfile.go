package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s inputfile\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

// Columnfile is for an ImageD11 columnfile
// Dictionary of parameters + double precision table
type Columnfile struct {
	parameters map[string]string
	titles     []string
	nrows      int
	ncols      int
	data       [][]float64
}

func readColumnfile(filename string) (cf Columnfile) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error! : ", err)
		return
	}
	defer file.Close()
	ierr := 0
	scanner := bufio.NewScanner(file)

	cf.parameters = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 { // skip empty lines
			continue
		}
		if line[0] == "#"[0] {
			if strings.Contains(line, "=") {
				fs := strings.Fields(line)
				cf.parameters[fs[1]] = fs[3]
			} else {
				cf.titles = strings.Fields(line)[1:]
				cf.ncols = len(cf.titles)
			}
		} else {
			ierr = 0
			s := make([]float64, cf.ncols)
			for i, v := range strings.Fields(line) {
				if i >= cf.ncols {
					ierr++
					fmt.Println("Error on line", cf.nrows)
					fmt.Println("line:", line)
					break // line too long
				}
				s[i], err = strconv.ParseFloat(v, 64)
				if err != nil {
					ierr++
				}
			}
			if ierr == 0 {
				cf.nrows++
				cf.data = append(cf.data, s)
			}
		}
	}
	return cf
}

func stats(data [][]float64, c int) (mini, maxi, mean float64) {
	mini = data[0][c]
	maxi = data[0][c]
	mean = 0.0
	for _, row := range data {
		if row[c] < mini {
			mini = row[c]
		}
		if row[c] > maxi {
			maxi = row[c]
		}
		mean += row[c]
	}
	mean = mean / float64(len(data))
	return mini, maxi, mean
}

func printColumnfile(cf Columnfile) {
	fmt.Printf("Columns: %d Rows %d\n", cf.ncols, cf.nrows)
	for k, v := range cf.parameters {
		fmt.Printf("Parameter %s = %s\n", k, v)
	}
	for i, t := range cf.titles {
		mini, maxi, mean := stats(cf.data, i)
		fmt.Printf("Column named %s : from %f to %f mean %f\n", t, mini, maxi, mean)
	}

}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
	cf := readColumnfile(args[0])
	printColumnfile(cf)
}
