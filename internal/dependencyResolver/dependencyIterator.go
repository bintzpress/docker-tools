package dependencyResolver

type DependencyIterator struct {
	stack []*DependencyIteratorState
}

type DependencyIteratorState struct {
	position   *node
	childIndex int
	inChildren bool
}

func NewDependencyIterator() *DependencyIterator {
	is := new(DependencyIteratorState)
	is.position = dependencyTree
	is.childIndex = 0
	is.inChildren = false
	i := new(DependencyIterator)
	i.stack = append([]*DependencyIteratorState(i.stack), is)
	return i
}

func (i *DependencyIterator) DependencyIteratorNext() string {
	n := i.dependencyIteratorFindNext()
	if n == nil {
		return ""
	} else {
		return n.key
	}
}

func (i *DependencyIterator) dependencyIteratorFindNext() *node {
	is := i.stack[len(i.stack)-1]

	if is.childIndex == -1 {
		is.childIndex = 0
		return is.position
	} else {
		if !is.inChildren {
			if len(is.position.children) == 0 {
				// time to ascend
				if len(i.stack) > 1 { // more left after truncate
					i.stack = i.stack[:len(i.stack)-1]
					is = i.stack[len(i.stack)-1]
					is.childIndex++
					return i.dependencyIteratorFindNext()
				} else {
					return nil // at end of iterating
				}
			} else if is.childIndex >= len(is.position.children) {
				// done going through children next need to decend
				is.inChildren = true
				is.childIndex = 0

				var nis = new(DependencyIteratorState)
				nis.position = is.position.children[is.childIndex]
				nis.childIndex = 0 // don't want to visit node key again
				nis.inChildren = false
				i.stack = append([]*DependencyIteratorState(i.stack), nis)
				return i.dependencyIteratorFindNext()
			} else {
				// still going through children
				n := is.position.children[is.childIndex]
				is.childIndex++
				return n
			}
		} else {
			// descending into children
			if is.childIndex >= len(is.position.children) {
				// done descending into children
				// time to ascend
				if len(i.stack) > 1 { // more after truncate
					i.stack = i.stack[:len(i.stack)-1]
					is = i.stack[len(i.stack)-1]
					is.childIndex++
					return i.dependencyIteratorFindNext()
				} else {
					return nil // at end of iterating
				}
			} else {
				// descending into next child
				n := is.position.children[is.childIndex]
				is.childIndex++
				var nis = new(DependencyIteratorState)
				nis.position = n
				nis.childIndex = 0
				nis.inChildren = false
				i.stack = append([]*DependencyIteratorState(i.stack), nis)
				return i.dependencyIteratorFindNext()
			}
		}
	}
}
