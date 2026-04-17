package hello

import "io"

// Echo reads from r and writes each line, unchanged, to w. It should
// terminate cleanly on EOF, not treat it as an error.
//
// Your task: implement this. The test will drive you.
func Echo(r io.Reader, w io.Writer) error {
	// TODO: implement
	_ = r
	_ = w
	return nil
}
