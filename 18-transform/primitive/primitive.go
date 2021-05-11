package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Mode defines the shapes used whe transforming images.
type Mode int

// Modes supported by the primitive package.
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

// WithMode is an option for the Transform function that will define the
// mode you want to use. By default ModeTriangle will be used.
func WithMode(mode Mode) func() []string {
	if mode < 0 || mode > ModePolygon {
		mode = ModeTriangle
	}
	return func() []string {
		return []string{"-m", strconv.Itoa(int(mode))}
	}
}

// Transform will take the provided image and apply a primitive
// transformation to it, then return a reader to the resulting image.
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	in, err := tempFile("__in_", ext)
	if err != nil {
		return nil, fmt.Errorf("tempFile: failed to create temporary input fule: %v", err)
	}
	defer os.Remove(in.Name())
	out, err := tempFile("__out_", ext)
	if err != nil {
		return nil, fmt.Errorf("tempFile: failed to create temporary output file: %v", err)
	}
	defer os.Remove(out.Name())

	if _, err = io.Copy(in, image); err != nil {
		return nil, fmt.Errorf("io.Copy: failed to copy image into temp input file: %v", err)
	}

	_, err = primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, fmt.Errorf("primitive: failed to run the primitive command: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, out); err != nil {
		return nil, fmt.Errorf("io.Copy: failed to copy output file into byte buffer: %v", err)
	}
	return buf, nil
}

func primitive(in, out string, numShapes int, args ...string) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d", in, out, numShapes)
	args = append(strings.Fields(argStr), args...)
	cmd := exec.Command("primitive", args...)
	bs, err := cmd.CombinedOutput()
	return string(bs), err
}

func tempFile(prefix, ext string) (*os.File, error) {
	tmp, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	return os.Create(fmt.Sprintf("%s.%s", tmp.Name(), ext))
}
