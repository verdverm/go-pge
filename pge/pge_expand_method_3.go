package pge

import (
	// "fmt"
	"sort"

	expr "github.com/verdverm/go-symexpr"
)

func (PS *PgeSearch) ExpandMethod3(O expr.Expr) (ret []expr.Expr) {
	O.Sort()
	ret = make([]expr.Expr, 0)
	// fmt.Printf("Expanding expression:  %v\n", O)

	add := O.(*expr.Add)

	// adding term to addition
	for _, B := range PS.ffxBases {
		found := false
		for _, C := range add.CS {
			// fmt.Printf("checking %v in %v\n", B, add)
			if C.AmISame(B) {
				// fmt.Println("found\n\n")
				found = true
				break
			}
		}
		if !found {
			e := O.Clone()
			a := e.(*expr.Add)
			a.Insert(B.Clone())
			sort.Sort(a)
			a.CalcExprStats()
			good := PS.cnfg.treecfg.CheckExpr(a)
			if good {
				ret = append(ret, a)
			}
			// fmt.Printf("grew  %v\n\n", a)
			// ret = append(ret, a)
		}
	}

	// extending terms in addition
	for _, B := range PS.ffxBases {
		for i, C := range add.CS {
			if C.ExprType() == expr.MUL {
				m := C.(*expr.Mul)
				if len(m.CS) > 3 {
					continue
				}
			}
			e := O.Clone()
			a := e.(*expr.Add)
			mul := expr.NewMul()
			mul.Insert(a.CS[i])
			mul.Insert(B.Clone())
			a.CS[i] = mul
			sort.Sort(a)
			a.CalcExprStats()
			good := PS.cnfg.treecfg.CheckExpr(a)
			if good {
				ret = append(ret, a)
			}
			// fmt.Printf("grew  %v\n\n", a)
			ret = append(ret, a)
		}
	}

	// fmt.Println("Len of ret = ", len(ret))
	return ret
}
