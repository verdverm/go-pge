package problems

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

type Point struct {
	indep []float64
	depnd []float64
}

func (d *Point) NumIndep() int             { return len(d.indep) }
func (d *Point) SetNumIndep(sz int)        { d.indep = make([]float64, sz) }
func (d *Point) Indep(p int) float64       { return d.indep[p] }
func (d *Point) SetIndep(p int, v float64) { d.indep[p] = v }
func (d *Point) Indeps() []float64         { return d.indep }
func (d *Point) SetIndeps(v []float64)     { d.indep = v }
func (d *Point) NumDepnd() int             { return len(d.depnd) }
func (d *Point) SetNumDepnd(sz int)        { d.depnd = make([]float64, sz) }
func (d *Point) Depnd(p int) float64       { return d.depnd[p] }
func (d *Point) SetDepnd(p int, v float64) { d.depnd[p] = v }
func (d *Point) Depnds() []float64         { return d.depnd }
func (d *Point) SetDepnds(v []float64)     { d.depnd = v }

type PointSet struct {
	filename string
	id       int

	numDim     int
	indepNames []string
	depndNames []string
	sysNames   []string

	dataPoints []Point
	sysVals    []float64
}

func (d *PointSet) FN() string           { return d.filename }
func (d *PointSet) SetFN(fn string)      { d.filename = fn }
func (d *PointSet) ID() int              { return d.id }
func (d *PointSet) SetID(id int)         { d.id = id }
func (d *PointSet) SetNumPoints(cnt int) { d.dataPoints = make([]Point, cnt) }

func (d *PointSet) NumIndep() int     { return len(d.indepNames) }
func (d *PointSet) NumDepnd() int     { return len(d.depndNames) }
func (d *PointSet) NumDim() int       { return d.numDim } // TODO check to see if TIME is a variable
func (d *PointSet) SetNumDim(dim int) { d.numDim = dim }  // TODO check to see if TIME is a variable
func (d *PointSet) NumSys() int       { return len(d.sysNames) }
func (d *PointSet) NumPoints() int    { return len(d.dataPoints) }

func (d *PointSet) IndepName(xi int) string      { return d.indepNames[xi] }
func (d *PointSet) GetIndepNames() []string      { return d.indepNames }
func (d *PointSet) SetIndepNames(names []string) { d.indepNames = names }
func (d *PointSet) DepndName(xi int) string      { return d.depndNames[xi] }
func (d *PointSet) GetDepndNames() []string      { return d.depndNames }
func (d *PointSet) SetDepndNames(names []string) { d.depndNames = names }

func (d *PointSet) SysName(si int) string      { return d.sysNames[si] }
func (d *PointSet) GetSysNames() []string      { return d.sysNames }
func (d *PointSet) SetSysNames(names []string) { d.sysNames = names }
func (d *PointSet) SetSysVals(sv []float64)    { d.sysVals = sv }

func (d *PointSet) Point(p int) *Point    { return &(d.dataPoints[p]) }
func (d *PointSet) Points() []Point       { return d.dataPoints }
func (d *PointSet) SetPoints(pts []Point) { d.dataPoints = pts }
func (d *PointSet) SysVal(p int) float64  { return d.sysVals[p] }
func (d *PointSet) SysVals() []float64    { return d.sysVals }

// read function at end of file  [ func (d *PointSet) Read(filename string) ]

type PntSubset struct {
	ds *PointSet

	index   []int
	input   []Point
	output  []Point
	sysVals []float64
}

func (s *PntSubset) ID() int            { return s.ds.id }
func (s *PntSubset) DS() *PointSet      { return s.ds }
func (s *PntSubset) SetDS(ds *PointSet) { s.ds = ds }

func (s *PntSubset) NumIndep() int  { return s.ds.NumIndep() }
func (s *PntSubset) NumDepnd() int  { return s.ds.NumDepnd() }
func (s *PntSubset) NumSys() int    { return s.ds.NumSys() }
func (s *PntSubset) NumPoints() int { return len(s.index) }

func (s *PntSubset) SysVals() []float64        { return s.sysVals }
func (s *PntSubset) SetSysVals(svls []float64) { s.sysVals = svls }

func (s *PntSubset) Index(p int) int       { return s.index[p] }
func (s *PntSubset) Indexes() []int        { return s.index }
func (s *PntSubset) SetIndexes(idxs []int) { s.index = idxs }
func (s *PntSubset) Input(p int) *Point    { return &s.input[p] }
func (s *PntSubset) Output(p int) *Point   { return &s.output[p] }

func (s *PntSubset) AddPoint(p int, input, output *Point) {
	s.index = append(s.index, p)
	s.input = append(s.input, *input)
	s.output = append(s.output, *output)
}

// using indexes, update the input/output data
func (s *PntSubset) Refresh() {
	L := len(s.index)
	if s.input == nil {
		s.input = make([]Point, L)
	}
	if s.output == nil {
		s.output = make([]Point, L)
	}

	for i := 0; i < L; i++ {
		s.input[i] = *s.ds.Point(s.index[i])
		if s.index[i]+1 >= s.ds.NumPoints() {
			continue
		}
		s.output[i] = *s.ds.Point(s.index[i] + 1)
	}
}

