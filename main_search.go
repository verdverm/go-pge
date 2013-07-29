package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	config "github.com/verdverm/go-pge/config"
	pge "github.com/verdverm/go-pge/pge"
	probs "github.com/verdverm/go-pge/problems"
)

// defines the interface to a search type [GP,PE]
type Search interface {

	// parse a params file
	ParseConfig(filename string)

	// initialize the search, sending signal on chan when done
	// the input will be something for the search to connect to
	// in order to provide updates, be monitored, and receive control signals
	Init(done chan int, prob *probs.ExprProblem, logdir string, input interface{})

	// start the actual search procedure (Init is required before a new call to Run)
	Run()

	// clean up internal structures (Init is required before a new call to Run)
	Clean()
}

// parameters to a main search, which sets up the global system
// and instructs in where to find the sub-searches
type mainConfig struct {
	dataDir string
	cfgDir  string
	logDir  string

	probCfg string
	srchCfg []string
}

func mainConfigParser(field, value string, config interface{}) (err error) {

	DC := config.(*mainConfig)

	switch strings.ToUpper(field) {
	case "CONFIGDIR":
		DC.cfgDir = value
	case "DATADIR":
		DC.dataDir = value
	case "LOGDIR":
		DC.logDir = value

	case "PROBLEMCFG":
		DC.probCfg = value
	case "SEARCHCFG":
		DC.srchCfg = strings.Fields(value)
	default:
		log.Printf("Main Not Implemented  %s, %s\n\n", field, value)

	}
	return err
}

// defines the global level search, which may use different sub-searches one or more times
type MainSearch struct {
	cnfg mainConfig

	// problem and best results
	prob     *probs.ExprProblem
	eqns     probs.ExprReportArray
	per_eqns []*probs.ExprReportArray

	// sub-searches and comm
	srch []Search
	comm []*probs.ExprProblemComm
	iter []int

	// logs
	logDir     string
	mainLog    *log.Logger
	mainLogBuf *bufio.Writer
	eqnsLog    *log.Logger
	eqnsLogBuf *bufio.Writer
	errLog     *log.Logger
	errLogBuf  *bufio.Writer
}

