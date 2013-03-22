package pge

import (
	// "fmt"
	"sort"

	expr "github.com/verdverm/go-symexpr"
)

func (PS *PgeSearch) ExpandMethod2(O expr.Expr) (ret []expr.Expr) {
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

	// deepening terms 
	// if len(add.CS) < 2 {
	// 	return ret
	// }
	for i, C := range add.CS {
		if C.ExprType() == expr.MUL {
			m := C.(*expr.Mul)
			if len(m.CS) != 2 {
				continue
			}
			if m.CS[1].ExprType() == expr.ADD {
				continue
			}
		} else {
			continue
		}

		for _, B := range PS.ffxBases {
			e := O.Clone()
			a := e.(*expr.Add)
			m := a.CS[i].(*expr.Mul)
			n := m.CS[1]

			switch n.ExprType() {
			case expr.SQRT:
				A := expr.NewAdd()
				A.Insert(n.(*expr.Sqrt).C)
				A.Insert(B.Clone())
				n.(*expr.Sqrt).C = A
			case expr.SIN:
				A := expr.NewAdd()
				A.Insert(n.(*expr.Sin).C)
				A.Insert(B.Clone())
				n.(*expr.Sin).C = A
			case expr.COS:
				A := expr.NewAdd()
				A.Insert(n.(*expr.Cos).C)
				A.Insert(B.Clone())
				n.(*expr.Cos).C = A
			case expr.TAN:
				A := expr.NewAdd()
				A.Insert(n.(*expr.Tan).C)
				A.Insert(B.Clone())
				n.(*expr.Tan).C = A
			case expr.EXP:
				A := expr.NewAdd()
				A.Insert(n.(*expr.Exp).C)
				A.Insert(B.Clone())
				n.(*expr.Exp).C = A
			case expr.LOG:
				A := expr.NewAdd()
				A.Insert(n.(*expr.Log).C)
				A.Insert(B.Clone())
				n.(*expr.Log).C = A

			default:
				A := expr.NewAdd()
				A.Insert(m.CS[1])
				A.Insert(B.Clone())
				m.CS[1] = A
			}
			sort.Sort(a)
			a.CalcExprStats()
			good := PS.cnfg.treecfg.CheckExpr(a)
			if good {
				ret = append(ret, a)
			}
			// fmt.Printf("grew  %v\n", a)
			ret = append(ret, a)
		}
	}
	for i, C := range add.CS {
		if C.ExprType() == expr.MUL {
			m := C.(*expr.Mul)
			if len(m.CS) != 2 {
				continue
			}
			if m.CS[1].ExprType() == expr.ADD {
				continue
			}
		} else {
			continue
		}

		for _, B := range PS.ffxBases {
			e := O.Clone()
			a := e.(*expr.Add)
			m := a.CS[i].(*expr.Mul)
			n := m.CS[1]

			switch n.ExprType() {
			case expr.SQRT:
				M := expr.NewMul()
				M.Insert(n.(*expr.Sqrt).C)
				M.Insert(B.Clone())
				n.(*expr.Sqrt).C = M
			case expr.SIN:
				M := expr.NewMul()
				M.Insert(n.(*expr.Sin).C)
				M.Insert(B.Clone())
				n.(*expr.Sin).C = M
			case expr.COS:
				M := expr.NewMul()
				M.Insert(n.(*expr.Cos).C)
				M.Insert(B.Clone())
				n.(*expr.Cos).C = M
			case expr.TAN:
				M := expr.NewMul()
				M.Insert(n.(*expr.Tan).C)
				M.Insert(B.Clone())
				n.(*expr.Tan).C = M
			case expr.EXP:
				M := expr.NewMul()
				M.Insert(n.(*expr.Exp).C)
				M.Insert(B.Clone())
				n.(*expr.Exp).C = M
			case expr.LOG:
				M := expr.NewMul()
				M.Insert(n.(*expr.Log).C)
				M.Insert(B.Clone())
				n.(*expr.Log).C = M
			}
			sort.Sort(a)
			a.CalcExprStats()
			good := PS.cnfg.treecfg.CheckExpr(a)
			if good {
				ret = append(ret, a)
			}
			// fmt.Printf("grew  %v\n", a)
			ret = append(ret, a)
		}
	}
	// fmt.Println("Len of ret = ", len(ret))
	return ret
}
