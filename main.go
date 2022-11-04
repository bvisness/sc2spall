package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bvisness/spall-go"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sc2spall",
		Short: "Converts the FlameGraph format to Spall.",
		Long: `Converts the FlameGraph format (stackcollapse-*) to the Spall format.

Spall is an extremely fast profiler by Colin Davidson, available
at https://gravitymoth.com/spall/. For optimal file size and load times, Spall
has a proprietary binary format. This tool produces files in that format.

The FlameGraph format was created for Brendan Gregg's FlameGraph tool,
available at https://www.brendangregg.com/FlameGraphs/cpuflamegraphs.html.
A wide variety of stackcollapse-* scripts are available for different
languages, e.g. stackcollapse-perf, stackcollapse-chrome-tracing, and
stackcollapse-xdebug. This tool converts the collapsed format to Spall, so
it should be compatible with any of those tools.`,
		Run: func(cmd *cobra.Command, args []string) {
			var f io.Writer = os.Stdout
			if out, err := cmd.PersistentFlags().GetString("out"); err == nil && out != "" {
				if out == "-" {
					f = os.Stdout
				} else {
					f, err = os.Create(out)
					if err != nil {
						panic(err)
					}
				}
			}

			p := spall.NewProfile(f, spall.UnitMilliseconds)
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
	rootCmd.PersistentFlags().StringP("out", "o", "", "The filename to write to. For stdout, use \"-\".")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
