package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bvisness/spall-go"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sc2spall",
	Run: func(cmd *cobra.Command, args []string) {
		p := spall.NewProfile(os.Stdout, spall.UnitMilliseconds)
		defer p.Close()
		e := p.NewEventer()
		defer e.Close()

		var currentStack []string
		var now float64 = 0

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			stack := strings.TrimRight(line, "0123456789")
			countStr := line[len(stack):]

			stackEntries := strings.Split(stack, ";")

			for i := 0; i < len(stackEntries); i++ {
				entry := stackEntries[i]
				if i < len(currentStack) && currentStack[i] != entry {
					// Different entry - end everything past this point
					for j := len(currentStack) - 1; j >= i; j-- {
						e.End(now)
					}
					currentStack = currentStack[:i]
				}
				if i >= len(currentStack) {
					// New stack entries; begin these events
					e.Begin(entry, now)
					currentStack = append(currentStack, entry)
				}
			}

			count, err := strconv.Atoi(countStr)
			if err != nil {
				panic(fmt.Errorf("'%s' is not a valid sample count", countStr))
			}
			now += float64(count)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		// Pop the remaining items
		for i := len(currentStack) - 1; i >= 0; i-- {
			e.End(now)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
