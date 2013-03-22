package pge

import (
	"fmt"
	// "sort"

	probs "github.com/verdverm/go-pge/problems"
	expr "github.com/verdverm/go-symexpr"
)

func (PS *PgeSearch) GenInitExpr() *probs.ReportQueue {
	switch PS.cnfg.initMethod {
	case "method1":
		return PS.GenInitExprMethod1()
	case "method2":
		return PS.GenInitExprMethod2()
	case "method3":
		return PS.GenInitExprMethod3()
	default:
		fmt.Println("Unknown init method")
	}
	return nil
}

func (PS *PgeSearch) Expand(O expr.Expr) (ret []expr.Expr) {
	var exprs []expr.Expr

	switch PS.cnfg.growMethod {
	case "method1":
		exprs = PS.ExpandMethod1(O)
	case "method2":
		exprs = PS.ExpandMethod2(O)
	case "method3":
		exprs = PS.ExpandMethod3(O)
	default:
		fmt.Println("Unknown expand method")
	}

	// convert and simplify
	for i, e := range exprs {
		if e == nil {
			continue
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered in Expand %v   %d %v", r, i, e)
				exprs[i] = nil
			}
		}()
		c := make([]float64, 0)
		// fmt.Printf("Preconv   %v\n", e)
		c, eqn := e.ConvertToConstants(c)
		e.CalcExprStats()
		// fmt.Printf("Postconv  %v\n", eqn)

		// fmt.Printf("Presimp  %v\n", eqn)
		exprs[i] = eqn.Simplify(PS.cnfg.simprules)
		// fmt.Printf("Postsimp  %v\n\n", exprs[i])
		// serial := make([]int, 0, 64)
		// serial = exprs[i].Serial(serial)
		// fmt.Printf("Postsimp  %v   %v\n\n", exprs[i], serial)
	}

	// sort.Sort(expr.ExprArray(exprs))

	// remove duplicates

	// last := 0
	// for last < len(exprs) && exprs[last] == nil {
	// 	last++
	// }
	// for i := last + 1; i < len(exprs); i++ {
	// 	if exprs[i] == nil {
	// 		continue
	// 	}
	// 	if exprs[i].AmIAlmostSame(exprs[last]) {
	// 		exprs[i] = nil
	// 	} else {
	// 		last = i
	// 	}
	// }

	// copy to ret, ignoring nil & bad expressions ...
	for _, e := range exprs {
		if e == nil {
			continue
		}
		// if !PS.cnfg.treecfg.CheckExpr(e) {
		// 	continue
		// }
		ret = append(ret, e)
	}
	return ret
}
