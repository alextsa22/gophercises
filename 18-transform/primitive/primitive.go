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
func Transform(image io.Reader, ext string, numShapes int, opts ...func()) (io.Reader, error) {
	in, err := tempFile("__in_", ext)
	if err != nil {
		return nil, fmt.Errorf("tempFile(): failed to create temporary input fule: %v", err)
	}
	defer os.Remove(in.Name())
	out, err := tempFile("__out_", ext)
	if err != nil {
		return nil, fmt.Errorf("tempFile(): failed to create temporary output file: %v", err)
	}
	defer os.Remove(out.Name())

	if _, err = io.Copy(in, image); err != nil {
		return nil, fmt.Errorf("io.Copy(): failed to copy image into temp input file: %v", err)
	}

	// stdCombo
	_, err = primitive(in.Name(), out.Name(), numShapes, ModeCombo)
	if err != nil {
		return nil, fmt.Errorf("primitive(): failed to run the primitive command: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, out); err != nil {
		return nil, fmt.Errorf("io.Copy(): failed to copy output file into byte buffer: %v", err)
	}
	return buf, nil
}

func primitive(in, out string, numShapes int, mode Mode) (string, error) {
	args := fmt.Sprintf("-i %s -o %s -n %d -m %d", in, out, numShapes, mode)
	cmd := exec.Command("primitive", strings.Fields(args)...)
	bs, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("CombineOutput(): %v", err)
	}
	return string(bs), nil
}

func tempFile(prefix, ext string) (*os.File, error) {
	tmp, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	return os.Create(fmt.Sprintf("%s.%s", tmp.Name(), ext))
}
