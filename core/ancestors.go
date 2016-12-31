package core

type ancestors []*File

func (a ancestors) push(value *File) ancestors {
	return append(a, value)
}

func (a ancestors) pop() (*File, ancestors) {
	if len(a) == 0 {
		return nil, a
	}
	newA := make(ancestors, len(a)-1)
	copy(newA, a[:len(a)-1])
	return a[len(a)-1], newA
}
