package btree

const (
	iRTypeExist    = 1
	iRTypeModified = 2
	iRTypeSplit    = 3

	dRType = 10
)

type insertResult struct {
	rtype    int // result type
	modified node
	left     node   // for split
	right    node   // for split
	pivot    []byte // pivot key, for split
}

type deleteResult struct {
	rtype int
}
