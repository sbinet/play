// Copyright 2018 The play Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package play

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kr/pty"
)

// RunPy executes the given code through the python3 VM.
func RunPy(code string) {
	dir, err := ioutil.TempDir("", "py-run-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "run.py")
	ioutil.WriteFile(fname, []byte(code), 0644)
	err = runCmd(exec.Command("python3", fname))
	if err != nil {
		log.Fatal(err)
	}
}

// RunCxx compiles the given code with 'c++' and executes the
// resulting binary.
func RunCxx(code string, args ...string) {
	dir, err := ioutil.TempDir("", "cxx-run-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "out.cxx")
	oname := filepath.Join(dir, "out.exe")
	ioutil.WriteFile(fname, []byte(code), 0644)

	cxxArgs := []string{fname, "-o", oname}
	switch len(args) {
	case 0:
		cxxArgs = append(cxxArgs, "--std=c++11", "-pthread")
	default:
		cxxArgs = append(cxxArgs, args...)
	}

	cxx := exec.Command("c++", cxxArgs...)
	cxx.Stdout = os.Stdout
	cxx.Stderr = os.Stderr

	err = cxx.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = runCmd(exec.Command(oname))
	if err != nil {
		log.Fatal(err)
	}
}
func runCmd(cmd *exec.Cmd) error {
	tty, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer tty.Close()

	done := make(chan int)
	go func() {
		tick := time.NewTicker(10 * time.Millisecond)
		defer tick.Stop()
		buf := make([]byte, 1024)
		for {
			select {
			case <-done:
				io.Copy(os.Stdout, tty)
				done <- 1
				return
			case <-tick.C:
				n, _ := tty.Read(buf)
				os.Stdout.Write(buf[:n])
				os.Stdout.Sync()
			}
		}
	}()
	err = cmd.Wait()
	done <- 1
	<-done
	return err
}
