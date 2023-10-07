package slice

type sliceUtil struct {
}

var SliceUtil = new(sliceUtil)

func (s *sliceUtil) ContainsInt(a []int, value int) bool {
	if len(a) == 0 {
		return false
	}
	for _, item := range a {
		if item == value {
			return true
		}
	}
	return false
}
