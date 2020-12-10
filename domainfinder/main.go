package main

import (
	"log"
	"os"
	"os/exec"
)

var cmdChain = []*exec.Cmd{
	exec.Command("lib/synonyms"),
	exec.Command("lib/sprinkle"),
	exec.Command("lib/coolify"),
	exec.Command("lib/domainify"),
	exec.Command("lib/available"),
}

func main() {
	// domainfinerのStdinがsynonymsのStdinに
	cmdChain[0].Stdin = os.Stdin
	// domaindinerのStdoutがavailableのStdoutに
	cmdChain[len(cmdChain)-1].Stdout = os.Stdout

	// それぞれのプログラムの標準出力が、直後のプログラムの標準入力と接続される
	for i := 0; i < len(cmdChain)-1; i++ {
		thisCmd := cmdChain[i]
		nextCmd := cmdChain[i+1]
		stdout, err := thisCmd.StdoutPipe()
		if err != nil {
			log.Panicln(err)
		}

		nextCmd.Stdin = stdout
	}

	// それぞれのコマンドのStartメソッドを呼び出し、プログラムをパッググラウンドで実行
	for _, cmd := range cmdChain {
		if err := cmd.Start(); err != nil {
			log.Panicln(err)
		} else {
			// Startに成功した場合、コマンドのプロセスを終了
			defer cmd.Process.Kill()
		}
	}

	for _, cmd := range cmdChain {
		if err := cmd.Wait(); err != nil {
			log.Panicln(err)
		}
	}
}
