go-pge
======

If you publish any work using this software, please cite:

Tony Worm and Kenneth Chiu. Prioritized Grammar Enumeration: Symbolic Regression by Dynamic Programming.
In Proceedings of the Genetic and Evolutionary Computation Conference (GECCO’2013). July 6-10 Amsterdam, Netherlands.

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


Installation
=====================================

1. install Go  (golang.org)
2. add $GOROOT to $PATH
3. sudo apt-get install g++ cmake liblapack3 liblapack-dev libblas3 libblas-dev f2c (ubuntu)

4. go get github.com/verdverm/go-pge
5. navigate to github.com/verdverm/go-levmar/levmar-2.6
  * cmake -DCMAKE_BUILD_TYPE=RelWithDebInfo -DLINSOLVERS_RETAIN_MEMORY=0 .
  * make
6. navigate to github.com/verdverm/go-pge
7. go build

now you should be able to run scripts/test.sh
