package input

import "io"

type Parser interface {
	Parse(r io.Reader, w io.Writer) error
}
