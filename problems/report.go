package problems

import (
	"container/list"
	"fmt"
	expr "github.com/verdverm/go-symexpr"
	"sort"
)

var MAX_BBQ_SIZE = 1 << 20

type SortType int

const (
	SORT_NULL SortType = iota

	GPSORT_SIZE

	GPSORT_PRE_HIT
	GPSORT_TRN_HIT
	GPSORT_TST_HIT
	GPSORT_PRE_ERR
	GPSORT_TRN_ERR
	GPSORT_TST_ERR

	GPSORT_SIZE_PRE_HIT
	GPSORT_SIZE_TRN_HIT
	GPSORT_SIZE_TST_HIT
	GPSORT_SIZE_PRE_ERR
	GPSORT_SIZE_TRN_ERR
	GPSORT_SIZE_TST_ERR

	GPSORT_PRE_HIT_SIZE
	GPSORT_TRN_HIT_SIZE
	GPSORT_TST_HIT_SIZE
	GPSORT_PRE_ERR_SIZE
	GPSORT_TRN_ERR_SIZE
	GPSORT_TST_ERR_SIZE

	GPSORT_PARETO_PRE_ERR
	GPSORT_PARETO_TRN_ERR
	GPSORT_PARETO_TST_ERR
	GPSORT_PARETO_PRE_HIT
	GPSORT_PARETO_TRN_HIT
	GPSORT_PARETO_TST_HIT

	PESORT_SIZE

	PESORT_PRE_HIT
	PESORT_TRN_HIT
	PESORT_TST_HIT
	PESORT_PRE_ERR
	PESORT_TRN_ERR
	PESORT_TST_ERR

	PESORT_SIZE_PRE_HIT
	PESORT_SIZE_TRN_HIT
	PESORT_SIZE_TST_HIT
	PESORT_SIZE_PRE_ERR
	PESORT_SIZE_TRN_ERR
	PESORT_SIZE_TST_ERR

	PESORT_PRE_HIT_SIZE
	PESORT_TRN_HIT_SIZE
	PESORT_TST_HIT_SIZE
	PESORT_PRE_ERR_SIZE
	PESORT_TRN_ERR_SIZE
	PESORT_TST_ERR_SIZE

	PESORT_PARETO_PRE_ERR
	PESORT_PARETO_TRN_ERR
	PESORT_PARETO_TST_ERR
	PESORT_PARETO_PRE_HIT
	PESORT_PARETO_TRN_HIT
	PESORT_PARETO_TST_HIT
)

type ExprReport struct {
	expr  expr.Expr
	coeff []float64

	// metrics
	size                             int
	predError, trainError, testError float64
	predScore, trainScore, testScore int

	// per data set metrics, if multiple data sets used
	predErrz  []float64
	predHitz  []int
	trainErrz []float64
	trainHitz []int
	testErrz  []float64
	testHitz  []int

	// ids
	uniqID int // unique ID among all exprs ?
	procID int // ID of the search process
	iterID int // iteration of the search
	unitID int // ID with respect to the search
	index  int // used internally in containers

	// production information
	// p1,p2 int  // parent IDs
	// method int // method that produced this expression
}

func (r *ExprReport) Size() int {
	if r.size == 0 {
		r.size = r.expr.Size()
	}
	return r.size
}

func (r *ExprReport) Clone() *ExprReport {
	ret := new(ExprReport)
	ret.expr = r.expr.Clone()
	ret.expr.CalcExprStats()
	ret.coeff = make([]float64, len(r.coeff))
	copy(ret.coeff, r.coeff)

	ret.predError = r.predError
	ret.trainError = r.trainError
	ret.testError = r.testError
	ret.predScore = r.predScore
	ret.trainScore = r.trainScore
	ret.testScore = r.testScore

	ret.predErrz = make([]float64, len(r.predErrz))
	copy(ret.predErrz, r.predErrz)
	ret.predHitz = make([]int, len(r.predHitz))
	copy(ret.predHitz, r.predHitz)
	ret.trainErrz = make([]float64, len(r.trainErrz))
	copy(ret.trainErrz, r.trainErrz)
	ret.trainHitz = make([]int, len(r.trainHitz))
	copy(ret.trainHitz, r.trainHitz)
	ret.testErrz = make([]float64, len(r.testErrz))
	copy(ret.testErrz, r.testErrz)
	ret.testHitz = make([]int, len(r.testHitz))
	copy(ret.testHitz, r.testHitz)

	ret.uniqID = r.uniqID
	ret.procID = r.procID
	ret.iterID = r.iterID
	ret.unitID = r.unitID

	return ret
}

