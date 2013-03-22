package problems

import (
	// "fmt"
	expr "github.com/verdverm/go-symexpr"
	"math"
	"math/rand"
)

type RangeType int

const (
	_ RangeType = iota
	Uniform
	Equal
)

type BenchmarkVar struct {
	Name    string
	Index   int
	Rtype   RangeType
	L, H, S float64 // low,high,step of range
}

type Benchmark struct {
	Name         string
	TrainVars    []BenchmarkVar
	TrainSamples int
	TestVars     []BenchmarkVar
	TestSamples  int

	Functions []string
	NonTrig   []string
	FuncText  string // function as text

}

func GenBenchmark(b Benchmark) (p *ExprProblem) {
	p = new(ExprProblem)
	p.Name = b.Name

	// set the function
	varNames := make([]string, 0)
	for _, v := range b.TrainVars {
		// fmt.Printf("  %v\n", v)
		varNames = append(varNames, v.Name)
	}
	eqn := expr.ParseFunc(b.FuncText, varNames)
	sort := eqn.Clone()
	rules := expr.DefaultRules()
	rules.GroupAddTerms = false
	sort = sort.Simplify(rules)
	p.VarNames = varNames
	p.FuncTree = sort

	trn := new(PointSet)
	trn.SetFN(b.Name + "_train")
	trn.SetID(0)
	trn.SetNumDim(len(varNames))
	trn.SetIndepNames(varNames)
	trn.SetDepndNames([]string{"f(xs)"})
	trn.SetPoints(GenBenchData(eqn, b.TrainVars, b.TrainSamples))
	p.Train = make([]*PointSet, 1)
	p.Train[0] = trn

	tst := new(PointSet)
	tst.SetFN(b.Name + "_test")
	tst.SetID(0)
	tst.SetNumDim(len(varNames))
	tst.SetIndepNames(varNames)
	tst.SetDepndNames([]string{"f(xs)"})
	tst.SetPoints(GenBenchData(eqn, b.TestVars, b.TestSamples))
	p.Test = make([]*PointSet, 1)
	p.Test[0] = tst

	return p
}

func GenBenchData(e expr.Expr, vars []BenchmarkVar, samples int) (pts []Point) {
	pts = make([]Point, 0)
	if vars[0].Rtype == Uniform {
		for i := 0; i < samples; i++ {
			input := make([]float64, len(vars))
		retry:
			for j, v := range vars {
				r := rand.Float64()
				input[j] = (r * (v.H - v.L)) + v.L
			}
			out := e.Eval(0, input, nil, nil)

			if math.IsNaN(out) || math.IsInf(out, 0) || math.Abs(out) > 100000.0 {
				goto retry
			}

			var pnt Point
			pnt.SetIndeps(input)
			pnt.SetDepnds([]float64{out})
			pts = append(pts, pnt)
		}
	} else { // RangeType == Equal
		counter := make([]float64, len(vars))
		for j, v := range vars {
			counter[j] = v.L
		}
		L1, L2 := len(vars)-1, vars[len(vars)-1].L

		for counter[L1] <= L2 {
			input := make([]float64, len(vars))
			copy(input, counter)

			out := e.Eval(0, input, nil, nil)
			var pnt Point
			pnt.SetIndeps(input)
			pnt.SetDepnds([]float64{out})
			pts = append(pts, pnt)

			// increment counter
			for j, v := range vars {
				counter[j] += v.S
				if counter[j] > v.H {
					counter[j] = v.L
				} else {
					break
				}
			}
		}

	}

	return
}

var xU11 = BenchmarkVar{"x", 0, Uniform, -1.0, 1.0, 0.0}
var xU22 = BenchmarkVar{"x", 0, Uniform, -2.0, 2.0, 0.0}
var xU33 = BenchmarkVar{"x", 0, Uniform, -3.0, 3.0, 0.0}
var xU01 = BenchmarkVar{"x", 0, Uniform, 0.0, 1.0, 0.0}
var xU02 = BenchmarkVar{"x", 0, Uniform, 0.0, 2.0, 0.0}
var xU04 = BenchmarkVar{"x", 0, Uniform, 0.0, 4.0, 0.0}
var yU01 = BenchmarkVar{"y", 1, Uniform, 0.0, 1.0, 0.0}

var xU5050 = BenchmarkVar{"x", 0, Uniform, -50.0, 50.0, 0.0}
var yU5050 = BenchmarkVar{"y", 1, Uniform, -50.0, 50.0, 0.0}
var zU5050 = BenchmarkVar{"z", 2, Uniform, -50.0, 50.0, 0.0}
var vU5050 = BenchmarkVar{"v", 3, Uniform, -50.0, 50.0, 0.0}
var wU5050 = BenchmarkVar{"w", 4, Uniform, -50.0, 50.0, 0.0}

var korns5 = []BenchmarkVar{xU5050, yU5050, zU5050, vU5050, wU5050}
var xyU01 = []BenchmarkVar{xU01, yU01}

