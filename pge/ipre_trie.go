package pge

type IpreNode struct {
	val int
	cnt int
	vst int

	next map[int]*IpreNode
}

func (n *IpreNode) InsertSerial(s []int) (did_ins bool) {
	in, _ := n.next[s[0]]

	// does this branch exist?
	if in == nil {
		in = new(IpreNode)
		in.val = s[0]
		in.next = make(map[int]*IpreNode)
		n.next[s[0]] = in
		did_ins = true
	}

	// recursive call to insert
	if len(s) > 1 {
		did_ins = in.InsertSerial(s[1:]) || did_ins
	}

	// visitation accounting  
	// (-1 is the root of the memoization tree)
	in.vst++
	if n.val == -1 {
		n.vst++
	}
	if did_ins {
		in.cnt++
		if n.val == -1 {
			n.cnt++
		}
	}

	return did_ins
}
