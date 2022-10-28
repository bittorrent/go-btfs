package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	testfile        = "QmZHeNJTU4jFzgBAouHSqbT2tyYJxgk6i15e7x5pudBune"
	testfilecontent = "Hello BTFS!"
	timeoutSeconds  = 100
)

// we need to delete the file for get test from last run
func prepare_test(btfsBinaryPath, onlineServerDomain, peerId, hValue string) bool {
	cmd := exec.Command(btfsBinaryPath, "rm", testfile)
	err := cmd.Start()

	if err != nil {
		errMsg := fmt.Sprintf("btfs rm failed with message: [%v]", err)
		log.Errorf(errMsg)
		SendError(errMsg, onlineServerDomain, peerId, hValue)
		return false
	} else {
		log.Info("btfs test preparation succeed\n")
	}
	return true
}

func get_functest(btfsBinaryPath string) error {
	// btfs get file saved to current working directory
	dir, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("get working directory failed: [%v]", err))
	}

	cmd := exec.Command(btfsBinaryPath, "get", "-o", dir, testfile)

	go func() {
		time.Sleep(timeoutSeconds * time.Second)
		_, err := os.FindProcess(int(cmd.Process.Pid))
		if err != nil {
			log.Info("process already finished\n")
		} else {
			err := cmd.Process.Kill()
			if err != nil {
				if !strings.Contains(err.Error(), "process already finished") {
					log.Errorf("cannot kill process: [%v] \n", err)
				}
			}
		}
	}()

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err = cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("btfs get test failed: [%v], [%s]", err, errbuf.String()))
	}

	data, err := ioutil.ReadFile(dir + "/" + testfile)
	if err != nil {
		log.Errorf("btfs get test: read file failed: [%v]\n", err)
		return errors.New(fmt.Sprintf("btfs get test: read file failed: [%v]", err))
	}

	// remote last "\n" before compare
	if string(data[:len(data)-1]) != testfilecontent {
		log.Errorf("btfs get test: get different content[%s]\n", string(data))
		return errors.New(fmt.Sprintf("btfs get test: get different content[%s]", string(data)))
	}

	return nil
}

func add_functest(btfsBinaryPath, peerId string) error {
	// write btfs peerId and current time to a file in current working directory
	// then btfs add that file for test
	dir, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("get working directory failed: [%v]", err))
	}

	out := []byte(peerId + "\n")
	// add current time stamp to file content so every time adding-file hash is different
	currentTime := time.Now().String()
	out = append(out, currentTime...)

	origin := out
	filename := dir + "/btfstest.txt"
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("btfs add test: write file failed: [%v]", err))
	}

	cmd := exec.Command(btfsBinaryPath, "add", filename)

	go func() {
		time.Sleep(timeoutSeconds * time.Second)
		_, err := os.FindProcess(int(cmd.Process.Pid))
		if err != nil {
			log.Info("process already finished\n")
		} else {
			err := cmd.Process.Kill()
			if err != nil {
				if !strings.Contains(err.Error(), "process already finished") {
					log.Errorf("cannot kill process: [%v] \n", err)
				}
			}
		}
	}()

	out, err = cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("btfs add test failed: [%v]", err))
	}

	s := strings.Split(string(out), " ")
	if len(s) < 2 {
		return errors.New(fmt.Sprintf("btfs add test failed: invalid add result[%s]", string(out)))
	}

	addfilehash := s[1]
	// btfs get the file to compare with original file content
	cmd = exec.Command(btfsBinaryPath, "get", "-o", dir, addfilehash)
	out, err = cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("btfs add test failed: get file error [%v]\n", err))
	}

	path := dir + "/" + addfilehash
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("btfs add test: read file failed: [%v]\n", err)
		return errors.New(fmt.Sprintf("btfs add test: read file failed: [%v]", err))
	}

	if string(data) != string(origin) {
		return errors.New(fmt.Sprintf("btfs add test failed: get different content, btfs add file:[%s], btfs get file:[%s]",
			string(origin), string(data)))
	}

	// clean up
	err = os.Remove(path)
	if err != nil {
		// if not deleted, no need to fail the whole function
		log.Errorf("btfs add test: clean up file failed: [%v]\n", err)
	}

	return nil
}
