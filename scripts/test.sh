#!/bin/bash

files=(
	# Koza_1
	# Koza_2
	# Koza_3
	Nguyen_01
	# Nguyen_02
	# Nguyen_03
	# Nguyen_04
	# Nguyen_05
	# Nguyen_06
	# Nguyen_07
	# Nguyen_08
	# Nguyen_09
	# Nguyen_10
	# Nguyen_11
	# Nguyen_12


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
	# Korns_14
	# Korns_15
)

methods=(
	method1
	# method2
	# method3
)

peels=(
	# 1
	# 2
	3
	# 4
)

iters=(
	# 25
	200
	# 400
	# 1000
)



mkdir -p runs 

echo "ULTIMATE START"
echo "--------------"
date
echo
echo


for F in ${files[@]}
do
	mkdir -p runs/$F/pge
	echo " " > runs/${F}_pge_tails.txt 
	echo " " > runs/${F}_pge_fit.txt 

	for M in ${methods[@]}; do
	for I in ${iters[@]}; do
	for P in ${peels[@]}; do

	echo "PGE $I $P $M :: $F  "
	echo "------------"
	date
	echo "./go-pge -pcfg=prob/bench/${F}.cfg -peel=${P} -iter=${I} -init=${M} -grow=${M} > 'runs/${F}/${F}_pge_${I}_${P}_${M}.out'"
	time ./go-pge -pcfg=prob/bench/${F}.cfg -peel=${P} -iter=${I} -init=${M} -grow=${M} > "runs/${F}/${F}_pge_${I}_${P}_${M}.out"
	# gdb ./go-pge

	cp runs/${F}/pge/pge/pge:fitness.log runs/${F}/${F}_pge_${I}_${P}_${M}.fit
	echo "${F}_pge_${I}_${P}_${M}" >> runs/${F}_pge_fit.txt 
	tail -n 7 runs/${F}/pge/pge/pge:fitness.log >> runs/${F}_pge_fit.txt
	for i in {1..4}; do
		echo "" >> runs/${F}_pge_fit.txt 
	done
	# 
	echo "${F}_pge_${I}_${P}_${M}.out" >> runs/${F}_pge_tails.txt 
	tail -n 333 runs/${F}/${F}_pge_${I}_${P}_${M}.out >> runs/${F}_pge_tails.txt 
	for i in {1..200}; do
		echo "" >> runs/${F}_pge_tails.txt 
	done
	
	echo
	echo
	echo

	done
		for i in {1..5}; do
			echo "" >> runs/${F}_pge_fit.txt 
		done
	done
		for i in {1..10}; do
			echo "" >> runs/${F}_pge_fit.txt 
		done
	done



done

echo
date
echo
echo
echo