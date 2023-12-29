package main

type Repetitions struct {
	_zobrists []uint64
	_indices  []int
	_index    int
}

const _RepetitionSize = 1024

func NewRepetitions() *Repetitions {
	return &Repetitions{
		_zobrists: make([]uint64, _RepetitionSize),
		_indices:  make([]int, _RepetitionSize+1),
	}
}

func (r *Repetitions) Push(zobrist uint64, reset bool) {
	r._zobrists[r._index] = zobrist
	r._indices[r._index+1] = ternary(reset, r._index, r._indices[r._index])

	r._index++
}

func (r *Repetitions) Pop() {
	if r._index > 0 {
		r._index--
	}
}

func (r *Repetitions) Contains(zobrist uint64) bool {
	for i := r._indices[r._index]; i < r._index-1; i++ {
		if r._zobrists[i] == zobrist {
			return true
		}
	}

	return false
}