var reportFormatString = `%v
coeff: %v
size: %4d   depth: %4d   
uId:  %4d   pId:   %4d
iId:  %4d   tId:   %4d

train NaNs: %4d
train L1:   %f

test  NaNs: %4d
test Evals: %4d   
test  L1:   %f
test  L2:   %f
`

func (r *ExprReport) String() string {

	return fmt.Sprintf(reportFormatString,
		r.expr, r.Coeff(),
		r.expr.Size(), r.expr.Height(),
		r.uniqID, r.procID, r.iterID, r.unitID,
		r.trainScore, r.trainError,
		r.predScore, r.testScore, r.testError, r.predError)
}
func (r *ExprReport) Latex(dnames, snames []string, cvals []float64) string {

	return fmt.Sprintf(reportFormatString,
		r.expr.PrettyPrint(dnames, snames, cvals), r.Coeff(),
		r.expr.Size(), r.expr.Height(),
		r.uniqID, r.procID, r.iterID, r.unitID,
		r.trainScore, r.testScore, r.trainError, r.testError, r.predError)
}

func (r *ExprReport) Expr() expr.Expr     { return r.expr }
func (r *ExprReport) SetExpr(e expr.Expr) { r.expr = e }

func (r *ExprReport) Coeff() []float64     { return r.coeff }
func (r *ExprReport) SetCoeff(c []float64) { r.coeff = c }

func (r *ExprReport) PredScore() int     { return r.predScore }
func (r *ExprReport) SetPredScore(s int) { r.predScore = s }

func (r *ExprReport) PredError() float64     { return r.predError }
func (r *ExprReport) SetPredError(e float64) { r.predError = e }

func (r *ExprReport) TrainScore() int     { return r.trainScore }
func (r *ExprReport) SetTrainScore(s int) { r.trainScore = s }

func (r *ExprReport) TrainError() float64     { return r.trainError }
func (r *ExprReport) SetTrainError(e float64) { r.trainError = e }

func (r *ExprReport) TestScore() int     { return r.testScore }
func (r *ExprReport) SetTestScore(s int) { r.testScore = s }

func (r *ExprReport) TestError() float64     { return r.testError }
func (r *ExprReport) SetTestError(e float64) { r.testError = e }

func (r *ExprReport) PredScoreZ() []int     { return r.predHitz }
func (r *ExprReport) SetPredScoreZ(s []int) { r.predHitz = s }

func (r *ExprReport) PredErrorZ() []float64     { return r.predErrz }
func (r *ExprReport) SetPredErrorZ(e []float64) { r.predErrz = e }

func (r *ExprReport) TrainScoreZ() []int     { return r.trainHitz }
func (r *ExprReport) SetTrainScoreZ(s []int) { r.trainHitz = s }

func (r *ExprReport) TrainErrorZ() []float64     { return r.trainErrz }
func (r *ExprReport) SetTrainErrorZ(e []float64) { r.trainErrz = e }

func (r *ExprReport) TestScoreZ() []int     { return r.testHitz }
func (r *ExprReport) SetTestScoreZ(s []int) { r.testHitz = s }

func (r *ExprReport) TestErrorZ() []float64     { return r.testErrz }
func (r *ExprReport) SetTestErrorZ(e []float64) { r.testErrz = e }

func (r *ExprReport) UniqID() int     { return r.uniqID }
func (r *ExprReport) SetUniqID(i int) { r.uniqID = i }

func (r *ExprReport) ProcID() int     { return r.procID }
func (r *ExprReport) SetProcID(i int) { r.procID = i }

func (r *ExprReport) IterID() int     { return r.iterID }
func (r *ExprReport) SetIterID(i int) { r.iterID = i }

func (r *ExprReport) UnitID() int     { return r.unitID }
func (r *ExprReport) SetUnitID(i int) { r.unitID = i }

// Array of ExprReports
type ExprReportArray []*ExprReport

