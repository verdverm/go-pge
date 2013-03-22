package problems

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	expr "github.com/verdverm/go-symexpr"
)

type TreeParams struct {

	// bounds on tree
	MaxSize, MaxDepth,
	MinSize, MinDepth int

	// usable terms at each location
	RootsS, NodesS, LeafsS, NonTrigS []string
	RootsT, NodesT, LeafsT, NonTrigT []expr.ExprType
	Roots, Nodes, Leafs, NonTrig     []expr.Expr

	// simplify options
	DoSimp bool
	SRules expr.SimpRules

	// bounds on some operands
	UsableVars               []int
	NumDim, NumSys, NumCoeff int

	// tpm bounds on tree (for subtree distributions)
	TmpMaxSize, TmpMaxDepth,
	TmpMinSize, TmpMinDepth int

	// Current values
	CurrSize, CurrDepth int
	InTrig              bool
	CoeffCount          int
}

func (t *TreeParams) Clone() *TreeParams {
	n := new(TreeParams)
	n.MaxSize = t.MaxSize
	n.MaxDepth = t.MaxDepth
	n.MinSize = t.MinSize
	n.MinDepth = t.MinDepth

	n.RootsS = make([]string, len(t.RootsS))
	copy(n.RootsS, t.RootsS)
	n.NodesS = make([]string, len(t.NodesS))
	copy(n.NodesS, t.NodesS)
	n.LeafsS = make([]string, len(t.LeafsS))
	copy(n.LeafsS, t.LeafsS)
	n.NonTrigS = make([]string, len(t.NonTrigS))
	copy(n.NonTrigS, t.NonTrigS)

	n.RootsT = make([]expr.ExprType, len(t.RootsT))
	copy(n.RootsT, t.RootsT)
	n.NodesT = make([]expr.ExprType, len(t.NodesT))
	copy(n.NodesT, t.NodesT)
	n.LeafsT = make([]expr.ExprType, len(t.LeafsT))
	copy(n.LeafsT, t.LeafsT)
	n.NonTrigT = make([]expr.ExprType, len(t.NonTrigT))
	copy(n.NonTrigT, t.NonTrigT)

	n.Roots = make([]expr.Expr, len(t.Roots))
	copy(n.Roots, t.Roots)
	n.Nodes = make([]expr.Expr, len(t.Nodes))
	copy(n.Nodes, t.Nodes)
	n.Leafs = make([]expr.Expr, len(t.Leafs))
	copy(n.Leafs, t.Leafs)
	n.NonTrig = make([]expr.Expr, len(t.NonTrig))
	copy(n.NonTrig, t.NonTrig)

	n.DoSimp = t.DoSimp
	n.SRules = t.SRules

	n.UsableVars = make([]int, len(t.UsableVars))
	copy(n.UsableVars, t.UsableVars)

	n.NumDim = t.NumDim
	n.NumSys = t.NumSys
	n.NumCoeff = t.NumCoeff

	return n
}
func ParseTreeParams(field, value string, config interface{}) (found bool, err error) {

	TP := config.(*TreeParams)
	found = true

	switch strings.ToUpper(field) {
	case "ROOTS":
		TP.RootsS = strings.Fields(value)
		TP.RootsT, TP.Roots = fillExprStuff(TP.RootsS)
	case "NODES":
		TP.NodesS = strings.Fields(value)
		TP.NodesT, TP.Nodes = fillExprStuff(TP.NodesS)
	case "NONTRIG":
		TP.NonTrigS = strings.Fields(value)
		TP.NonTrigT, TP.NonTrig = fillExprStuff(TP.NonTrigS)
	case "LEAFS":
		TP.LeafsS = strings.Fields(value)
		TP.LeafsT, TP.Leafs = fillExprStuff(TP.LeafsS)

	case "USABLEVARS":
		usable := strings.Fields(value)
		for _, v := range usable {
			ival, cerr := strconv.Atoi(v)
			if cerr != nil {
				log.Printf("Expected integer for UsableDim\n")
				return found, cerr
			}
			TP.UsableVars = append(TP.UsableVars, ival)
		}

	case "MAXSIZE":
		ival, cerr := strconv.Atoi(value)
		if cerr != nil {
			log.Printf("Expected integer for MaxSize\n")
			return found, cerr
		}
		TP.MaxSize = ival
	case "MINSIZE":
		ival, cerr := strconv.Atoi(value)
		if cerr != nil {
			log.Printf("Expected integer for UsableDim\n")
			return found, cerr
		}
		TP.MinSize = ival
	case "MAXDEPTH":
		ival, cerr := strconv.Atoi(value)
		if cerr != nil {
			log.Printf("Expected integer for UsableDim\n")
			return found, cerr
		}
		TP.MaxDepth = ival
	case "MINDEPTH":
		ival, cerr := strconv.Atoi(value)
		if cerr != nil {
			log.Printf("Expected integer for UsableDim\n")
			return found, cerr
		}
		TP.MinDepth = ival

	default:
		found = false
		log.Printf("Problem Config Not Implemented: %s, %s\n\n", field, value)

	}
	return
}

