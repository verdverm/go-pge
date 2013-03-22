package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	exec "os/exec"
	"sort"
	"text/template"
)

func post(DS *MainSearch) {

	fmt.Printf("\n\nMain - Post Processing\n======================\n\n")

	DC := DS.cnfg
	// setup log dir and open main log files
	basedir, _ := os.Open(DC.logDir)
	subdirs, _ := basedir.Readdirnames(0)
	sort.Strings(subdirs)

	for _, sdir := range subdirs {
		processProblemDir(sdir, DC.logDir+sdir)
	}

}

func processProblemDir(probstr, dirstr string) {
	outdir := "out/bench/" + probstr

	pdir, _ := os.Open(dirstr)
	rundirs, _ := pdir.Readdirnames(0)
	sort.Strings(rundirs)
	fmt.Printf("ProbDir %s\n", probstr)

	var results []runResults
	for r, rdir := range rundirs {
		ret := processRunDir(r, dirstr+"/"+rdir)
		results = append(results, ret)
	}

	var genErr []aveErr
	for j := 1; j < len(results[0].genErr); j++ {
		ave, best := 0.0, 0.0
		for i := 0; i < len(results); i++ {
			if results[i].genErr[j].gen < 1 {
				continue
			}
			// fmt.Printf("%d %d %d\n", i, j, results[i].genErr[j].gen)
			ave += results[i].genErr[j].err
			best += results[i].genErr[j].best
		}
		gen := j * 10
		ave /= float64(len(results))
		best /= float64(len(results))
		tmp := aveErr{gen, ave, best}
		fmt.Printf("Gen: %v\n", tmp)

		genErr = append(genErr, tmp)
	}

	makeGraph(probstr, outdir, genErr)

	fmt.Println("\n\n\n")
}

type aveErr struct {
	gen  int
	err  float64
	best float64
}

type bestEqn struct {
	pos     int
	eqn_str string
	latex   string
	size    int
	err     float64
}

type runResults struct {
	genErr []aveErr
	eqns   []bestEqn
}

func processRunDir(run int, dirstr string) (ret runResults) {

	fmt.Println("  run ", run, " ", dirstr)

	ret.genErr = make([]aveErr, 101)

	gpsr_eqn_data, _ := ioutil.ReadFile(dirstr + "/gpsr/gpsr:eqns.log")
	lines := bytes.Split(gpsr_eqn_data, []byte("\n"))

	var (
		str, tmp         string
		gen, cnt         int
		pos, size        int
		terr, tsum, tmin float64
		latex, eqn_str   string
	)

	for l := 0; l < len(lines); l++ {
		_, err := fmt.Sscanf(string(lines[l]), "%s %d %d", &str, &gen, &cnt)
		if err != nil {
			// fmt.Println("Err on GEN line: ", err)
			break
		}
		// fmt.Println("    Gen: ", gen)
		tsum = 0.0
		tmin = 1000000.0
		for i := 0; i < cnt; i++ {
			l++
			fmt.Sscanf(string(lines[l]), "%d", &pos)
			latex = string(bytes.TrimSpace(lines[l][2:]))
			l++
			eqn_str = string(lines[l])
			l++
			fmt.Sscanf(string(lines[l]), "%s %d", &str, &size)
			l += 5
			fmt.Sscanf(string(lines[l]), "%s %s %f", &str, &tmp, &terr)
			if terr < tmin {
				tmin = terr
			}
			tsum += terr
			l += 2
			_ = latex
			_ = eqn_str

			if gen == 1000 {
				var eqn bestEqn
				eqn.pos = i
				eqn.eqn_str = eqn_str
				eqn.latex = latex
				eqn.err = terr
				eqn.size = size
				ret.eqns = append(ret.eqns, eqn)
				// if i < 12 {
				// 	fmt.Printf("    %d %s  %d %f\n", i, latex, size, terr)
				// }
			}
		}
		var gerr aveErr
		gerr.gen = gen
		gerr.err = terr
		gerr.best = tmin
		ret.genErr[gen/10] = gerr
		// fmt.Printf("Gen: %d  %f  %f\n", gen, tsum/float64(cnt), tmin)

	}
	return
}

func makeGraph(prob, dir string, genErr []aveErr) {

	dataFN := dir + "/errData.txt"
	gnuFN := dir + "/plotData.gnu"
	plotFN := dir + "/" + prob

	dataFile, _ := os.Create(dataFN)
	out := bufio.NewWriter(dataFile)

	for _, pnt := range genErr {
		fmt.Fprintln(out, pnt.gen, pnt.err, pnt.best)
	}

	out.Flush()
	dataFile.Close()

	pInfo := plotInfo{plotFN, dataFN, prob}

	gnuFile, _ := os.Create(gnuFN)
	gout := bufio.NewWriter(gnuFile)

	tmpl, err := template.New("gnu").Parse(plot_file)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(gout, pInfo)
	if err != nil {
		panic(err)
	}

	gout.Flush()
	gnuFile.Close()

	cmd := exec.Command("gnuplot", gnuFN)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

type plotInfo struct {
	Outfn string
	Infn  string
	Prob  string
}

var plot_file = `
set term postscript enhanced
set output '{{.Outfn}}.png'

set style line 1 lt 1 lc rgb "blue" lw 1
set style line 2 lt 1 lc rgb "black" lw 1
set style line 3 lt 1 lc rgb "red" lw 1
set autoscale
plot '{{.Infn}}' using 1:2 with lines ls 2, '' using 1:3 with lines ls 1

`

var eqn_tex_file = `
\subsection*\{ {{}} \}

`