func (p ExprReportArray) Len() int      { return len(p) }
func (p ExprReportArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ExprReportArray) Less(i, j int) bool {
	if p[i] == nil || p[i].expr == nil {
		return false
	}
	if p[j] == nil || p[j].expr == nil {
		return true
	}
	return p[i].expr.AmILess(p[j].expr)
}

type ExprReportArrayPredError struct {
	Array ExprReportArray
}

func (p ExprReportArrayPredError) Len() int      { return len(p.Array) }
func (p ExprReportArrayPredError) Swap(i, j int) { p.Array[i], p.Array[j] = p.Array[j], p.Array[i] }
func (p ExprReportArrayPredError) Less(i, j int) bool {
	if p.Array[i] == nil || p.Array[i].expr == nil {
		return false
	}
	if p.Array[j] == nil || p.Array[j].expr == nil {
		return true
	}
	return p.Array[i].predError < p.Array[j].predError
}

type ReportListNode struct {
	Rpt      *ExprReport
	prv, nxt *ReportListNode
}

func (self *ReportListNode) Next() *ReportListNode {
	return self.nxt
}
func (self *ReportListNode) Prev() *ReportListNode {
	return self.prv
}

type ReportList struct {
	length int
	Head   *ReportListNode
	Tail   *ReportListNode
}

func (self *ReportList) Len() int {
	return self.length
}

func (self *ReportList) Front() (node *ReportListNode) {
	return self.Head
}
func (self *ReportList) Back() (node *ReportListNode) {
	return self.Tail
}

func (self *ReportList) PushFront(rpt *ExprReport) {
	tmp := new(ReportListNode)
	tmp.Rpt = rpt
	tmp.nxt = self.Head
	if self.Head != nil {
		self.Head.prv = tmp
	}
	self.Head = tmp
	if self.Tail == nil {
		self.Tail = tmp
	}
	self.length++
}
func (self *ReportList) PushBack(rpt *ExprReport) {
	tmp := new(ReportListNode)
	tmp.Rpt = rpt
	tmp.prv = self.Tail
	if self.Tail != nil {
		self.Tail.nxt = tmp
	}
	self.Tail = tmp
	if self.Head == nil {
		self.Head = tmp
	}
	self.length++
}
func (self *ReportList) Remove(node *ReportListNode) {
	if node.prv != nil {
		node.prv.nxt = node.nxt
	} else {
		self.Head = node.nxt
	}
	if node.nxt != nil {
		node.nxt.prv = node.prv
	} else {
		self.Tail = node.prv
	}
	node.prv = nil
	node.nxt = nil
	self.length--
}

// Enhanced ExprReport Array with variable sorting methods
// NOTE!!!  the storage is reverse order
type ReportQueue struct {
	queue      []*ExprReport
	less       func(i, j *ExprReport) bool
	sortmethod SortType
}

func NewQueueFromArray(era ExprReportArray) *ReportQueue {
	Q := new(ReportQueue)
	Q.queue = era
	return Q
}

func NewReportQueue() *ReportQueue {
	B := new(ReportQueue)
	B.queue = make([]*ExprReport, 0, MAX_BBQ_SIZE)
	return B
}

func (bb ReportQueue) Len() int { return len(bb.queue) }
func (bb ReportQueue) Less(i, j int) bool {
	if bb.queue[i] == nil {
		return false
	}
	if bb.queue[j] == nil {
		return true
	}
	return bb.less(bb.queue[i], bb.queue[j])
}
func (bb ReportQueue) Swap(i, j int) {
	bb.queue[i], bb.queue[j] = bb.queue[j], bb.queue[i]
	if bb.queue[i] != nil {
		bb.queue[i].index = i
	}
	if bb.queue[j] != nil {
		bb.queue[j].index = j
	}
}
func (bb *ReportQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	// To simplify indexing expressions in these methods, we save a copy of the
	// slice object. We could instead write (*pq)[i].
	if len(bb.queue) == cap(bb.queue) {
		B := make([]*ExprReport, 0, len(bb.queue)*2)
		copy(B[:len(bb.queue)], (bb.queue)[:])
		bb.queue = B
	}
	a := bb.queue[:]
	n := len(a)
	a = a[0 : n+1]
	item := x.(*ExprReport)
	item.index = n
	a[n] = item
	bb.queue = a
}

