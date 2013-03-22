package problems

import (
	"log"
	"sort"
	"strconv"
	"strings"

	expr "github.com/verdverm/go-symexpr"
)

type ExprProblemType int

const (
	UnknownPType ExprProblemType = iota
	ExprBenchmark
	ExprDiffeq
	ExprRealData
)

func (ept ExprProblemType) String() string {
	switch ept {
	case ExprBenchmark:
		return "benchmark"
	case ExprDiffeq:
		return "diffeq"
	case ExprRealData:
		return "real"
	default:
		return "UnknownPType"
	}
	return ""
}

func ProblemTypeFromString(ptype string) ExprProblemType {
	switch ptype {
	case "benchmark":
		return ExprBenchmark
	case "diffeq":
		return ExprDiffeq
	case "real":
		return ExprRealData
	default:
		return UnknownPType
	}
	return UnknownPType
}

type ExprProblem struct {
	Name     string
	VarNames []string
	FuncTree expr.Expr // function as tree
	MaxIter  int

	// type of evaluation / data
	SearchType ExprProblemType

	// data
	TrainFns []string
	TestFns  []string
	Train    []*PointSet
	Test     []*PointSet
	HitRatio float64

	// variable information
	SearchVar  int
	UsableVars []int
	IndepNames []string
	DepndNames []string
	SysNames   []string

	// tree gen/restrict information
	TreeCfg *TreeParams
}

type ExprProblemComm struct {
	// incoming channels
	Cmds chan int

	// outgoing channels
	Rpts chan *ExprReportArray
	Gen  chan [2]int
}

func ProbConfigParser(field, value string, config interface{}) (err error) {

	EP := config.(*ExprProblem)
	if EP.TreeCfg == nil {
		EP.TreeCfg = new(TreeParams)
	}

	switch strings.ToUpper(field) {
	case "NAME":
		EP.Name = value
	case "PROBLEMTYPE":
		typ := ProblemTypeFromString(strings.ToLower(value))
		if typ == UnknownPType {
			log.Fatalf("Unknown ProblemType in Problem config file\n")
		} else {
			EP.SearchType = typ
		}
	case "MAXITER":
		EP.MaxIter, err = strconv.Atoi(value)

	case "HITRATIO":
		fval, cerr := strconv.ParseFloat(value, 64)
		if cerr != nil {
			log.Printf("Expected float64 for HitRatio\n")
			return cerr
		}
		EP.HitRatio = fval

	case "TRAINDATA":
		EP.TrainFns = strings.Fields(value)
	case "TESTDATA":
		EP.TestFns = strings.Fields(value)

	case "USABLEVARS":
		usable := strings.Fields(value)
		for _, v := range usable {
			ival, cerr := strconv.Atoi(v)
			if cerr != nil {
				log.Printf("Expected integer for UsableDim\n")
				return cerr
			}
			EP.UsableVars = append(EP.UsableVars, ival)
		}
		EP.UsableVars = unique(EP.UsableVars)
		EP.TreeCfg.UsableVars = EP.UsableVars[:]
		EP.TreeCfg.UsableVars = unique(EP.TreeCfg.UsableVars)

	case "SEARCHVAR":
		ival, cerr := strconv.Atoi(value)
		if cerr != nil {
			log.Printf("Expected integer for UsableDim\n")
			return cerr
		}
		EP.SearchVar = ival

	default:
		// check augillary parsable structures [only TreeParams for now]
		found, ferr := ParseTreeParams(field, value, EP.TreeCfg)
		if ferr != nil {
			log.Fatalf("error parsing Problem Config\n")
			return ferr
		}
		if !found {
			log.Printf("Problem Config Not Implemented: %s, %s\n\n", field, value)
		}

	}
	return
}

func unique(list []int) []int {
	sort.Ints(list)
	var last int
	i := 0
	for _, x := range list {
		if i == 0 || x != last {
			last = x
			list[i] = x
			i++
		}
	}
	return list[0:i]
}
