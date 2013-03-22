package main

import (
	"fmt"
	"strings"

	expr "github.com/verdverm/go-symexpr"
	probs "go-pge/problems"
)

func genBenchmark(probname string) {
	fmt.Printf("Generating files for %s\n", probname)
	bname := probname[strings.Index(probname, ":")+1:]
	fmt.Printf("bname:  '%s'\n", bname)
	benches := probs.BenchmarkList
	switch bname {
	case "all":
		for i := 0; i < len(benches); i++ {
			fmt.Printf("bm: %v\n", benches[i])
			bname = benches[i].Name
			symprob := probs.GenBenchmark(benches[i])
			symprob.Train[0].WritePointSet("data/benchmark/" + bname + ".trn")
			symprob.Test[0].WritePointSet("data/benchmark/" + bname + ".tst")
			fmt.Println()
		}
	case "list":
		printBenchNames()
	default:
		i := 0
		for ; i < len(benches); i++ {
			if benches[i].Name == bname {
				fmt.Printf("bm: %v\n", benches[i])
				symprob := probs.GenBenchmark(benches[i])
				symprob.Train[0].WritePointSet("data/benchmark/" + bname + ".trn")
				symprob.Test[0].WritePointSet("data/benchmark/" + bname + ".tst")
				fmt.Println()
				return
			}
		}
		if i == len(benches) {
			fmt.Printf("benchmark problem not found:  %s\n", probname)
			printBenchNames()
		}
	}

}

func printBenchNames() {
	benches := probs.BenchmarkList
	fmt.Printf("Available Benchmarks:  |%d|\n[ %s", len(benches), benches[0].Name)
	for i := 0; i < len(benches); i++ {
		fmt.Printf(", %s", benches[i].Name)
	}
	fmt.Println(" ]")
}

func printBenchLatex() {

	for _, B := range probs.BenchmarkList {

		varNames := make([]string, 0)
		for _, v := range B.TrainVars {
			// fmt.Printf("  %v\n", v)
			varNames = append(varNames, v.Name)
		}
		eqn := expr.ParseFunc(B.FuncText, varNames)
		sort := eqn.Clone()
		rules := expr.DefaultRules()
		rules.GroupAddTerms = false
		sort = sort.Simplify(rules)

		latex := sort.Latex(varNames, nil, nil)

		fmt.Println(B.Name, " \t&\t $", latex, "$ \\\\")
	}
}
