go-pge
======

If you publish any work using this software, please cite:

Tony Worm and Kenneth Chiu. Prioritized Grammar Enumeration: Symbolic Regression by Dynamic Programming.
In Proceedings of the Genetic and Evolutionary Computation Conference (GECCO’2013). July 6-10 Amsterdam, Netherlands.

Also, check out www.symbolicregression.org   [I'm new to web-dev'n ;]

Sorry there is a bunch of garbage hanging around in the code base.  
(and a dependency go-levmar [which can be a pain to get setup])


Some insight into how to use go-pge
===================================

scripts/test.sh has an example of usage, plus batch running

It basically uses a bunch of config files as parameters

some places to look for example files:
  - the config/* directories
  - main.go  (that's where all of the cmd line args are)


go-pge processes data sets

I set it up to read a particular data format, maybe more than one, can't recall for sure
basically white space seperated with the variable names at the top of the file
(sample benchmark files are around)

There may be code paths hanging around to allow you to use multiple input sets
