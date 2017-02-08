package btree

const (
	iRTypeExist    = 1
	iRTypeModified = 2
	iRTypeSplit    = 3

	dRTypeNotPresent      = 11
	dRTypeRemoved         = 12
	dRTypeBorrowFromLeft  = 13
	dRTypeBorrowFromRight = 14
	dRTypeMergeWithLeft   = 15
	dRTypeMergeWithRight  = 16
)

type insertResult struct {
	rtype    int // result type
	modified node
	left     node   // for split
	right    node   // for split
	pivot    []byte // pivot key, for split
}

type deleteResult struct {
	rtype           int
	modified        node
	modifiedSibling node
}
