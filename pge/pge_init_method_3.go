package pge

import (
	"fmt"

	probs "github.com/verdverm/go-pge/problems"
	expr "github.com/verdverm/go-symexpr"
)

// This is the FFXish style init function
func (PS *PgeSearch) GenInitExprMethod3() *probs.ReportQueue {
	fmt.Printf("generating initial expressions\n")

	GP := PS.cnfg.treecfg
	fmt.Printf("%v\n", GP)

	bases := make([]expr.Expr, 0)

	// single constant
	bases = append(bases, expr.NewConstant(-1))

	// c*x_i
	for _, v := range PS.prob.UsableVars {
		mul := expr.NewMul()
		mul.Insert(expr.NewConstant(-1))
		mul.Insert(expr.NewVar(v))
		bases = append(bases, mul)
	}

	// c*x_i^p & c*x_i^-p
	for p := 2; p <= 4; p++ {
		for _, v := range PS.prob.UsableVars {
			// positive powers
			mul := expr.NewMul()
			mul.Insert(expr.NewConstant(-1))
			mul.Insert(expr.NewPowI(expr.NewVar(v), p))
			bases = append(bases, mul)

			// negative powers
			nul := expr.NewMul()
			nul.Insert(expr.NewConstant(-1))
			nul.Insert(expr.NewPowI(expr.NewVar(v), -p))
			bases = append(bases, nul)
		}
	}

	// c*N(d*x_i)
	for _, N := range PS.GenNodes {
		if N.ExprType() == expr.DIV || N.ExprType() == expr.ADD || N.ExprType() == expr.MUL {
			continue
		}
		for _, v := range PS.prob.UsableVars {
			// positive powers
			mul := expr.NewMul()
			mul.Insert(expr.NewConstant(-1))

			nul := expr.NewMul()
			nul.Insert(expr.NewConstant(-1))
			nul.Insert(expr.NewVar(v))

			n := N.Clone()
			p := 1
			n.SetExpr(&p, nul)
			mul.Insert(n)
			bases = append(bases, mul)
		}
	}

	// copy bases for use later in expand
	PS.ffxBases = make([]expr.Expr, len(bases))
	for i, b := range bases {
		PS.ffxBases[i] = b.Clone()
	}

	exprs := probs.NewReportQueue()
	// exprs.SetSort(PS.cnfg.sortType)
	exprs.SetSort(probs.PESORT_PARETO_TST_ERR)

	for i, e := range bases {
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
