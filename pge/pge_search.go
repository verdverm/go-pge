package pge

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	levmar "github.com/verdverm/go-levmar"
	config "github.com/verdverm/go-pge/config"
	probs "github.com/verdverm/go-pge/problems"
	expr "github.com/verdverm/go-symexpr"
)

type pgeConfig struct {
	// search params
	maxGen        int
	pgeRptEpoch   int
	pgeRptCount   int
	pgeArchiveCap int

	simprules expr.SimpRules
	treecfg   *probs.TreeParams

	// PGE specific options
	peelCnt     int
	sortType    probs.SortType
	zeroEpsilon float64

	initMethod string
	growMethod string

	evalrCount int
}

func pgeConfigParser(field, value string, config interface{}) (err error) {

	PC := config.(*pgeConfig)

	switch strings.ToUpper(field) {
	case "MAXGEN":
		PC.maxGen, err = strconv.Atoi(value)
	case "PGERPTEPOCH":
		PC.pgeRptEpoch, err = strconv.Atoi(value)
	case "PGERPTCOUNT":
		PC.pgeRptCount, err = strconv.Atoi(value)
	case "PGEARCHIVECAP":
		PC.pgeArchiveCap, err = strconv.Atoi(value)

	case "PEELCOUNT":
		PC.peelCnt, err = strconv.Atoi(value)

	case "EVALRCOUNT":
		PC.evalrCount, err = strconv.Atoi(value)

	case "SORTTYPE":
		switch strings.ToLower(value) {
		case "paretotrainerror":
			PC.sortType = probs.PESORT_PARETO_TRN_ERR
		case "paretotesterror":
			PC.sortType = probs.PESORT_PARETO_TST_ERR

		default:
			log.Printf("PGE Config Not Implemented: %s, %s\n\n", field, value)
		}

	case "ZEROEPSILON":
		PC.zeroEpsilon, err = strconv.ParseFloat(value, 64)

	default:
		// check augillary parsable structures [only TreeParams for now]
		if PC.treecfg == nil {
			PC.treecfg = new(probs.TreeParams)
		}
		found, ferr := probs.ParseTreeParams(field, value, PC.treecfg)
		if ferr != nil {
			log.Fatalf("error parsing PGE - treecfg Config\n")
			return ferr
		}
		if !found {
			log.Printf("PGE Config Not Implemented: %s, %s\n\n", field, value)
		}

	}
	return
}

type PgeSearch struct {
	id   int
	cnfg pgeConfig
	prob *probs.ExprProblem
	iter int
	stop bool

	// comm up
	commup *probs.ExprProblemComm

	// comm down

	// best exprs
	Best *probs.ReportQueue

	// training data in C format
	c_input  []levmar.C_double
	c_ygiven []levmar.C_double

	// logs
	logDir     string
	mainLog    *log.Logger
	mainLogBuf *bufio.Writer
	eqnsLog    *log.Logger
	eqnsLogBuf *bufio.Writer
	errLog     *log.Logger
	errLogBuf  *bufio.Writer

	fitnessLog    *log.Logger
	fitnessLogBuf *bufio.Writer
	ipreLog       *log.Logger
	ipreLogBuf    *bufio.Writer

	// equations visited
	Trie  *IpreNode
	Queue *probs.ReportQueue

	// eval channels
	eval_in  chan expr.Expr
	eval_out chan *probs.ExprReport

	// genStuff
	GenRoots   []expr.Expr
	GenLeafs   []expr.Expr
	GenNodes   []expr.Expr
	GenNonTrig []expr.Expr

	// FFXish stuff
	ffxBases []expr.Expr

	// statistics
	neqns    int
	ipreCnt  int
	maxSize  int
	maxScore int
	minError float64
}

func (PS *PgeSearch) GetMaxIter() int {
	return PS.cnfg.maxGen
}
func (PS *PgeSearch) SetMaxIter(iter int) {
	PS.cnfg.maxGen = iter
}
func (PS *PgeSearch) SetPeelCount(cnt int) {
	PS.cnfg.peelCnt = cnt
}
func (PS *PgeSearch) SetInitMethod(init string) {
	PS.cnfg.initMethod = init
}
func (PS *PgeSearch) SetGrowMethod(grow string) {
	PS.cnfg.growMethod = grow
}
func (PS *PgeSearch) SetEvalrCount(cnt int) {
	PS.cnfg.evalrCount = cnt
}

