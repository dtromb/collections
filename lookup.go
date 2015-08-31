package collections

// LookupType is the "direction" of the lookup to perform.
type LookupType int

const (
	// GTE searches for the smallest value >= search argument.
	GTE LookupType = iota
	// LTE searches for the largest value <= search argument.
	LTE
)
