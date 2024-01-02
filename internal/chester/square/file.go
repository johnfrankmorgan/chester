package square

import "github.com/johnfrankmorgan/chester/internal/chester/util"

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

func (f File) String() string {
	if f.Valid() {
		return string('a' + byte(f))
	}

	return util.UnknownNumeric(f)
}

func (f File) Valid() bool {
	return f >= FileFirst && f <= FileLast
}
