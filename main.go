package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
)

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

	d, err := parseArg(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := d.Milliseconds()
	w := int(s)
	if w < 1 {
		w = 1
	}
	uiprogress.Start()

	bar := uiprogress.AddBar(w / 100)
	bar.Width = 68
	bar.AppendCompleted()

	t := time.NewTicker(100 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
loop:
	for {
		select {
		case <-t.C:
			bar.Incr()
		case <-ctx.Done():
			break loop
		}
	}
	t.Stop()
	bar.Set(w)
	uiprogress.Stop()
}
