go-pge
======

Implementation of Prioritized Grammar Enumeration


go-pge -help for options

check out the scripts/test.sh for an example of how to run go-pge


If you publish any work and used this software, please cite:

Tony Worm and Kenneth Chiu. Prioritized Grammar Enumeration: Symbolic Regression by Dynamic Programming.
In Proceedings of the Genetic and Evolutionary Computation Conference (GECCO’2013). July 6-10 Amsterdam, Netherlands.



Sorry there is a bunch of garbage hanging around.  
(and some dependencies go-symexpr & go-levmar [which is harder to get setup])


Some insight into how to use go-pge
===================================

scripts/test.sh has an example of usage, plus batch running

It basically uses a bunch of config files as parameters

some places to look for example files:
  - the config/* directories
  - main.go  (that's where all of the cmd line args are)


go-pge processes data sets

I set it up to read a particular data format, maybe more than one, can't recall for sure
There may be code paths hanging around to allow you to use multiple input sets



