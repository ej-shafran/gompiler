package location

type FileInfo struct {
	FileName string
	Contents string
}

type Location struct {
	FileInfo FileInfo
	Cursor   int
}

func (loc Location) LineAndOffset() (line int, offset int) {
	cursor := 0

	line = 1
	offset = 1
	for cursor < loc.Cursor {
		c := loc.FileInfo.Contents[cursor]

		if c == '\n' {
			line++
			offset = 0
		} else {
			offset++
		}

		cursor++
	}

	return line, offset
}