var xE11_01 = BenchmarkVar{"x", 0, Equal, -1.0, 1.0, 0.01}
var xE22_01 = BenchmarkVar{"x", 0, Equal, -2.0, 2.0, 0.01}
var xE33_01 = BenchmarkVar{"x", 0, Equal, -3.0, 3.0, 0.01}
var xE11_001 = BenchmarkVar{"x", 0, Equal, -1.0, 1.0, 0.001}
var xE22_001 = BenchmarkVar{"x", 0, Equal, -2.0, 2.0, 0.001}
var xE33_001 = BenchmarkVar{"x", 0, Equal, -3.0, 3.0, 0.001}

var xE55_4 = BenchmarkVar{"x", 0, Equal, -5.0, 5.0, 0.4}
var yE55_4 = BenchmarkVar{"y", 1, Equal, -5.0, 5.0, 0.4}
var xyE55_4 = []BenchmarkVar{xE55_4, yE55_4}

var kozaFuncs = []string{"Add", "Mul", "Div", "Sin", "Cos", "Exp", "Log"}
var kornsFuncs = []string{"Add", "Mul", "Div", "Sin", "Cos", "Exp", "Log", "Sqrt"}
var vladiFuncsA = []string{"Add", "Mul", "Div"}
var vladiFuncsB = []string{"Add", "Mul", "Div", "Exp"}
var vladiFuncsC = []string{"Add", "Mul", "Div", "Sin", "Cos", "Exp"}

var kozaFuncs_NT = []string{"Add", "Mul", "Div", "Exp", "Log"}
var kornsFuncs_NT = []string{"Add", "Mul", "Div", "Exp", "Log", "Sqrt"}
var vladiFuncsA_NT = []string{"Add", "Mul", "Div"}
var vladiFuncsB_NT = []string{"Add", "Mul", "Div", "Exp"}
var vladiFuncsC_NT = []string{"Add", "Mul", "Div", "Exp"}

