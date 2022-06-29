package flowhash

import "errors"

type RSSContext uint32

func (c RSSContext) IsNew() bool {
	return c == ETH_RXFH_CONTEXT_ALLOC
}

const (
	ETH_RXFH_CONTEXT_ALLOC   = 0xffffffff
	ETH_RXFH_INDIR_NO_CHANGE = 0xffffffff
)

type Action interface {
	Fill(table IndirectTable) (int, error)
}

var (
	_ Action = (*Equal)(nil)
	_ Action = (*Weight)(nil)
	_ Action = (*Default)(nil)
	_ Action = (*Delete)(nil)
)

// Equal sets the receive flow hash indirection table to spread flows evenly between the first N receive queues.
type Equal struct {
	Start int // Sets the starting receive queue for spreading flows to N.
	N     int
}

func (e *Equal) Fill(table IndirectTable) (int, error) {

	for i := range table {
		table[i] = uint32(e.Start + i%e.N)
	}

	return len(table), nil
}

// Weight sets the receive flow hash indirection table to spread flows between receive queues according to the given weights.
// The sum of the weights must be non-zero and must not exceed the size of the indirection table.
type Weight struct {
	Start   int // Sets the starting receive queue for spreading flows to N.
	Weights []int
}

func (w *Weight) Fill(table IndirectTable) (int, error) {
	var sum int
	for _, n := range w.Weights {
		sum += n
	}

	if sum == 0 {
		return 0, errors.New("At least one weight must be non-zero")
	}

	if sum > len(table) {
		return 0, errors.New("Total weight exceeds the size of the indirection table")
	}

	var partial int

	j := -1
	for i := range table {
		for i >= len(table)*partial/sum {
			j += 1
			partial += w.Weights[j]
		}
		table[i] = uint32(w.Start + j)
	}

	return len(table), nil
}

// Default sets the receive flow hash indirection table to its default value.
type Default struct{}

func (d *Default) Fill(table IndirectTable) (int, error) {
	return 0, nil
}

// Delete the specified RSS context.
type Delete struct{}

func (d *Delete) Fill(table IndirectTable) (int, error) {
	return ETH_RXFH_INDIR_NO_CHANGE, nil
}
