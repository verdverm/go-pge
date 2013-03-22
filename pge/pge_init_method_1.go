package pge

import (
	"fmt"

	probs "github.com/verdverm/go-pge/problems"
	expr "github.com/verdverm/go-symexpr"
)

func (PS *PgeSearch) GenInitExprMethod1() *probs.ReportQueue {
	fmt.Printf("generating initial expressions\n")

	GP := PS.cnfg.treecfg
	fmt.Printf("%v\n", GP)

	eList := make([]expr.Expr, 0)

	for _, T := range GP.Roots {
		switch T.ExprType() {
		case expr.ADD:
			eList = append(eList, PS.GenInitExprAddMethod1()[:]...)
		case expr.MUL:
			eList = append(eList, PS.GenInitExprMulMethod1()[:]...)
		case expr.DIV:
			eList = append(eList, PS.GenInitExprDivMethod1()[:]...)
		default:
			fmt.Printf("Error in GenInitExpr: unknown ROOT %d\n", T)
		}
	}

	exprs := probs.NewReportQueue()
	// exprs.SetSort(PS.cnfg.sortType)
	exprs.SetSort(probs.PESORT_PARETO_TST_ERR)

	for i, e := range eList {
		fmt.Printf("%d:  %v\n", i, e)
		serial := make([]int, 0, 64)
		serial = e.Serial(serial)
		PS.Trie.InsertSerial(serial)
		// on train data
		re := RegressExpr(e, PS.prob)
		re.SetUnitID(i)
		exprs.Push(re)
	}
	exprs.Sort()
	return exprs
}

func (PS *PgeSearch) GenInitExprAddMethod1() []expr.Expr {

	exprs := make([]expr.Expr, 0)

	// WARNING adding a single coefficient results in the same traversal but with the added extra coefficient
	cf := new(expr.Constant)
	cf.P = 1
	af := new(expr.Add)
	af.Insert(cf)
	exprs = append(exprs, af)

	for _, i := range PS.cnfg.treecfg.UsableVars {
		c := new(expr.Constant)
		c.P = -1
		v := new(expr.Var)
		v.P = i

		m := expr.NewMul()
		m.Insert(c)
		m.Insert(v)

		a := expr.NewAdd()
		a.Insert(m)
		exprs = append(exprs, a)
	}

	fmt.Println("Initial Add:  ", exprs)
	return exprs
}

func (PS *PgeSearch) GenInitExprMulMethod1() []expr.Expr {
	exprs := make([]expr.Expr, 0)

	for _, i := range PS.cnfg.treecfg.UsableVars {
		c := new(expr.Constant)
		c.P = -1
		v := new(expr.Var)
		v.P = i

		m := expr.NewMul()
		m.Insert(c)
		m.Insert(v)

		exprs = append(exprs, m)
	}

	fmt.Println("Initial Mul:  ", exprs)
	return exprs
}

func (PS *PgeSearch) GenInitExprDivMethod1() []expr.Expr {
	exprs := make([]expr.Expr, 0)

	for _, i := range PS.cnfg.treecfg.UsableVars {
		c := new(expr.Constant)
		c.P = -1
		v := new(expr.Var)
		v.P = i

		d := new(expr.Div)
		d.Numer = c
		d.Denom = v

		exprs = append(exprs, d)
	}

	fmt.Println("Initial Div:  ", exprs)
	return exprs
}
