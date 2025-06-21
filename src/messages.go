package src

// Messages
type brightnessMsg struct {
	brightness int
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}