func (d *PointSet) ReadPointSet(filename string) {
	ftotal, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer ftotal.Close()
	file := bufio.NewReader(ftotal)

	var word string

	// get independent variables (x_i...)
	for i := 0; ; i++ {
		_, err := fmt.Fscanf(file, "%s", &word)
		if err != nil {
			break
		}
		d.indepNames = append(d.indepNames, word)
	}
	d.numDim = len(d.indepNames)

	// get dependent variables (y_j...)
	for i := 0; ; i++ {
		_, err := fmt.Fscanf(file, "%s", &word)
		if err != nil {
			break
		}
		d.depndNames = append(d.depndNames, word)
	}

	fmt.Printf("Var Names = %v | %v\n", d.depndNames, d.indepNames)

	for i := 0; ; i++ {
		var pnt Point
		var dval, ival float64
		if err != nil {
			break
		}

		for j := 0; j < len(d.indepNames); j++ {
			_, err = fmt.Fscanf(file, "%f", &ival)
			if err != nil {
				break
			}
			pnt.indep = append(pnt.indep, ival)
		}

		for j := 0; j < len(d.depndNames); j++ {
			_, err = fmt.Fscanf(file, "%f\n", &dval)
			if err != nil {
				break
			}
			pnt.depnd = append(pnt.depnd, dval)
		}

		if len(pnt.indep) > 0 {
			d.dataPoints = append(d.dataPoints, pnt)
		}
		if i%100 == 0 {
			fmt.Println("Point(%d): %v\n", i, pnt)
		}
	}
	fmt.Printf("Num Points: %v\n", len(d.dataPoints))
}

func (d *PointSet) WritePointSet(filename string) {
	fmt.Printf("Writing file: %s\n", filename)
	ftotal, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %s  %v\n", filename, err)
		return
	}
	defer ftotal.Close()
	file := bufio.NewWriter(ftotal)
	defer file.Flush()

	// write independent variable names (x_i...)
	for i := 0; i < d.NumIndep(); i++ {
		_, err := fmt.Fprintf(file, "%s ", d.IndepName(i))
		if err != nil {
			fmt.Errorf("error writing pointset to file: %v\n", err)
			break
		}
	}
	fmt.Fprintln(file)

	// trite dependent variable names (y_j...)
	for i := 0; i < d.NumDepnd(); i++ {
		_, err := fmt.Fprintf(file, "%s ", d.DepndName(i))
		if err != nil {
			break
		}
	}
	fmt.Fprintln(file)

	// write points 
	points := d.Points()
	for i := 0; i < d.NumPoints(); i++ {
		indep := points[i].Indeps()
		depnd := points[i].Depnds()
		for j := 0; j < len(indep); j++ {
			_, err = fmt.Fprintf(file, "%f ", indep[j])
			if err != nil {
				break
			}
		}
		for j := 0; j < len(depnd); j++ {
			_, err = fmt.Fprintf(file, "%f ", depnd[j])
			if err != nil {
				break
			}
		}
		fmt.Fprintln(file)
	}
}

func SplitPointSetTrainTest(pnts *PointSet, pcnt_train float64, seed int) (train, test *PointSet) {

	train = new(PointSet)
	test = new(PointSet)

	train.filename, test.filename = pnts.filename, pnts.filename
	train.id, test.id = pnts.id, pnts.id
	train.indepNames, test.indepNames = pnts.indepNames, pnts.indepNames
	train.depndNames, test.depndNames = pnts.depndNames, pnts.depndNames
	train.sysNames, test.sysNames = pnts.sysNames, pnts.sysNames
	train.sysVals, test.sysVals = pnts.sysVals, pnts.sysVals

	L := len(pnts.dataPoints)
	Tst := int(float64(L) * (1.0 - pcnt_train))

	tmp := make([]Point, L)
	copy(tmp, pnts.dataPoints)

	rand.Seed(int64(seed))

	for i := 0; i < Tst; i++ {
		p := rand.Intn(L - i)
		tmp[i], tmp[p] = tmp[p], tmp[i]
	}

	test.dataPoints = tmp[:Tst]
	train.dataPoints = tmp[Tst:]

	return
}

