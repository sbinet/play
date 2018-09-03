package play

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func RunPy(code string) {
	dir, err := ioutil.TempDir("", "py-run-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "run.py")
	ioutil.WriteFile(fname, []byte(code), 0644)
	cmd := exec.Command("python3", fname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
