# sc2spall (stackcollapse to spall)

Converts the FlameGraph format (stackcollapse-*) to the Spall format.

Spall is an extremely fast profiler by Colin Davidson, available
at https://gravitymoth.com/spall/. For optimal file size and load times, Spall
has a proprietary binary format. This tool produces files in that format.

The FlameGraph format was created for Brendan Gregg's FlameGraph tool,
available at https://www.brendangregg.com/FlameGraphs/cpuflamegraphs.html.
A wide variety of stackcollapse-* scripts are available for different
languages, e.g. stackcollapse-perf, stackcollapse-chrome-tracing, and
stackcollapse-xdebug. This tool converts the collapsed format to Spall, so
it should be compatible with any of those tools.

## Installing

[Go 1.19](https://go.dev/) or higher is required. Make sure that `$GOBIN` or `$HOME/go/bin` is on your PATH.

```
go install github.com/bvisness/sc2spall@latest
```

## Usage

```
Usage:
  sc2spall [flags]

Flags:
  -h, --help         help for sc2spall
  -o, --out string   The filename to write to. For stdout, use "-".
```
