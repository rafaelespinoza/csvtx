package task_test

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/rafaelespinoza/csvtx/internal/task"
)

func TestMintToYNAB(t *testing.T) {
	baseDir, err := os.MkdirTemp("", strings.Replace(t.Name(), "/", "_", -1)+"_*")
	if err != nil {
		t.Fatal(err)
	}

	params := task.Params{
		Infile:  "fixtures/mint.csv",
		Outdir:  baseDir,
		LogDest: &testLogger{t},
	}
	err = task.MintToYNAB(params)
	if err != nil {
		t.Fatal(err)
	}

	readOutput(t, path.Join(baseDir, "checking.csv"))
	readOutput(t, path.Join(baseDir, "personal-savings.csv"))
	readOutput(t, path.Join(baseDir, "credit.csv"))
}

type testLogger struct{ t *testing.T }

func (w *testLogger) Write(in []byte) (n int, e error) {
	w.t.Logf("%s", in)
	n = len(in)
	return
}

func readOutput(t *testing.T, filename string) {
	t.Helper()

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("could not open %q", filename)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("could not read data from %q", filename)
	}
	t.Logf("-- %s\n%s", filename, data)
}
