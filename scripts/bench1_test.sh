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
	# Korns_05 ? no good equations ?
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

# cp scripts/post/eqnMakefile out/bench
# cat scripts/post/eqns_top.tex > out/bench/eqns.tex

echo "ULTIMATE START"
echo "--------------"
date
echo
echo
for F in ${files[@]}
do
	echo "GPSR: $F"
	mkdir -p runs/$F/{gpsr,pesr}
	mkdir -p out/bench/$F
	echo "\include{$F}" >> out/bench/eqns.tex

	for I in {1..50}
	do
		str="$F"
		if [[ "$I" -lt "10" ]]; then
			dir="runs/$F/gpsr/run_0$I"
			mkdir -p $dir 
			str="$dir/gpsr.txt" 
		else
			dir="runs/$F/gpsr/run_$I"
			mkdir -p $dir 
			str="$dir/gpsr.txt" 
		fi

		echo $str
		date
		time ./damd -pcfg=prob/bench/$F.cfg > "$str"
		echo 
		echo
	done
	
	echo "PESR: $F"
	date
	echo "runs/$F/pesr/pesr.txt"
	time ./damd -pcfg=prob/bench/$F.cfg -scfg=pesr/pesr_default.cfg > "runs/$F/pesr/pesr.txt"

	echo
	echo
done

# cat scripts/post/eqns_bot.tex >> out/bench/eqns.tex

# run post processing


# -pcfg=prob/bench/Korns_01.cfg -scfg=pesr/pesr_default.cfg
