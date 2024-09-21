package main

import "iter"

type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH

	FileFirst = FileA
	FileLast  = FileH

	FileCount = 8
)

func FileFromString(s string) (File, bool) {
	if len(s) != 1 || s[0] < 'a' || s[0] > 'h' {
		return 0, false
	}

	return File(s[0] - 'a'), true
}

func (f File) String() string {
	if f >= FileFirst && f <= FileLast {
		return string('a' + byte(f))
	}

	return repr(f)
}

func Files() iter.Seq[File] {
	return func(yield func(File) bool) {
		for file := FileFirst; file <= FileLast; file++ {
			if !yield(file) {
				break
			}
		}
	}
}
