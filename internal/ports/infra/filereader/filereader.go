package filereader

import (
	"encoding/json"
	"io"
	"os"

	"github.com/guil95/ports/internal/ports"
)

type readFile struct{}

func NewReadFile() ports.FileReader {
	return &readFile{}
}

func (rf *readFile) ReadFileStream(path string, port chan ports.Ports, quit chan error) {
	defer close(port)

	file, err := os.Open(path)
	if err != nil {
		quit <- err
		return
	}

	dec := json.NewDecoder(file)

	_, err = dec.Token()
	if err != nil {
		quit <- err
	}

	for {
		var p ports.Ports

		_, err = dec.Token()
		if err == io.EOF {
			quit <- ports.EofError{}
		}
		if err != nil {
			quit <- err
		}

		if dec.More() {
			err = dec.Decode(&p)
			if err != nil {
				quit <- err
			}

			p.SetID()

			port <- p
		}
	}
}