func (bb *ReportQueue) Pop() interface{} {
	a := bb.queue
	n := len(a)
	item := a[n-1]
	item.index = -1 // for safety
	bb.queue = a[0 : n-1]
	return item
}

func (bb *ReportQueue) GetReport(pos int) (er *ExprReport) {
	return bb.queue[pos]
}
func (bb *ReportQueue) SetReport(pos int, er *ExprReport) {
	bb.queue[pos] = er
}
func (bb *ReportQueue) GetQueue() ExprReportArray {
	return bb.queue
}

func (bb *ReportQueue) SetSort(sortmethod SortType) {
	// all of the methods need to be reversed conceptually except for paretos
	// this is because the queue push/pops from the back
	bb.sortmethod = sortmethod

	switch sortmethod {
	case GPSORT_SIZE:
		bb.less = lessSize

	case GPSORT_PRE_HIT:
		bb.less = lessPredHits
	case GPSORT_TRN_HIT:
		bb.less = lessTrainHits
	case GPSORT_TST_HIT:
		bb.less = lessTestHits
	case GPSORT_PRE_ERR:
		bb.less = lessPredError
	case GPSORT_TRN_ERR:
		bb.less = lessTrainError
	case GPSORT_TST_ERR:
		bb.less = lessTestError

	case GPSORT_SIZE_PRE_HIT:
		bb.less = lessSizePredHits
	case GPSORT_SIZE_TRN_HIT:
		bb.less = lessSizeTrainHits
	case GPSORT_SIZE_TST_HIT:
		bb.less = lessSizeTestHits
	case GPSORT_SIZE_PRE_ERR:
		bb.less = lessSizePredError
	case GPSORT_SIZE_TRN_ERR:
		bb.less = lessSizeTrainError
	case GPSORT_SIZE_TST_ERR:
		bb.less = lessSizeTestError

	case GPSORT_PRE_HIT_SIZE:
		bb.less = lessPredHitsSize
	case GPSORT_TRN_HIT_SIZE:
		bb.less = lessTrainHitsSize
	case GPSORT_TST_HIT_SIZE:
		bb.less = lessTestHitsSize
	case GPSORT_PRE_ERR_SIZE:
		bb.less = lessPredErrorSize
	case GPSORT_TRN_ERR_SIZE:
		bb.less = lessTrainErrorSize
	case GPSORT_TST_ERR_SIZE:
		bb.less = lessTestErrorSize

	// we use more here because of the heap / queue
	case PESORT_SIZE:
		bb.less = moreSize

	case PESORT_PRE_HIT:
		bb.less = morePredHits
	case PESORT_TRN_HIT:
		bb.less = moreTrainHits
	case PESORT_TST_HIT:
		bb.less = moreTestHits
	case PESORT_PRE_ERR:
		bb.less = morePredError
	case PESORT_TRN_ERR:
		bb.less = moreTrainError
	case PESORT_TST_ERR:
		bb.less = moreTestError

	case PESORT_SIZE_PRE_HIT:
		bb.less = moreSizePredHits
	case PESORT_SIZE_TRN_HIT:
		bb.less = moreSizeTrainHits
	case PESORT_SIZE_TST_HIT:
		bb.less = moreSizeTestHits
	case PESORT_SIZE_PRE_ERR:
		bb.less = moreSizePredError
	case PESORT_SIZE_TRN_ERR:
		bb.less = moreSizeTrainError
	case PESORT_SIZE_TST_ERR:
		bb.less = moreSizeTestError

	case PESORT_PRE_HIT_SIZE:
		bb.less = morePredHitsSize
	case PESORT_TRN_HIT_SIZE:
		bb.less = moreTrainHitsSize
	case PESORT_TST_HIT_SIZE:
		bb.less = moreTestHitsSize
	case PESORT_PRE_ERR_SIZE:
		bb.less = morePredErrorSize
	case PESORT_TRN_ERR_SIZE:
		bb.less = moreTrainErrorSize
	case PESORT_TST_ERR_SIZE:
		bb.less = moreTestErrorSize

	default:
		// fmt.Printf("unknown sort method\n")
	}

}

