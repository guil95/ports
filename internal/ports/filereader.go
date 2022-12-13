package ports

type FileReader interface {
	ReadFileStream(path string, port chan Ports, quit chan error)
}

const eofErrorMessage = "eof error"

type EofError struct{}

func (ee EofError) Error() string {
	return eofErrorMessage
}
