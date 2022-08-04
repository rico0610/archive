package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func Run(try int, name string, arg ...string) (*bytes.Buffer, error) {
	var err error
	for i := 0; i < try; i++ {
		log.Println("running round:", i, name, arg)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command(name, arg...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.Println(fmt.Sprint(err), stderr.String())
			continue
		}
		return &stdout, nil
	}
	return nil, err
}