func (bb *ReportQueue) Sort() {
	// all of the methods need to be reversed conceptually except for paretos
	// this is because the queue push/pops from the back
	switch bb.sortmethod {
	case GPSORT_PARETO_PRE_ERR:
		bb.GP_ParetoPredError()
		bb.reverseQueue()
	case GPSORT_PARETO_TRN_ERR:
		bb.GP_ParetoTrainError()
		bb.reverseQueue()
	case GPSORT_PARETO_TST_ERR:
		bb.GP_ParetoTestError()
		bb.reverseQueue()
	case GPSORT_PARETO_PRE_HIT:
		bb.GP_ParetoPredHits()
		bb.reverseQueue()
	case GPSORT_PARETO_TRN_HIT:
		bb.GP_ParetoTrainHits()
		bb.reverseQueue()
	case GPSORT_PARETO_TST_HIT:
		bb.GP_ParetoTestHits()
		bb.reverseQueue()

	case PESORT_PARETO_PRE_ERR:
		bb.PE_ParetoPredError()
	case PESORT_PARETO_TRN_ERR:
		bb.PE_ParetoTrainError()
	case PESORT_PARETO_TST_ERR:
		bb.PE_ParetoTestError()
	case PESORT_PARETO_PRE_HIT:
		bb.PE_ParetoPredHits()
	case PESORT_PARETO_TRN_HIT:
		bb.PE_ParetoTrainHits()
	case PESORT_PARETO_TST_HIT:
		bb.PE_ParetoTestHits()

	default:
		sort.Sort(bb)
	}

}

func (self *ReportQueue) Reverse() {
	self.reverseQueue()
}

func (bb *ReportQueue) reverseQueue() {
	i, j := 0, len(bb.queue)-1
	for i < j {
		bb.Swap(i, j)
		i++
		j--
	}
}