func (PS *PgeSearch) ParseConfig(filename string) {
	fmt.Printf("Parsing PGE Config: %s\n", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = config.ParseConfig(data, pgeConfigParser, &PS.cnfg)
	if err != nil {
		log.Fatal(err)
	}
}

func (PS *PgeSearch) Init(done chan int, prob *probs.ExprProblem, logdir string, input interface{}) {
	fmt.Printf("Init'n PGE\n")
	// setup data

	// open logs
	PS.initLogs(logdir)

	// copy in common config options
	PS.prob = prob
	if PS.cnfg.treecfg == nil {
		PS.cnfg.treecfg = PS.prob.TreeCfg.Clone()
	}
	srules := expr.DefaultRules()
	srules.ConvertConsts = true
	PS.cnfg.simprules = srules

	fmt.Println("Roots:   ", PS.cnfg.treecfg.RootsS)
	fmt.Println("Nodes:   ", PS.cnfg.treecfg.NodesS)
	fmt.Println("Leafs:   ", PS.cnfg.treecfg.LeafsS)
	fmt.Println("NonTrig: ", PS.cnfg.treecfg.NonTrigS)

	PS.GenRoots = make([]expr.Expr, len(PS.cnfg.treecfg.Roots))
	for i := 0; i < len(PS.GenRoots); i++ {
		PS.GenRoots[i] = PS.cnfg.treecfg.Roots[i].Clone()
	}
	PS.GenNodes = make([]expr.Expr, len(PS.cnfg.treecfg.Nodes))
	for i := 0; i < len(PS.GenNodes); i++ {
		PS.GenNodes[i] = PS.cnfg.treecfg.Nodes[i].Clone()
	}
	PS.GenNonTrig = make([]expr.Expr, len(PS.cnfg.treecfg.NonTrig))
	for i := 0; i < len(PS.GenNonTrig); i++ {
		PS.GenNonTrig[i] = PS.cnfg.treecfg.NonTrig[i].Clone()
	}

	PS.GenLeafs = make([]expr.Expr, 0)
	for _, t := range PS.cnfg.treecfg.LeafsT {
		switch t {
		case expr.TIME:
			PS.GenLeafs = append(PS.GenLeafs, expr.NewTime())

		case expr.VAR:
			fmt.Println("Use Vars: ", PS.cnfg.treecfg.UsableVars)
			for _, i := range PS.cnfg.treecfg.UsableVars {
				PS.GenLeafs = append(PS.GenLeafs, expr.NewVar(i))
			}

		case expr.SYSTEM:
			for i := 0; i < PS.prob.Train[0].NumSys(); i++ {
				PS.GenLeafs = append(PS.GenLeafs, expr.NewSystem(i))
			}

		}
	}
	/*** FIX ME
	PS.GenLeafs = make([]expr.Expr, len(PS.cnfg.treecfg.Leafs))
	for i := 0; i < len(PS.GenLeafs); i++ {
		PS.GenLeafs[i] = PS.cnfg.treecfg.Leafs[i].Clone()
	}
	***/

	fmt.Println("Roots:   ", PS.GenRoots)
	fmt.Println("Nodes:   ", PS.GenNodes)
	fmt.Println("Leafs:   ", PS.GenLeafs)
	fmt.Println("NonTrig: ", PS.GenNonTrig)

	// setup communication struct
	PS.commup = input.(*probs.ExprProblemComm)

	// initialize bbq
	PS.Trie = new(IpreNode)
	PS.Trie.val = -1
	PS.Trie.next = make(map[int]*IpreNode)

	PS.Best = probs.NewReportQueue()
	PS.Best.SetSort(probs.GPSORT_PARETO_TST_ERR)
	PS.Queue = PS.GenInitExpr()
	PS.Queue.SetSort(probs.PESORT_PARETO_TST_ERR)

	PS.neqns = PS.Queue.Len()

	PS.minError = math.Inf(1)

	PS.eval_in = make(chan expr.Expr, 4048)
	PS.eval_out = make(chan *probs.ExprReport, 4048)

	for i := 0; i < PS.cnfg.evalrCount; i++ {
		go PS.Evaluate()
	}
}

func (PS *PgeSearch) Evaluate() {

	for !PS.stop {
		e := <-PS.eval_in
		if e == nil {
			continue
		}
		PS.eval_out <- RegressExpr(e, PS.prob)
	}

}

func (PS *PgeSearch) Run() {
	fmt.Printf("Running PGE\n")

	PS.loop()

	fmt.Println("PGE exitting")

	PS.Clean()
	PS.commup.Cmds <- -1
}

func (PS *PgeSearch) loop() {

	PS.checkMessages()
	for !PS.stop {

		fmt.Println("in: PS.step() ", PS.iter)
		PS.step()

		// if PS.iter%PS.cnfg.pgeRptEpoch == 0 {
		PS.reportExpr()
		// }

		// report current iteration
		PS.commup.Gen <- [2]int{PS.id, PS.iter}
		PS.iter++

		PS.Clean()

		PS.checkMessages()

	}

	// done expanding, pull the rest of the regressed solutions from the queue
	p := 0
	for PS.Queue.Len() > 0 {
		e := PS.Queue.Pop().(*probs.ExprReport)

		bPush := true
		if len(e.Coeff()) == 1 && math.Abs(e.Coeff()[0]) < PS.cnfg.zeroEpsilon {
			// fmt.Println("No Best Push")
			bPush = false
		}

		if bPush {
			// fmt.Printf("pop/push(%d,%d): %v\n", p, PS.Best.Len(), e.Expr())
			PS.Best.Push(e)
			p++
		}

		if e.TestScore() > PS.maxScore {
			PS.maxScore = e.TestScore()
		}
		if e.TestError() < PS.minError {
			PS.minError = e.TestError()
			fmt.Printf("EXITING New Min Error:  %v\n", e)
		}
		if e.Size() > PS.maxSize {
			PS.maxSize = e.Size()
		}
	}

	fmt.Println("PGE sending last report")
	PS.reportExpr()

}

func (PS *PgeSearch) step() {

	loop := 0
	eval_cnt := 0 // for channeled eval

	es := PS.peel()

	ex := PS.expandPeeled(es)

	for cnt := range ex {
		E := ex[cnt]

		if E == nil {
			continue
		}

		for _, e := range E {
			if e == nil {
				continue
			}
			if !PS.cnfg.treecfg.CheckExpr(e) {
				continue
			}

			// check ipre_trie
			serial := make([]int, 0, 64)
			serial = e.Serial(serial)
			ins := PS.Trie.InsertSerial(serial)
			if !ins {
				continue
			}

			// for serial eval
			// re := RegressExpr(e, PS.prob)

			// start channeled eval
			PS.eval_in <- e
			eval_cnt++
		}
	}
	for i := 0; i < eval_cnt; i++ {
		re := <-PS.eval_out
		// end channeled eval

		// check for NaN/Inf in re.error  and  if so, skip
		if math.IsNaN(re.TestError()) || math.IsInf(re.TestError(), 0) {
			// fmt.Printf("Bad Error\n%v\n", re)
			continue
		}

		if re.TestError() < PS.minError {
			PS.minError = re.TestError()
		}

		// check for coeff == 0
		doIns := true
		for _, c := range re.Coeff() {
			// i > 0 for free coeff
			if math.Abs(c) < PS.cnfg.zeroEpsilon {
				doIns = false
				break
			}
		}

		if doIns {
			re.SetProcID(PS.id)
			re.SetIterID(PS.iter)
			re.SetUnitID(loop)
			re.SetUniqID(PS.neqns)
			loop++
			PS.neqns++
			// fmt.Printf("Queue.Push(): %v\n%v\n\n", re.Expr(), serial)
			// fmt.Printf("Queue.Push(): %v\n", re)
			// fmt.Printf("Queue.Push(): %v\n", re.Expr())

			PS.Queue.Push(re)

		}
	}
	// } // for sequential eval
	PS.Queue.Sort()

}

func (PS *PgeSearch) peel() []*probs.ExprReport {
	es := make([]*probs.ExprReport, PS.cnfg.peelCnt)
	for p := 0; p < PS.cnfg.peelCnt && PS.Queue.Len() > 0; p++ {

		e := PS.Queue.Pop().(*probs.ExprReport)

		bPush := true
		if len(e.Coeff()) == 1 && math.Abs(e.Coeff()[0]) < PS.cnfg.zeroEpsilon {
			fmt.Println("No Best Push")
			p--
			continue
		}

		if bPush {
			fmt.Printf("pop/push(%d,%d): %v\n", p, PS.Best.Len(), e.Expr())
			PS.Best.Push(e)
		}

		es[p] = e

		if e.TestScore() > PS.maxScore {
			PS.maxScore = e.TestScore()
		}
		if e.TestError() < PS.minError {
			PS.minError = e.TestError()
			fmt.Printf("Best New Min Error:  %v\n", e)
		}
		if e.Size() > PS.maxSize {
			PS.maxSize = e.Size()
		}

	}
	return es
}

func (PS *PgeSearch) expandPeeled(es []*probs.ExprReport) [][]expr.Expr {
	eqns := make([][]expr.Expr, PS.cnfg.peelCnt)
	for p := 0; p < PS.cnfg.peelCnt; p++ {
		if es[p] == nil {
			continue
		}
		// fmt.Printf("expand(%d): %v\n", p, es[p].Expr())
		if es[p].Expr().ExprType() != expr.ADD {
			add := expr.NewAdd()
			add.Insert(es[p].Expr())
			add.CalcExprStats()
			es[p].SetExpr(add)
		}
		eqns[p] = PS.Expand(es[p].Expr())
		// fmt.Printf("Results:\n")
		// for i, e := range eqns[p] {
		// 	fmt.Printf("%d,%d:  %v\n", p, i, e)
		// }
		// fmt.Println()
	}
	fmt.Println("\n")
	return eqns
}

func (PS *PgeSearch) reportExpr() {

	cnt := PS.cnfg.pgeRptCount
	PS.Best.Sort()

	// repot best equations
	rpt := make(probs.ExprReportArray, cnt)
	if PS.Best.Len() < cnt {
		cnt = PS.Best.Len()
	}
	copy(rpt, PS.Best.GetQueue()[:cnt])

	errSum, errCnt := 0.0, 0
	PS.eqnsLog.Println("\n\nReport", PS.iter)
	for i, r := range rpt {
		PS.eqnsLog.Printf("\n%d:  %v\n", i, r)
		if r != nil && r.Expr() != nil {
			errSum += r.TestError()
			errCnt++
		}
	}

	PS.mainLog.Printf("Iter: %d  %f  %f\n", PS.iter, errSum/float64(errCnt), PS.minError)

	PS.ipreLog.Println(PS.iter, PS.neqns, PS.Trie.cnt, PS.Trie.vst)
	PS.fitnessLog.Println(PS.iter, PS.neqns, PS.Trie.cnt, PS.Trie.vst, errSum/float64(errCnt), PS.minError)

	PS.commup.Rpts <- &rpt

}

func (PS *PgeSearch) Clean() {
	// fmt.Printf("Cleaning PGE\n")

	PS.errLogBuf.Flush()
	PS.mainLogBuf.Flush()
	PS.eqnsLogBuf.Flush()
	PS.fitnessLogBuf.Flush()
	PS.ipreLogBuf.Flush()

}

func (PS *PgeSearch) initLogs(logdir string) {
	// open logs
	PS.logDir = logdir + "pge/"
	os.Mkdir(PS.logDir, os.ModePerm)
	tmpF0, err5 := os.Create(PS.logDir + "pge:err.log")
	if err5 != nil {
		log.Fatal("couldn't create errs log")
	}
	PS.errLogBuf = bufio.NewWriter(tmpF0)
	PS.errLogBuf.Flush()
	PS.errLog = log.New(PS.errLogBuf, "", log.LstdFlags)

	tmpF1, err1 := os.Create(PS.logDir + "pge:main.log")
	if err1 != nil {
		log.Fatal("couldn't create main log")
	}
	PS.mainLogBuf = bufio.NewWriter(tmpF1)
	PS.mainLogBuf.Flush()
	PS.mainLog = log.New(PS.mainLogBuf, "", log.LstdFlags)

	tmpF2, err2 := os.Create(PS.logDir + "pge:eqns.log")
	if err2 != nil {
		log.Fatal("couldn't create eqns log")
	}
	PS.eqnsLogBuf = bufio.NewWriter(tmpF2)
	PS.eqnsLogBuf.Flush()
	PS.eqnsLog = log.New(PS.eqnsLogBuf, "", 0)

	tmpF3, err3 := os.Create(PS.logDir + "pge:fitness.log")
	if err3 != nil {
		log.Fatal("couldn't create eqns log")
	}
	PS.fitnessLogBuf = bufio.NewWriter(tmpF3)
	PS.fitnessLogBuf.Flush()
	PS.fitnessLog = log.New(PS.fitnessLogBuf, "", log.Ltime|log.Lmicroseconds)

	tmpF4, err4 := os.Create(PS.logDir + "pge:ipre.log")
	if err4 != nil {
		log.Fatal("couldn't create eqns log")
	}
	PS.ipreLogBuf = bufio.NewWriter(tmpF4)
	PS.ipreLogBuf.Flush()
	PS.ipreLog = log.New(PS.ipreLogBuf, "", log.Ltime|log.Lmicroseconds)
}

func (PS *PgeSearch) checkMessages() {

	// check messages from superior
	select {
	case cmd, ok := <-PS.commup.Cmds:
		if ok {
			if cmd == -1 {
				fmt.Println("PGE: stop sig recv'd")
				PS.stop = true
				return
			}
		}
	default:
		return
	}
}

var c_input, c_ygiven []levmar.C_double

func RegressExpr(E expr.Expr, P *probs.ExprProblem) (R *probs.ExprReport) {

	guess := make([]float64, 0)
	guess, eqn := E.ConvertToConstants(guess)

	var coeff []float64
	if len(guess) > 0 {

		// fmt.Printf("x_dims:  %d  %d\n", x_dim, x_dim2)

		// Callback version
		coeff = levmar.LevmarExpr(eqn, P.SearchVar, P.SearchType, guess, P.Train, P.Test)

		// Stack version
		// x_dim := P.Train[0].NumDim()
		// if c_input == nil {
		// 	ps := P.Train[0].NumPoints()
		// 	PS := len(P.Train) * ps
		// 	x_tot := PS * x_dim

		// 	c_input = make([]levmar.C_double, x_tot)
		// 	c_ygiven = make([]levmar.C_double, PS)

		// 	for i1, T := range P.Train {
		// 		for i2, p := range T.Points() {
		// 			i := i1*ps + i2
		// 			c_ygiven[i] = levmar.MakeCDouble(p.Depnd(P.SearchVar))
		// 			for i3, x_p := range p.Indeps() {
		// 				j := i1*ps*x_dim + i2*x_dim + i3
		// 				c_input[j] = levmar.MakeCDouble(x_p)
		// 			}
		// 		}
		// 	}
		// }
		// coeff = levmar.StackLevmarExpr(eqn, x_dim, guess, c_ygiven, c_input)

		// serial := make([]int, 0)
		// serial = eqn.StackSerial(serial)
		// fmt.Printf("StackSerial: %v\n", serial)
		// fmt.Printf("%v\n%v\n%v\n\n", eqn, coeff, steff)
	}

	R = new(probs.ExprReport)
	R.SetExpr(eqn) /*.ConvertToConstantFs(coeff)*/
	R.SetCoeff(coeff)
	R.Expr().CalcExprStats()

	// hitsL1, hitsL2, evalCnt, nanCnt, infCnt, l1_err, l2_err := scoreExpr(E, P, coeff)
	_, _, _, trnNanCnt, _, trn_l1_err, _ := scoreExpr(E, P, P.Train, coeff)
	_, _, tstEvalCnt, tstNanCnt, _, tst_l1_err, tst_l2_err := scoreExpr(E, P, P.Test, coeff)

	R.SetTrainScore(trnNanCnt)
	R.SetTrainError(trn_l1_err)

	R.SetPredScore(tstNanCnt)
	R.SetTestScore(tstEvalCnt)
	R.SetTestError(tst_l1_err)
	R.SetPredError(tst_l2_err)

	return R
}

func scoreExpr(e expr.Expr, P *probs.ExprProblem, dataSets []*probs.PointSet, coeff []float64) (hitsL1, hitsL2, evalCnt, nanCnt, infCnt int, l1_err, l2_err float64) {
	var l1_sum, l2_sum float64
	for _, PS := range dataSets {
		for _, p := range PS.Points() {
			y := p.Depnd(P.SearchVar)
			var out float64
			if P.SearchType == probs.ExprBenchmark {
				out = e.Eval(0, p.Indeps(), coeff, PS.SysVals())
			} else if P.SearchType == probs.ExprDiffeq {
				out = e.Eval(p.Indep(0), p.Indeps()[1:], coeff, PS.SysVals())
			}

			if math.IsNaN(out) {
				nanCnt++
				continue
			} else if math.IsInf(out, 0) {
				infCnt++
				continue
			} else {
				evalCnt++
			}

			diff := out - y
			l1_val := math.Abs(diff)
			l2_val := diff * diff
			l1_sum += l1_val
			l2_sum += l2_val

			if l1_val < P.HitRatio {
				hitsL1++
			}
			if l2_val < P.HitRatio {
				hitsL2++
			}
		}
	}

	if evalCnt == 0 {
		l1_err = math.NaN()
		l2_err = math.NaN()
	} else {
		fEvalCnt := float64(evalCnt + 1)
		l1_err = l1_sum / fEvalCnt
		l2_err = math.Sqrt(l2_sum / fEvalCnt)
	}

	return
}