func (d *PointSet) ReadLakeFile(filename string) {
	ftotal, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer ftotal.Close()
	file := bufio.NewReader(ftotal)

	var word string

	// get independent variables (x_i...)
	for i := 0; ; i++ {
		_, err := fmt.Fscanf(file, "%s", &word)
		if err != nil {
			break
		}
		d.indepNames = append(d.indepNames, word)
	}
	d.numDim = len(d.indepNames)

	// // get dependent variables (y_j...)
	// for i := 0; ; i++ {
	// 	_, err := fmt.Fscanf(file, "%s", &word)
	// 	if err != nil {
	// 		break
	// 	}
	// 	d.depndNames = append(d.depndNames, word)
	// }

	// remove time names from indepNames
	d.indepNames = d.indepNames[2:]

	// fmt.Printf("Var Names = %v | %v\n", d.depndNames, d.indepNames)

	for i := 0; ; i++ {
		var pnt Point
		var ival float64
		if err != nil {
			break
		}

		// read time values and disgaurd
		var dummy string
		for j := 0; j < 2; j++ {
			_, err = fmt.Fscanf(file, "%s", &dummy)
			// fmt.Println(i, dummy)
			if err != nil {
				// fmt.Println("err:", err)
				break
			}
		}

		// append dummy time value
		pnt.indep = append(pnt.indep, -0.1)

		// read real data
		for j := 0; j < len(d.indepNames); j++ {
			_, err = fmt.Fscanf(file, "%f\n", &ival)
			if err != nil {
				// fmt.Println("err:", err)
				break
			}
			// hack for PAR
			// if j == 1 {
			// 	ival /= 2000.0 * (math.Pi / 2.0)
			// }

			pnt.indep = append(pnt.indep, ival)
		}

		// fmt.Println(i, pnt)

		// for j := 0; j < len(d.depndNames); j++ {
		// 	_, err = fmt.Fscanf(file, "%f\n", &dval)
		// 	if err != nil {
		// 		break
		// 	}
		// 	pnt.depnd = append(pnt.depnd, dval)
		// }

		if len(pnt.indep) > 1 {
			for p := 0; p < len(pnt.indep); p++ {
				if math.IsNaN(pnt.indep[p]) {
					// fmt.Println("NaN @ ", i, p)
					goto skip
				}
			}
			d.dataPoints = append(d.dataPoints, pnt)
		}
	skip:
	}

	// calculate numerical derivatives
	calcDerivs(d.dataPoints)

	for i := 0; i < len(d.dataPoints); i++ {
		p := d.dataPoints[i]
		if len(p.Indeps()) == 0 || len(p.Depnds()) == 0 {
			fmt.Println("Bad Point @", i)
		}
		// if i%100 == 0 {
		// 	fmt.Printf("Point(%d): %v\n", i, d.dataPoints[i])
		// }
	}

	// fmt.Printf("Num Points: %v\n", len(d.dataPoints))
}

/* Calculate the first derivative of four points with: h = 0.25
 * (from: http://www.trentfguidry.net/post/2010/09/04/Numerical-differentiation-formulas.aspx)
 *
 * xF0 = ( -3.0*xF4 + 16.0*xF3 - 36.0*xF2 + 48.0*xF1 - 25.0*xF0) / (12.0*h)
 * xF1 = (      xF4 -  6.0-xF3 + 18.0*xF2 - 10.0*xF1 -  3.0*xF0) / (12.0*h)
 * xF2 = (     -xF4 +  8.0*xF3           -   8.0*xF1 +      xF0) / (12.0*h)
 * xF3 = (  3.0*xF4 + 10.0*xF3 - 18.0*xF2 +  6.0*xF1 -      xF0) / (12.0*h)
 * xF4 = ( 25.0*xF4 - 48.0*xF3 + 36.0*xF2 - 16.0*xF1 +  3.0*xF0) / (12.0*h)
 */
func calcDerivs(pts []Point) {
	h := 24.0 // / (24.0 * 60.0)

	NP := len(pts)
	ND := pts[0].NumIndep()

	// for summing and averaging on a point/variable~wise basis
	cnts := make([][]int, NP)
	vals := make([][]float64, NP)
	for p := 0; p < NP; p++ {
		cnts[p] = make([]int, ND)
		vals[p] = make([]float64, ND)
	}

	for p := 0; p < NP-5; p++ {
		for i := 0; i < ND; i++ {
			var F [5]float64
			var dF [5]float64

			for j := 0; j < 5; j++ {
				F[j] = pts[p+j].Indep(i)
			}
			dF[0] = (-3.0*F[4] + 16.0*F[3] - 36.0*F[2] + 48.0*F[1] - 25.0*F[0]) / (12.0 * h)
			dF[1] = (F[4] - 6.0 - F[3] + 18.0*F[2] - 10.0*F[1] - 3.0*F[0]) / (12.0 * h)
			dF[2] = (-F[4] + 8.0*F[3] - 8.0*F[1] + F[0]) / (12.0 * h)
			dF[3] = (3.0*F[4] + 10.0*F[3] - 18.0*F[2] + 6.0*F[1] - F[0]) / (12.0 * h)
			dF[4] = (25.0*F[4] - 48.0*F[3] + 36.0*F[2] - 16.0*F[1] + 3.0*F[0]) / (12.0 * h)

			for j := 0; j < 5; j++ {

				vals[p+j][i] += dF[j]
				cnts[p+j][i]++
			}
		}
	}

	for p := 0; p < NP; p++ {
		depnds := make([]float64, ND)
		for i := 0; i < ND; i++ {
			depnds[i] = vals[p][i] / float64(cnts[p][i])
		}
		pts[p].SetDepnds(depnds)
	}
}
