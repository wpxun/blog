package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 可以用这个方法进行简单的测试，运行时要 build 成执行文件直接运行，不要用 go run 间接（包装）运行
func test() {

	//可以直接 kill -s sigkill 58148 或 kill -n 9 58148 方式发送信号
	//ctrl-c 发送 SIGINT 信号给前台进程组中的所有进程。常用于终止正在运行的程序。
	//ctrl-z 发送 SIGTSTP 信号给前台进程组中的所有进程，常用于挂起一个进程。
	//ctrl-\  发送 SIGQUIT
	//ctrl-d 不是发送信号，而是表示一个特殊的二进制值，表示 EOF。
	//SIGKILL 和 SIGSTOP 两种信息不能自行处理

	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR2, syscall.SIGTSTP}
	fmt.Printf("Set notification for %s... [sigRecv1]\n", sigs1)
	signal.Notify(sigRecv1, sigs1...)
	for sig := range sigRecv1 {
		fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
	}
}

func main() {
	go func() {
		time.Sleep(5 * time.Second)
		sendSignal()
	}()
	handleSignal()
}

func handleSignal() {
	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv1]\n", sigs1)
	signal.Notify(sigRecv1, sigs1...)
	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv2]\n", sigs2)
	signal.Notify(sigRecv2, sigs2...)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for sig := range sigRecv1 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv1]\n")
		wg.Done()
	}()
	go func() {
		for sig := range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv2]\n")
		wg.Done()
	}()

	fmt.Println("Wait for 2 seconds... ")
	time.Sleep(2 * time.Second)
	fmt.Printf("Stop notification...")
	signal.Stop(sigRecv1)
	close(sigRecv1)
	fmt.Printf("done. [sigRecv1]\n")
	wg.Wait()
}

func sendSignal() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal Error: %s\n", err)
			debug.PrintStack()
		}
	}()
	// ps aux | grep "signal" | grep -v "grep" | grep -v "go run" | awk '{print $2}'
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "signal"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{print $2}"),
	}
	output, err := runCmds(cmds)
	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}
	pids, err := getPids(output)
	if err != nil {
		fmt.Printf("PID Parsing Error: %s\n", err)
		return
	}
	fmt.Printf("Target PID(s):\n%v\n", pids)
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Process Finding Error: %s\n", err)
			return
		}
		sig := syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process (pid=%d)...\n", sig, pid)
		err = proc.Signal(sig)
		if err != nil {
			fmt.Printf("Signal Sending Error: %s\n", err)
			return
		}
	}
}

func getPids(strs []string) ([]int, error) {
	var pids []int
	for _, str := range strs {
		pid, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func runCmds(cmds []*exec.Cmd) ([]string, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("The cmd slice is invalid!")
	}
	first := true
	var output []byte
	var err error
	for _, cmd := range cmds {
		fmt.Printf("Run command: %v\n", getCmdPlaintext(cmd))
		if !first {
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}
		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}
		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}
		output = stdoutBuf.Bytes()
		//fmt.Printf("Output:\n%s\n", string(output))
		if first {
			first = false
		}
	}
	var lines []string
	var outputBuf bytes.Buffer
	outputBuf.Write(output)
	for {
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, getError(err, nil)
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

func getCmdPlaintext(cmd *exec.Cmd) string {
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)
	for _, arg := range cmd.Args[1:] {
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}
	return buf.String()
}

func getError(err error, cmd *exec.Cmd, extraInfo ...string) error {
	var errMsg string
	if cmd != nil {
		errMsg = fmt.Sprintf("%s  [%s %v]", err, (*cmd).Path, (*cmd).Args)
	} else {
		errMsg = fmt.Sprintf("%s", err)
	}
	if len(extraInfo) > 0 {
		errMsg = fmt.Sprintf("%s (%v)", errMsg, extraInfo)
	}
	return errors.New(errMsg)
}
