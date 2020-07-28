package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/boson-project/grid"
	"github.com/boson-project/grid/knative"
	"github.com/boson-project/grid/local"
	"github.com/boson-project/grid/platformb"
	"github.com/boson-project/grid/platformc"
)

var usage = `grid

Provide a unified, local interface to underlying serverless infrastructure in a manner befitting a utility grid.
`
var (
	Version = flag.Bool("version", false, "Print version [$GRID_VERSION]")
	Verbose = flag.Bool("verbose", false, "Print verbose logs [$GRID_VERBOSE]")
	Address = flag.String("address", grid.DefaultAddress, "Listen address [$GRID_ADDRESS]")
	Adapter = flag.String("adapter", "local", "Underlying architecture (knative|local) [$GRID_ADAPTER]")

	date string
	vers string
	hash string
)

func parseEnv() {
	parseBool("GRID_VERSION", Version)
	parseBool("GRID_VERBOSE", Verbose)
	parseString("GRID_ADDRESS", Address)
	parseString("GRID_ADAPTER", Adapter)
}

func printCfg() {
	fmt.Printf("GRID_VERSION=%v\n", *Version)
	fmt.Printf("GRID_VERBOSE=%v\n", *Verbose)
	fmt.Printf("GRID_ADDRESS=%v\n", *Address)
	fmt.Printf("GRID_ADAPTER=%v\n", *Adapter)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
	parseEnv()
	flag.Parse()

	if *Verbose {
		printCfg()
	}

	if *Version || (len(os.Args) > 1 && os.Args[1] == "version") {
		fmt.Println(version())
		return
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	fmt.Println("start", version())

	g := grid.New(
		grid.WithAddress(*Address),
		grid.WithVerbose(*Verbose),
		grid.WithAdapter(adapter()),
		grid.WithVersion(version()),
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Printf("%v\n", <-c)
		cancel()
	}()

	return g.Serve(ctx)
}

func adapter() grid.Adapter {
	// TODO: detect environment rather than use the flag.
	switch *Adapter {
	case "knative":
		return knative.NewAdapter()
	case "platformb":
		return platformb.NewAdapter()
	case "platformc":
		return platformc.NewAdapter()
	default:
		return local.NewAdapter()
	}
}

func parseBool(key string, value *bool) {
	if val, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			panic(err)
		}
		*value = b
	}
}

func parseString(key string, value *string) {
	if val, ok := os.LookupEnv(key); ok {
		*value = val
	}
}

func version() string {
	// If 'vers' is not a semver already, then the binary was built either
	// from an untagged git commit or was built directly from source
	// (set semver to v0.0.0)

	var elements = []string{}
	if strings.HasPrefix(vers, "v") {
		elements = append(elements, vers) // built via make with a tagged commit
	} else {
		elements = append(elements, "v0.0.0") // from source or untagged commit
	}

	if date != "" {
		elements = append(elements, date)
	}

	if hash != "" {
		elements = append(elements, hash)
	}

	return strings.Join(elements, "-")

}