func (DS *MainSearch) ParseConfig(filename string) {
	fmt.Printf("Parsing Main Config: %s\n", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = config.ParseConfig(data, mainConfigParser, &DS.cnfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", DS.cnfg)
}

func (DS *MainSearch) Init(done chan int, input interface{}) {
	fmt.Printf("Init'n PGE1\n----------\n")

	DC := DS.cnfg

	// read and setup problem
	eprob := new(probs.ExprProblem)

	fmt.Printf("Parsing Problem Config: %s\n", DC.probCfg)
	data, err := ioutil.ReadFile(DC.cfgDir + DC.probCfg)
	if err != nil {
		log.Fatal(err)
	}
	err = config.ParseConfig(data, probs.ProbConfigParser, eprob)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Prob: %v\n", eprob)
	fmt.Printf("TCfg: %v\n\n", eprob.TreeCfg)

	// setup log dir and open main log files
	DC.logDir += eprob.Name + "/"
	if DC.srchCfg[0][:4] == "pge1" {
		DC.logDir += "pge1/"
	} else if DC.srchCfg[0][:3] == "pge" {
		DC.logDir += "pge/"
	}

	os.MkdirAll(DC.logDir, os.ModePerm)

	now := time.Now()
	fmt.Println("LogDir: ", DC.logDir)
	os.MkdirAll(DC.logDir, os.ModePerm)
	DS.initLogs(DC.logDir)

	DS.mainLog.Println(DC.logDir, now)

	// // setup data
	fmt.Printf("Setting up problem: %s\n", eprob.Name)

	eprob.Train = make([]*probs.PointSet, len(eprob.TrainFns))
	for i, fn := range eprob.TrainFns {
		fmt.Printf("Reading Trainging File: %s\n", fn)
		eprob.Train[i] = new(probs.PointSet)
		if strings.HasSuffix(fn, ".dataF2") || strings.HasSuffix(fn, ".mat") {
			eprob.Train[i].ReadLakeFile(DC.dataDir + fn)
		} else {
			eprob.Train[i].ReadPointSet(DC.dataDir + fn)
		}
	}
	eprob.Test = make([]*probs.PointSet, len(eprob.TestFns))
	for i, fn := range eprob.TestFns {
		fmt.Printf("Reading Testing File: %s\n", fn)
		eprob.Test[i] = new(probs.PointSet)
		if strings.HasSuffix(fn, ".dataF2") || strings.HasSuffix(fn, ".mat") {
			eprob.Test[i].ReadLakeFile(DC.dataDir + fn)
		} else {
			eprob.Test[i].ReadPointSet(DC.dataDir + fn)
		}
	}

	DS.prob = eprob
	fmt.Println()

	// read search configs
	for _, cfg := range DC.srchCfg {
		if cfg[:4] == "pge1" {
			GS := new(pge.PgeSearch)
			GS.ParseConfig(DC.cfgDir + cfg)
			DS.srch = append(DS.srch, GS)
		} else if cfg[:3] == "pge" {
			PS := new(pge.PgeSearch)
			PS.ParseConfig(DC.cfgDir + cfg)
			DS.srch = append(DS.srch, PS)

			/************/
			// temporary hack
			DS.prob.MaxIter = PS.GetMaxIter()
			/************/
			if *arg_pge_iter >= 0 {
				DS.prob.MaxIter = *arg_pge_iter
				PS.SetMaxIter(*arg_pge_iter)
			}
			if *arg_pge_peel >= 0 {
				PS.SetPeelCount(*arg_pge_peel)
			}
			if *arg_pge_init != "" {
				PS.SetInitMethod(*arg_pge_init)
			}
			if *arg_pge_grow != "" {
				PS.SetGrowMethod(*arg_pge_grow)
			}
			PS.SetEvalrCount(*arg_pge_evals)

		} else {
			log.Fatalf("unknown config type: %v  from  %v\n", cfg[:4], cfg)
		}
	}

	// setup best results
	DS.eqns = make(probs.ExprReportArray, 32)
	DS.per_eqns = make([]*probs.ExprReportArray, len(DS.srch))

	// setup communication struct
	DS.comm = make([]*probs.ExprProblemComm, len(DS.srch))
	for i, _ := range DS.comm {
		DS.comm[i] = new(probs.ExprProblemComm)
		DS.comm[i].Cmds = make(chan int)
		DS.comm[i].Rpts = make(chan *probs.ExprReportArray, 64)
		DS.comm[i].Gen = make(chan [2]int, 64)
	}

	DS.iter = make([]int, len(DS.srch))

	fmt.Println("\n******************************************************\n")

	// initialize searches
	sdone := make(chan int)
	for i, _ := range DS.srch {
		DS.srch[i].Init(sdone, eprob, DC.logDir, DS.comm[i])
	}
	fmt.Println("\n******************************************************\n")

}

func (DS *MainSearch) Run() {
	fmt.Printf("Running Main\n")
	fmt.Println("numSrch = ", len(DS.srch))
	for i := 0; i < len(DS.srch); i++ {
		go DS.srch[i].Run()
	}
	counter := 0
	for {
		// fmt.Println("DS: ", counter)
		DS.checkMessages()

		// time.Sleep(time.Second / 20)
		counter++

		if DS.checkStop() {
			DS.doStop()
			break
		}
	}

	for i, R := range DS.eqns {
		if R == nil || R.Expr() == nil {
			continue
		}
		trn := DS.prob.Train[0]
		f_x := "df(" + trn.GetIndepNames()[DS.prob.SearchVar] + ")"
		str := R.Expr().PrettyPrint(trn.GetIndepNames(), trn.GetSysNames(), R.Coeff())
		fmt.Printf("%d: %s = %s\n%v\n\n", i, f_x, str, R)
	}

	DS.Clean()

	fmt.Println("DS leaving Run()")
}

func (DS *MainSearch) Clean() {
	fmt.Printf("Cleaning Main\n")

	DS.errLogBuf.Flush()
	DS.mainLogBuf.Flush()
	DS.eqnsLogBuf.Flush()

}

func (DS *MainSearch) checkStop() bool {
	if DS.iter[0] > DS.prob.MaxIter {
		return true
	}
	return false
}

func (DS *MainSearch) doStop() {
	done := make(chan int)

	for i, _ := range DS.comm {
		func() {
			c := i
			C := DS.comm[i]
			go func() {
				C.Cmds <- -1
				fmt.Printf("DS sent -1 to Srch %d\n", c)
				<-C.Cmds
				done <- 1
			}()
		}()
	}

	cnt := 0
	for cnt < len(DS.comm) {
		DS.checkMessages()
		_, ok := <-done
		if ok {
			cnt++
			fmt.Println("DS done = ", cnt, len(DS.comm))
		}

	}

	fmt.Println("DAMD checking last messages")
	DS.checkMessages()

	fmt.Println("DS done stopping")
}

func (DS *MainSearch) checkMessages() {
	msg := false
	for i := 0; i < len(DS.comm); i++ {
		select {
		case gen, ok := <-DS.comm[i].Gen:
			if ok {
				DS.iter[i] = gen[1]
				if gen[0] == 0 {
					fmt.Println("Gen: ", gen[1])
				}
				i--
				msg = true
			}
		case rpt, ok := <-DS.comm[i].Rpts:
			if ok {
				msg = true
				DS.per_eqns[i] = rpt
				i--
			}
		default:
			continue
		}
	}
	if !msg {
		time.Sleep(time.Millisecond)
	}
	DS.accumExprs()
}

func (DS *MainSearch) accumExprs() {
	union := make(probs.ExprReportArray, 0)
	for i := 0; i < len(DS.per_eqns); i++ {
		if DS.per_eqns[i] != nil {
			union = append(union, (*DS.per_eqns[i])[:]...)
		}
	}
	union = append(union, DS.eqns[:]...)

	// remove duplicates
	sort.Sort(union)
	last := 0
	for last < len(union) && union[last] == nil {
		last++
	}
	for i := last + 1; i < len(union); i++ {
		if union[i] == nil {
			continue
		}
		if union[i].Expr().AmIAlmostSame(union[last].Expr()) {
			union[i] = nil
		} else {
			last = i
		}
	}

	queue := probs.NewQueueFromArray(union)
	queue.SetSort(probs.GPSORT_PARETO_TST_ERR)
	queue.Sort()

	copy(DS.eqns, union[:len(DS.eqns)])

	// DS.eqnsLog.Printf("\n\n\nLatest Eqns:\n")
	// DS.eqnsLog.Println(DS.eqns)
	// DS.Clean()

}

func (DS *MainSearch) initLogs(logdir string) {

	// open logs
	DS.logDir = logdir
	os.Mkdir(DS.logDir, os.ModePerm)
	tmpF0, err5 := os.Create(DS.logDir + "main:err.log")
	if err5 != nil {
		log.Fatal("couldn't create errs log", err5)
	}
	DS.errLogBuf = bufio.NewWriter(tmpF0)
	DS.errLogBuf.Flush()
	DS.errLog = log.New(DS.errLogBuf, "", log.LstdFlags)

	tmpF1, err1 := os.Create(DS.logDir + "main:main.log")
	if err1 != nil {
		log.Fatal("couldn't create main log", err1)
	}
	DS.mainLogBuf = bufio.NewWriter(tmpF1)
	DS.mainLogBuf.Flush()
	DS.mainLog = log.New(DS.mainLogBuf, "", log.LstdFlags)

	tmpF2, err2 := os.Create(DS.logDir + "main:eqns.log")
	if err2 != nil {
		log.Fatal("couldn't create eqns log", err2)
	}
	DS.eqnsLogBuf = bufio.NewWriter(tmpF2)
	DS.eqnsLogBuf.Flush()
	DS.eqnsLog = log.New(DS.eqnsLogBuf, "", log.LstdFlags)

}
