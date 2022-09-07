package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
)

const barWidth = 68

func parseArg(s string) (time.Duration, error) {
	d, err := time.ParseDuration(os.Args[1])
	if err == nil {
		return d, nil
	}
	i, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err == nil {
		return time.Duration(i) * time.Second, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(f * float64(time.Second)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: sleepy NUMBER[SUFFIX]...")
		os.Exit(1)
	}

	d, err := parseArg(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := d.Seconds()
	uiprogress.Start()

	if s < 1 {
		bar := uiprogress.AddBar(1)
		bar.Width = barWidth
		bar.AppendCompleted()
		time.Sleep(time.Duration(s * float64(time.Second)))
		bar.Incr()
		uiprogress.Stop()
		return
	}

	bar := uiprogress.AddBar(int(math.Round(s)))
	bar.Width = barWidth
	bar.AppendCompleted()

	t := time.NewTicker(time.Second)

loop:
	for {
		select {
		case <-t.C:
			s--
			bar.Incr()
			if s < 1 {
				break loop
			}
		}
	}
	t.Stop()
	time.Sleep(time.Duration(s * float64(time.Second)))
	bar.Incr()
	uiprogress.Stop()
}
