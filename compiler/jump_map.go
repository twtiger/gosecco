package compiler

type jumpMap struct {
	labelToPosition map[label][]int
	positionToLabel map[int]label
}

func createJumpMap() *jumpMap {
	return &jumpMap{
		labelToPosition: make(map[label][]int),
		positionToLabel: make(map[int]label),
	}
}

func (j *jumpMap) registerJump(l label, i int) {
	j.labelToPosition[l] = append(j.labelToPosition[l], i)
	// If the jump maps are used correctly, this should never overwrite an index
	j.positionToLabel[i] = l
}

func (j *jumpMap) allJumpsTo(l label) []int {
	return j.labelToPosition[l]
}

func (j *jumpMap) hasJumpFrom(i int) bool {
	_, ok := j.positionToLabel[i]
	return ok
}

func (j *jumpMap) shift(from int, incr int) {
	newLabelToPosition := make(map[label][]int)
	newPositionToLabel := make(map[int]label)

	for l, v := range j.labelToPosition {
		for _, ix := range v {
			if ix > from {
				ix += incr
			}

			newLabelToPosition[l] = append(newLabelToPosition[l], ix)
			newPositionToLabel[ix] = l
		}
	}

	j.labelToPosition = newLabelToPosition
	j.positionToLabel = newPositionToLabel
}

func jumpMapFrom(m map[label][]int) *jumpMap {
	jm := createJumpMap()
	for l, vs := range m {
		for _, v := range vs {
			jm.registerJump(l, v)
		}
	}
	return jm
}
