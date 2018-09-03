package play

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func RunCxx(code string, args ...string) {
	dir, err := ioutil.TempDir("", "cxx-run-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "out.cxx")
	oname := filepath.Join(dir, "out.exe")
	ioutil.WriteFile(fname, []byte(code), 0644)

	cxx := exec.Command("c++", fname, "-o", oname, "--std=c++11", "-pthread")
	cxx.Stdout = os.Stdout
	cxx.Stderr = os.Stderr

	err = cxx.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(oname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
