package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner error
	Message string
	StackTrace string
	Misc map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:err,
		Message:fmt.Sprintf(messagef, msgArgs...),
		StackTrace:string(debug.Stack()),
		Misc:make(map[string]interface{}),
	}
}

func (err MyError) Error () string {
	return err.Message
}

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{wrapError(err, err.Error())}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)

	if err != nil {
		return err
	} else if isExecutable == false {
		return wrapError(nil, "job binary is not executable")
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

var myLog = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile|log.LUTC)

func handleError(key int, err error, message string)  {
	myLog.SetPrefix(fmt.Sprintf("[logID: %v]:", key))
	myLog.Printf("%#v\n", err)
	myLog.Output(2, fmt.Sprintf("%v\n", "12313"))
	//fmt.Printf("[%v] %v", key, message)
}

func main()  {
	err := runJob("1")

	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}

