package main

import (
	"os/exec"
)

var recordingMap = make(map[string]*exec.Cmd)
