package pge

import (
	// "fmt"
	"sort"

	expr "github.com/verdverm/go-symexpr"
)

func (PS *PgeSearch) ExpandMethod1(O expr.Expr) (ret []expr.Expr) {
	O.Sort()
	ret = make([]expr.Expr, 0)
	// fmt.Printf("Expanding expression:  %v\n", O)

	for i := 0; i < O.Size(); i++ {
		I := i
		E := O.GetExpr(&I)
		switch E.ExprType() {
		case expr.ADD:
			tmp := PS.AddTermToExprMethod1(O, E, i)
			ret = append(ret, tmp[:]...)

		case expr.MUL:
			tmp := PS.WidenTermInExprMethod1(O, E, i)
			ret = append(ret, tmp[:]...)

		case expr.VAR:
			tmp := PS.DeepenTermInExprMethod1(O, E, i)
			ret = append(ret, tmp[:]...)

		default: // expr.DIV,expr.COS,expr.SIN,expr.EXP,expr.LOG,expr.ABS,expr.POW
			continue

		}
	}

	return ret
}

// add another term to an add expr
func (PS *PgeSearch) AddTermToExprMethod1(O, E expr.Expr, pos int) (ret []expr.Expr) {
	ret = make([]expr.Expr, 0)
	A := E.(*expr.Add)

	// f() + cL
	for _, L := range PS.GenLeafs {
		c := new(expr.Constant)
		c.P = -1
		l := L.Clone()

		// mul it
		M := expr.NewMul()
		M.Insert(c)
		M.Insert(l)

		// skip if the same term already exists in the add
		skip := false
		for _, e := range A.CS {
			if e == nil {
				continue
			}
			// fmt.Printf("ACMP  %v  %v\n", M, e)
			if e.AmIAlmostSame(M) || M.AmIAlmostSame(e) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		C := O.Clone()
		P := pos
		AM := C.GetExpr(&P).(*expr.Add)
		AM.Insert(M)
		sort.Sort(AM)
		C.CalcExprStats()
		good := PS.cnfg.treecfg.CheckExpr(C)
		if good {
			ret = append(ret, C)
		}
	}

	// f() + c*node(L)
	for _, N := range PS.GenNodes {
		for _, L := range PS.GenLeafs {

			c := new(expr.Constant)
			c.P = -1

			l := L.Clone()
			n := N.Clone()
			p := 1
			n.SetExpr(&p, l)

			var E expr.Expr
			if N.ExprType() == expr.DIV {
				E = expr.NewDiv(c, l)
			} else {
				// mul it
				M := expr.NewMul()
				M.Insert(c)
				M.Insert(n)
				E = M
			}

			// skip if the same term already exists in the add
			skip := false
			for _, e := range A.CS {
				if e == nil {
					continue
				}
				// fmt.Printf("ACMP  %v  %v\n", M, e)
				if e.AmIAlmostSame(E) || E.AmIAlmostSame(e) {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			// fmt.Println(E.String())

			C := O.Clone()
			P := pos
			AM := C.GetExpr(&P).(*expr.Add)
			AM.Insert(E)
			sort.Sort(AM)
			C.CalcExprStats()
			good := PS.cnfg.treecfg.CheckExpr(C)
			if good {
				ret = append(ret, C)
			}
		}
	}

	return ret
}

// add complexity to a single multiplication term
func (PS *PgeSearch) WidenTermInExprMethod1(O, E expr.Expr, pos int) (ret []expr.Expr) {
	ret = make([]expr.Expr, 0)

	// insert leafs  f()*L
	for _, L := range PS.GenLeafs {
		l := L.Clone()
		C := O.Clone()
		P := pos
		e := C.GetExpr(&P)
		// fmt.Printf("pos(%d): %v\n", pos, e)
		M := e.(*expr.Mul)
		M.Insert(l)
		sort.Sort(M)
		C.CalcExprStats()
		good := PS.cnfg.treecfg.CheckExpr(C)
		if good {
			ret = append(ret, C)
		}
	}

	// insert node(L)  :  f() * c*node(L)
	for _, N := range PS.GenNodes {
		for _, L := range PS.GenLeafs {
			c := new(expr.Constant)
			c.P = -1

			l := L.Clone()
			n := N.Clone()
			p := 1
			n.SetExpr(&p, l)

			var E expr.Expr
			if N.ExprType() == expr.DIV {
				E = expr.NewDiv(c, l)
			} else {
				// mul it
				M := expr.NewMul()
				M.Insert(c)
				M.Insert(n)
				E = M
			}

			C := O.Clone()
			P := pos
			e := C.GetExpr(&P)
			// fmt.Printf("pos(%d): %v\n", pos, e)
			M := e.(*expr.Mul)

			M.Insert(E)
			sort.Sort(M)
			C.CalcExprStats()
			good := PS.cnfg.treecfg.CheckExpr(C)
			if good {
				ret = append(ret, C)
			}
		}
	}
	return ret
}

// change any term to something more complex...
func (PS *PgeSearch) DeepenTermInExprMethod1(O, E expr.Expr, pos int) []expr.Expr {
	exprs := make([]expr.Expr, 0)

	// make into add
	A := expr.NewAdd()
	A.Insert(E.Clone())
	OA := A.Clone()
	exprs = append(exprs, PS.AddTermToExprMethod1(OA, A, 0)[:]...)

	// make into mul
	M := expr.NewMul()
	M.Insert(E.Clone())
	OM := M.Clone()
	exprs = append(exprs, PS.WidenTermInExprMethod1(OM, M, 0)[:]...)

	// // make into div
	// if E.ExprType() != expr.DIV {
	// 	D := new(expr.Div)
	// 	c := new(expr.Constant)
	// 	c.P = -1
	// 	D.Numer = c
	// 	D.Denom = E.Clone()
	// 	exprs = append(exprs, D)
	// }

	// make inside of nodes
	for _, N := range PS.GenNodes {
		if N.ExprType() == expr.DIV {
			continue
		}
		T := N.Clone()
		P := 1
		T.SetExpr(&P, E.Clone())
		exprs = append(exprs, T)
	}

	ret := make([]expr.Expr, 0)
	for _, e := range exprs {
		C := O.Clone()
		P := pos
		C.SetExpr(&P, e)
		ret = append(ret, C)
	}
	return ret
}