func lessSize(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func moreSize(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}

func lessPredHits(l, r *ExprReport) bool {
	sc := l.predScore - r.predScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTrainHits(l, r *ExprReport) bool {
	sc := l.trainScore - r.trainScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTestHits(l, r *ExprReport) bool {
	sc := l.testScore - r.testScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}

func lessPredError(l, r *ExprReport) bool {
	sc := l.predError - r.predError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTrainError(l, r *ExprReport) bool {
	sc := l.trainError - r.trainError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTestError(l, r *ExprReport) bool {
	sc := l.testError - r.testError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}

func lessSizePredHits(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.predScore - r.predScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessSizeTrainHits(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.trainScore - r.trainScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessSizeTestHits(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.testScore - r.testScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}

func lessSizePredError(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.predError - r.predError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessSizeTrainError(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.trainError - r.trainError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessSizeTestError(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.testError - r.testError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}

func lessPredHitsSize(l, r *ExprReport) bool {
	sc := l.predScore - r.predScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTrainHitsSize(l, r *ExprReport) bool {
	sc := l.trainScore - r.trainScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTestHitsSize(l, r *ExprReport) bool {
	sc := l.testScore - r.testScore
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}

func lessPredErrorSize(l, r *ExprReport) bool {
	sc := l.predError - r.predError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTrainErrorSize(l, r *ExprReport) bool {
	sc := l.trainError - r.trainError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}
func lessTestErrorSize(l, r *ExprReport) bool {
	sc := l.testError - r.testError
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.expr.AmILess(r.expr)
}

func morePredHits(l, r *ExprReport) bool {
	sc := l.predScore - r.predScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTrainHits(l, r *ExprReport) bool {
	sc := l.trainScore - r.trainScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTestHits(l, r *ExprReport) bool {
	sc := l.testScore - r.testScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func morePredError(l, r *ExprReport) bool {
	sc := l.predError - r.predError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTrainError(l, r *ExprReport) bool {
	sc := l.trainError - r.trainError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTestError(l, r *ExprReport) bool {
	sc := l.testError - r.testError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}

func moreSizePredHits(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.predScore - r.predScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreSizeTrainHits(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.trainScore - r.trainScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreSizeTestHits(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.testScore - r.testScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}

func moreSizePredError(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.predError - r.predError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreSizeTrainError(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.trainError - r.trainError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreSizeTestError(l, r *ExprReport) bool {
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.testError - r.testError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}

func morePredHitsSize(l, r *ExprReport) bool {
	sc := l.predScore - r.predScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTrainHitsSize(l, r *ExprReport) bool {
	sc := l.trainScore - r.trainScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTestHitsSize(l, r *ExprReport) bool {
	sc := l.testScore - r.testScore
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}

func morePredErrorSize(l, r *ExprReport) bool {
	sc := l.predError - r.predError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTrainErrorSize(l, r *ExprReport) bool {
	sc := l.trainError - r.trainError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}
func moreTestErrorSize(l, r *ExprReport) bool {
	sc := l.testError - r.testError
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	sz := l.Size() - r.Size()
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.expr.AmILess(r.expr)
}

func (bb *ReportQueue) GP_ParetoPredHits() {
	bb.less = lessSizePredHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.predScore
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.Size() > cSize {
				cSize = pb.Size()
				if pb.predScore < cScore {
					cScore = pb.predScore
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) GP_ParetoTrainHits() {
	bb.less = lessSizeTrainHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.trainScore
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.Size() > cSize {
				cSize = pb.Size()
				if pb.trainScore < cScore {
					cScore = pb.trainScore
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) GP_ParetoTestHits() {
	bb.less = lessSizeTestHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.testScore
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.Size() > cSize {
				cSize = pb.Size()
				if pb.testScore < cScore {
					cScore = pb.testScore
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) GP_ParetoPredError() {
	bb.less = lessSizePredError
	sort.Sort(bb)

	// var pareto list.List
	// pareto.Init()
	var pareto ReportList
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		// pb := pe.Value.(*ExprReport)
		pb := pe.Rpt
		cSize := pb.Size()
		cScore := pb.predError
		pe = pe.Next()
		for pe != nil && over >= 0 {
			// pb := pe.Value.(*ExprReport)
			pb := pe.Rpt
			sz := pb.Size()
			if sz > cSize {
				cSize = sz
				if pb.predError < cScore {
					cScore = pb.predError
					// bb.queue[over] = eLast.Value.(*ExprReport)
					bb.queue[over] = eLast.Rpt
					over--
					pareto.Remove(eLast)
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		// bb.queue[over] = eLast.Value.(*ExprReport)
		bb.queue[over] = eLast.Rpt
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) GP_ParetoTrainError() {
	bb.less = lessSizeTrainError
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.trainError
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.Size() > cSize {
				cSize = pb.Size()
				if pb.trainError < cScore {
					cScore = pb.trainError
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) GP_ParetoTestError() {
	bb.less = lessSizeTestError
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.testError
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.Size() > cSize {
				cSize = pb.Size()
				if pb.testError < cScore {
					cScore = pb.testError
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) PE_ParetoPredHits() {
	bb.less = lessSizePredHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.predScore
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.predScore > cScore {
				cScore = pb.predScore
				if pb.Size() > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.Size()
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) PE_ParetoTrainHits() {
	bb.less = lessSizeTrainHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.trainScore
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.trainScore > cScore {
				cScore = pb.trainScore
				if pb.Size() > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.Size()
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) PE_ParetoTestHits() {
	bb.less = lessSizeTestHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.testScore
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.testScore > cScore {
				cScore = pb.testScore
				if pb.Size() > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.Size()
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) PE_ParetoPredError() {
	bb.less = lessSizePredError
	sort.Sort(bb)

	// var pareto list.List
	// pareto.Init()
	var pareto ReportList
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		// pb := pe.Value.(*ExprReport)
		pb := pe.Rpt
		cSize := pb.Size()
		cScore := pb.predError
		pe = pe.Next()
		for pe != nil && over >= 0 {
			// pb = pe.Value.(*ExprReport)
			pb = pe.Rpt
			if pb.predError < cScore {
				cScore = pb.predError
				if pb.Size() > cSize {
					// bb.queue[over] = eLast.Value.(*ExprReport)
					bb.queue[over] = eLast.Rpt
					over--
					pareto.Remove(eLast)
					cSize = pb.Size()
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		// bb.queue[over] = eLast.Value.(*ExprReport)
		bb.queue[over] = eLast.Rpt
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) PE_ParetoTrainError() {
	bb.less = lessSizeTrainError
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.trainError
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.trainError < cScore {
				cScore = pb.trainError
				if pb.Size() > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.Size()
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) PE_ParetoTestError() {
	bb.less = lessSizeTestError
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		if bb.queue[i] == nil {
			continue
		}
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.Size()
		cScore := pb.testError
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.testError < cScore {
				cScore = pb.testError
				if pb.Size() > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.Size()
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}