var BenchmarkList = []Benchmark{
	Benchmark{"Koza_1", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^4 + x^3 + x^2 + x"},
	Benchmark{"Koza_2", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^5 - 2x^3 + x"},
	Benchmark{"Koza_3", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^6 - 2x^4 + x^2"},

	Benchmark{"Nguyen_01", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^3 + x^2 + x"},
	Benchmark{"Nguyen_02", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^4 + x^3 + x^2 + x"},
	Benchmark{"Nguyen_03", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^5 + x^4 + x^3 + x^2 + x"},
	Benchmark{"Nguyen_04", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, vladiFuncsA, vladiFuncsA_NT, "x^6 + x^5 + x^4 + x^3 + x^2 + x"},
	Benchmark{"Nguyen_05", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, kozaFuncs, kozaFuncs_NT, "sin(x^2)*cos(x) - 1"},
	Benchmark{"Nguyen_06", []BenchmarkVar{xU11}, 200, []BenchmarkVar{xU11}, 2000, kozaFuncs, kozaFuncs_NT, "sin(x) + sin(x + x^2)"},
	Benchmark{"Nguyen_07", []BenchmarkVar{xU02}, 200, []BenchmarkVar{xU02}, 2000, kozaFuncs, kozaFuncs_NT, "ln(x+1) + ln(x^2 + 1)"},
	Benchmark{"Nguyen_08", []BenchmarkVar{xU04}, 200, []BenchmarkVar{xU04}, 2000, kozaFuncs, kozaFuncs_NT, "sqrt(x)"},
	Benchmark{"Nguyen_09", xyU01, 200, xyU01, 2000, kozaFuncs, kozaFuncs_NT, "sin(x) + sin(y^2)"},
	Benchmark{"Nguyen_10", xyU01, 200, xyU01, 2000, kozaFuncs, kozaFuncs_NT, "2*sin(x)*cos(y)"},
	Benchmark{"Nguyen_11", xyU01, 200, xyU01, 2000, kozaFuncs, kozaFuncs_NT, "x^y"},
	Benchmark{"Nguyen_12", xyU01, 200, xyU01, 2000, vladiFuncsA, vladiFuncsA_NT, "x^4 - x^3 + 0.5*y^2 - y"},

	Benchmark{"Pagie_1", xyE55_4, 0, xyE55_4, 0, kozaFuncs, kozaFuncs_NT, "1 / (1 + x^-4) + 1 / (1 + y^-4)"},

	// 5 inputs: x,y,z,v,w
	Benchmark{"Korns_01", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "1.57 + 24.3*v"},
	Benchmark{"Korns_02", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "0.23 + 14.2*(v+y)/3w"},
	Benchmark{"Korns_03", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "-5.41 + 4.9*(v-x+y/w)/3w"},
	Benchmark{"Korns_04", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "-2.3 + 0.13sin(z)"},
	Benchmark{"Korns_05", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "3 + 2.13*ln(w)"},
	Benchmark{"Korns_06", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "1.3 + 0.13*sqrt(x)"},
	Benchmark{"Korns_07", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "213.80940889*(1 - e^(-0.54723748542*x))"},
	Benchmark{"Korns_08", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "6.87 + 11*sqrt(7.23*x*v*w)"},
	Benchmark{"Korns_09", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "(sqrt(x)/ln(y)) * (e^z / v^2)"},
	Benchmark{"Korns_10", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "0.81 + 24.3*(2y+3*z^2)/(4*(v)^3+5*(w)^4)"},
	Benchmark{"Korns_11", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "6.87 + 11*cos(7.23*x^3)"},
	Benchmark{"Korns_12", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "2 - 2.1*cos(9.8*x)*sin(1.3*w)"},
	Benchmark{"Korns_13", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "32 - 3*(tan(x)*tan(z))/(tan(y)*tan(v))"},
	Benchmark{"Korns_14", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "22 - 4.2*(cos(x)-tan(y))*(tanh(z)/sin(v))"},
	Benchmark{"Korns_15", korns5, 200, korns5, 2000, kornsFuncs, kornsFuncs_NT, "12 - 6*(tan(x)/e^y)(ln(z)-tan(v))"},

	Benchmark{"Keijzer_01", []BenchmarkVar{xE11_01}, 0, []BenchmarkVar{xE11_001}, 0, vladiFuncsA, vladiFuncsA_NT, "0.3*x*sin(2*PI*x)"},
	Benchmark{"Keijzer_02", []BenchmarkVar{xE22_01}, 0, []BenchmarkVar{xE22_001}, 0, vladiFuncsA, vladiFuncsA_NT, "0.3*x*sin(2*PI*x)"},
	Benchmark{"Keijzer_03", []BenchmarkVar{xE33_01}, 0, []BenchmarkVar{xE33_001}, 0, vladiFuncsA, vladiFuncsA_NT, "0.3*x*sin(2*PI*x)"},
	// Benchmark{"Keijzer_04", "x", "E[0,10,0.05]", "x^3*e^-x*cos(x)*sin(x)*((sin(x))^2*cos(x) - 1)"},
	// Benchmark{"Keijzer_05", "x,y,z", "x,z: U[-1,1,1000] y: U[1,2,1000]", "(30*x*z) / ((x-10)*y^2)"},
	// // Benchmark{"Keijzer_06", "x", "E[1,50,1]", "\\SUM_i^x (1/i)"},
	// Benchmark{"Keijzer_07", "x", "E[1,100,1]", "ln(x)"},
	// Benchmark{"Keijzer_08", "x", "E[0,100,1]", "sqrt(x)"},
	// // arcsinh(x) == ln(x+sqrt(x^2+1))
	// Benchmark{"Keijzer_09", "x", "E[0,100,1]", "ln(x+sqrt(x^2+1))"},
	// Benchmark{"Keijzer_10", "x,y", "U[0,1,100]", "x^y"},
	// // xy ? sin((x-1)(y-1))
	// // Benchmark{"Keijzer_11", "x,y", "U[-3,3,20]", "xy ? sin((x-1)*(y-1))"},
	// Benchmark{"Keijzer_12", "x,y", "U[-3,3,20]", "x^4 - x^3 + 0.5*y^2 - y "},
	// Benchmark{"Keijzer_13", "x,y", "U[-3,3,20]", "6*sin(x)*cos(y)"},
	// Benchmark{"Keijzer_14", "x,y", "U[-3,3,20]", "8 / (2 + x^2 + y^2) "},
	// Benchmark{"Keijzer_15", "x,y", "U[-3,3,20]", "0.2*x^3 + 0.5*y^2 - y - x "},

	// // (e^{-(x-1)^2}) / (1.2 + (y-2.5)^2)
	// Benchmark{"Vladislavleva_1", "x,y", "?U[0.3,4,10]?", "(e^{-(x-1)^2}) / (1.2 + (y-2.5)^2)"},
	// Benchmark{"Vladislavleva_2", "x", "E[0.05,10,0.1]", "e^-x* x^3 * (cos(x)*sin(x)) * (cos(x)*(sin(x))^2-1)"},
	// Benchmark{"Vladislavleva_3", "x,y", "x: E[0.05,10,0.1]  y: E[0.05,10.05,2]", "e^-y * x^3 * (cos(x)*sin(x)) * (cos(x)*(sin(x))^2-1)"},
	// // Benchmark{"Vladislavleva_4", "x_i", " U[0.05, 6.05, 1024]", "10 / (5 + \\SUM_1^5 (x_i - 3)^2)"},
	// Benchmark{"Vladislavleva_5", "x,y,z", "x,z: U[0.05,2,300]  y: U[1,2,300]", "(30*(x-1)*(z-1)) / (y^2*(x-10))  "},
	// Benchmark{"Vladislavleva_6", "x,y", "U[0.1,5.9,30]", "6*sin(x)*cos(y)"},
	// Benchmark{"Vladislavleva_7", "x,y", "U[0.05,6.05,300]", "(x-3)*(y-3) + 2sin((x-4)*(y-4))"},
	// Benchmark{"Vladislavleva_8", "x,y", "U[0.05,6.05,50]", "((x-3)^4 + (y-3)^3 - (y-3)) / ((y-2)^4 + 10)"},
}
