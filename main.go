package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: sleepy NUMBER[SUFFIX]...")
		os.Exit(1)
	}

	d, err := time.ParseDuration(os.Args[1])
	if err != nil {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		d = time.Duration(i) * time.Second
	}

	uiprogress.Start()
	bar := uiprogress.AddBar(int(d.Seconds()))
	bar.Width = 68
	bar.AppendCompleted()

	t := time.NewTicker(time.Second)

loop:
	for {
		select {
		case <-t.C:
			if !bar.Incr() {
				break loop
			}
		}
	}
	t.Stop()
	uiprogress.Stop()
}