func fillExprStuff(names []string) (types []expr.ExprType, exprs []expr.Expr) {
	types = make([]expr.ExprType, len(names))
	exprs = make([]expr.Expr, len(names))
	for i, n := range names {
		switch strings.ToLower(n) {
		case "constant":
			types[i] = expr.CONSTANT
			exprs[i] = new(expr.Constant)
		case "constantf":
			types[i] = expr.CONSTANTF
			exprs[i] = new(expr.ConstantF)
		case "time":
			types[i] = expr.TIME
			exprs[i] = new(expr.Time)
		case "system":
			types[i] = expr.SYSTEM
			exprs[i] = new(expr.System)
		case "var":
			types[i] = expr.VAR
			exprs[i] = new(expr.Var)

		case "neg":
			types[i] = expr.NEG
			exprs[i] = new(expr.Neg)
		case "abs":
			types[i] = expr.ABS
			exprs[i] = new(expr.Abs)
		case "sqrt":
			types[i] = expr.SQRT
			exprs[i] = new(expr.Sqrt)
		case "sin":
			types[i] = expr.SIN
			exprs[i] = new(expr.Sin)
		case "cos":
			types[i] = expr.COS
			exprs[i] = new(expr.Cos)
		case "tan":
			types[i] = expr.TAN
			exprs[i] = new(expr.Tan)
		case "exp":
			types[i] = expr.EXP
			exprs[i] = new(expr.Exp)
		case "log":
			types[i] = expr.LOG
			exprs[i] = new(expr.Log)

		case "powi":
			types[i] = expr.POWI
			exprs[i] = new(expr.PowI)
		case "powf":
			types[i] = expr.POWF
			exprs[i] = new(expr.PowF)
		case "powe":
			types[i] = expr.POWE
			exprs[i] = new(expr.PowE)

		case "add":
			types[i] = expr.ADD
			exprs[i] = expr.NewAdd()
		case "mul":
			types[i] = expr.MUL
			exprs[i] = expr.NewMul()
		case "div":
			types[i] = expr.DIV
			exprs[i] = new(expr.Div)
		default:
			log.Fatalf("Unknown ExprType:  %s\n", n)
		}
	}
	return
}

func (tp *TreeParams) CheckExprTmp(e expr.Expr) bool {
	if e.Size() < tp.TmpMinSize || e.Size() > tp.TmpMaxSize ||
		e.Height() < tp.TmpMinDepth || e.Height() > tp.TmpMaxDepth {
		return false
	}
	return true
}
func (tp *TreeParams) CheckExpr(e expr.Expr) bool {
	if e.Size() < tp.MinSize {
		//    fmt.Printf( "Too SMALL:  e:%v  l:%v\n", e.Size(), tp.TmpMinSize )
		return false
	} else if e.Size() > tp.MaxSize {
		//    fmt.Printf( "Too LARGE:  e:%v  l:%v\n", e.Size(), tp.TmpMaxSize )
		return false
	} else if e.Height() < tp.MinDepth {
		//    fmt.Printf( "Too SHORT:  e:%v  l:%v\n", e.Height(), tp.TmpMinDepth )
		return false
	} else if e.Height() > tp.MaxDepth {
		//    fmt.Printf( "Too TALL:  e:%v  l:%v\n", e.Height(), tp.TmpMaxDepth )
		return false
	}
	return true
}

func (tp *TreeParams) CheckExprLog(e expr.Expr, log *bufio.Writer) bool {
	//   if e.Size() < tp.TmpMinSize || e.Size() > tp.TmpMaxSize ||
	//     e.Height() < tp.TmpMinDepth || e.Height() > tp.TmpMaxDepth {
	//       return false
	//     }
	if e.Size() < tp.MinSize {
		fmt.Fprintf(log, "Too SMALL:  e:%v  l:%v\n", e.Size(), tp.MinSize)
		return false
	} else if e.Size() > tp.MaxSize {
		fmt.Fprintf(log, "Too LARGE:  e:%v  l:%v\n", e.Size(), tp.MaxSize)
		return false
	} else if e.Height() < tp.MinDepth {
		fmt.Fprintf(log, "Too SHORT:  e:%v  l:%v\n", e.Height(), tp.MinDepth)
		return false
	} else if e.Height() > tp.MaxDepth {
		fmt.Fprintf(log, "Too TALL:  e:%v  l:%v\n", e.Height(), tp.MaxDepth)
		return false
	}
	return true
}
func (tp *TreeParams) CheckExprPrint(e expr.Expr) bool {
	//   if e.Size() < tp.TmpMinSize || e.Size() > tp.TmpMaxSize ||
	//     e.Height() < tp.TmpMinDepth || e.Height() > tp.TmpMaxDepth {
	//       return false
	//     }
	if e.Size() < tp.MinSize {
		fmt.Printf("Too SMALL:  e:%v  l:%v\n", e.Size(), tp.MinSize)
		return false
	} else if e.Size() > tp.MaxSize {
		fmt.Printf("Too LARGE:  e:%v  l:%v\n", e.Size(), tp.MaxSize)
		return false
	} else if e.Height() < tp.MinDepth {
		fmt.Printf("Too SHORT:  e:%v  l:%v\n", e.Height(), tp.MinDepth)
		return false
	} else if e.Height() > tp.MaxDepth {
		fmt.Printf("Too TALL:  e:%v  l:%v\n", e.Height(), tp.MaxDepth)
		return false
	}
	return true
}

func (tp *TreeParams) ResetCurr() {
	tp.CurrSize, tp.CurrDepth, tp.InTrig, tp.CoeffCount = 0, 0, false, 0
}
func (tp *TreeParams) ResetTemp() {
	tp.TmpMaxSize, tp.TmpMaxDepth = tp.MaxSize, tp.MaxDepth
	tp.TmpMinSize, tp.TmpMinDepth = tp.MinSize, tp.MinDepth
}
