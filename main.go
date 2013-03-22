package main

import (
	"bufio"
	"flag"
	"fmt"
	log "log"
	rand "math/rand"
	os "os"
	rt "runtime"
	pprof "runtime/pprof"
	"strings"

	probs "github.com/verdverm/go-pge/problems"
	expr "github.com/verdverm/go-symexpr"
	tmplt "text/template"
)

var debug = false
var numProcs = 12

var cpuprofile = flag.String("prof", "", "write cpu profile to file")

var arg_cfg = flag.String("cfg", "config/main/main_default.cfg", cfg_help_str)
var arg_pcfg = flag.String("pcfg", "", pcfg_help_str)
var arg_scfg = flag.String("scfg", "", scfg_help_str)
var arg_gen = flag.String("gen", "", gen_help_str)
var arg_tmp = flag.Bool("tmp", false, "run tmp code and exit")
var arg_post = flag.Bool("post", false, "run output processing code and exit")

var arg_pge_iter = flag.Int("iter", -1, "iterations for PGE")
var arg_pge_peel = flag.Int("peel", -1, "peel count for PGE")
var arg_pge_init = flag.String("init", "", "initialize function for PGE")
var arg_pge_grow = flag.String("grow", "", "expansion function for PGE")

var cfg_help_str = "A main config file"
var pcfg_help_str = "A Problem config file"
var scfg_help_str = "A Search config file"
var gen_help_str = "Generate Data [bench,diffeq]:[list,all,probname]"

func main() {

	flag.Parse()

	initGo()

	expr.DumpExprTypes()

	if *arg_tmp {
		printBenchLatex()
		// tmp()
		return
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// if arg_gen = something, then generate data and exit
	if *arg_gen != "" {
		if strings.HasPrefix(strings.ToLower(*arg_gen), "bench") {
			genBenchmark(*arg_gen)
		} else {
			fmt.Printf("NOT generating %s data, mwahahaha\n", *arg_gen)
		}
		return
	}

	/*******************************
				Main Code			
	 *******************************/

	var DS MainSearch

	DS.ParseConfig(*arg_cfg)

	if *arg_pcfg != "" {
		DS.cnfg.probCfg = *arg_pcfg
	}
	if *arg_scfg != "" {
		DS.cnfg.srchCfg = []string{*arg_scfg}
	}

	if *arg_post {
		post(&DS)
		return
	}

	initDone := make(chan int)

	DS.Init(initDone, nil)
	// fmt.Printf("initd: %v\n", DS)

	DS.Run()

}

func initGo() {
	fmt.Printf("Initializing Go System (%d threads)\n", numProcs)
	if debug {
		rt.GOMAXPROCS(2)
		rand.Seed(0)
	} else {
		rt.GOMAXPROCS(numProcs)
		rand.Seed(rand.Int63())
	}
}

func tmp() {

	T, err := tmplt.ParseFiles("config/prob/prob_template.txt")
	if err != nil {
		log.Fatal("Error reading template file: ", err)
	}

	fmt.Println("T.name: ", T.Name())

	for _, B := range probs.BenchmarkList {

		ftotal, err2 := os.Create("config/prob/bench/" + B.Name + ".cfg")
		if err2 != nil {
			fmt.Printf("Error creating config file: %s  %v\n", B.Name, err)
			return
		}
		file := bufio.NewWriter(ftotal)

		err = T.Execute(os.Stdout, B)
		err = T.Execute(file, B)

		file.Flush()
		ftotal.Close()
	}
}
