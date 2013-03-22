#!/bin/bash

files=(
	Koza_1
	Koza_2
	Koza_3
	Nguyen_01
	Nguyen_02
	Nguyen_03
	Nguyen_04
	Nguyen_05
	Nguyen_06
	Nguyen_07
	Nguyen_08
	Nguyen_09
	Nguyen_10
	Nguyen_11
	Nguyen_12

	# Korns_01
	# Korns_02
	# Korns_03
	# Korns_04
	# Korns_05
	# Korns_06
	# Korns_07
	# Korns_08
	# Korns_09
	# Korns_10
	# Korns_11
	# Korns_12
	# Korns_13
	# Korns_15

	
)

mkdir -p runs out/bench

cp scripts/post/eqnMakefile out/bench
cat scripts/post/eqns_top.tex > out/bench/eqns.tex

for F in ${files[@]}
do
	mkdir -p runs/$F/{gpsr,pesr}
	mkdir -p out/bench/$F
	echo "\include{$F}" >> out/bench/eqns.tex

	for I in {0..29}
	do
		str="$F"
		if [[ "$I" -lt "10" ]]; then
			str="runs/$F/gpsr/gpsr_0$I" 
		else
			str="runs/$F/gpsr/gpsr_$I" 
		fi
		echo $str
		time ./damd -pcfg=prob/bench/$F.cfg > "$str.txt"
	
	done
	# time ./damd -pcfg=prob/bench/$F.cfg > "runs/$F/pesr.txt"
	echo
done

cat scripts/post/eqns_bot.tex >> out/bench/eqns.tex

# run post processing
