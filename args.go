package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// headerArgs is the type used to store the header arguments
type headerArgs []string

func (h *headerArgs) Set(val string) error {
	*h = append(*h, val)
	return nil
}

func (h headerArgs) String() string {
	return "string"
}

type saveStatusArgs []int

func (s *saveStatusArgs) Set(val string) error {
	i, _ := strconv.Atoi(val)
	*s = append(*s, i)
	return nil
}

func (s saveStatusArgs) String() string {
	return "string"
}

func (s saveStatusArgs) Includes(search int) bool {
	for _, status := range s {
		if status == search {
			return true
		}
	}
	return false
}

type config struct {
	body           string
	concurrency    int
	delay          int
	discResp       string
	headers        headerArgs
	jsonOut        bool
	followLocation bool
	method         string
	saveResp       string
	saveStatus     saveStatusArgs
	timeout        int
	verbose        bool

	noHeaders bool
	paths     string
	hosts     string
	output    string

	requester requester
}

func processArgs() config {

	// body param
	body := ""
	flag.StringVar(&body, "body", "", "")
	flag.StringVar(&body, "b", "", "")

	// concurrency param
	concurrency := 20
	flag.IntVar(&concurrency, "concurrency", 20, "")
	flag.IntVar(&concurrency, "c", 20, "")

	// delay param
	delay := 5000
	flag.IntVar(&delay, "delay", 5000, "")
	flag.IntVar(&delay, "d", 5000, "")

	// discResp params
	discResp := ""
	flag.StringVar(&discResp, "discresp", "", "")
	flag.StringVar(&discResp, "dr", "", "")

	// headers params
	var headers headerArgs
	flag.Var(&headers, "header", "")
	flag.Var(&headers, "H", "")

	// Json Output param
	jsonOut := false
	flag.BoolVar(&jsonOut, "jsonOut", false, "")
	flag.BoolVar(&jsonOut, "j", false, "")

	// follow location param
	followLocation := false
	flag.BoolVar(&followLocation, "location", false, "")
	flag.BoolVar(&followLocation, "L", false, "")

	// rawhttp param
	rawHTTP := false
	flag.BoolVar(&rawHTTP, "rawhttp", false, "")
	flag.BoolVar(&rawHTTP, "r", false, "")

	// saveResp params
	saveResp := ""
	flag.StringVar(&saveResp, "saveresp", "", "")
	flag.StringVar(&saveResp, "sr", "", "")

	// savestatus params
	var saveStatus saveStatusArgs
	flag.Var(&saveStatus, "savestatus", "")
	flag.Var(&saveStatus, "s", "")

	// timeout param
	timeout := 10000
	flag.IntVar(&timeout, "timeout", 10000, "")
	flag.IntVar(&timeout, "t", 10000, "")

	// verbose param
	verbose := false
	flag.BoolVar(&verbose, "verbose", false, "")
	flag.BoolVar(&verbose, "v", false, "")

	// method param
	method := "GET"
	flag.StringVar(&method, "method", "GET", "")
	flag.StringVar(&method, "X", "GET", "")

	// no headers
	noHeaders := false
	flag.BoolVar(&noHeaders, "no-headers", false, "")

	flag.Parse()

	// paths might be in a file, or it might be a single value
	paths := flag.Arg(0)
	if paths == "" {
		paths = defaultPathsFile
	}

	// hosts are always in a file
	hosts := flag.Arg(1)
	if hosts == "" {
		hosts = defaultHostsFile
	}

	// default the output directory to ./out
	output := flag.Arg(2)
	if output == "" {
		output = defaultOutputDir
	}

	// set the requester function to use
	requesterFn := goRequest
	if rawHTTP {
		requesterFn = rawRequest
	}

	return config{
		body:           body,
		concurrency:    concurrency,
		delay:          delay,
		discResp:       discResp,
		headers:        headers,
		jsonOut:        jsonOut,
		followLocation: followLocation,
		method:         method,
		saveResp:       saveResp,
		saveStatus:     saveStatus,
		timeout:        timeout,
		requester:      requesterFn,
		verbose:        verbose,
		noHeaders:      noHeaders,
		paths:          paths,
		hosts:          hosts,
		output:         output,
	}
}

func init() {
	flag.Usage = func() {
		h := "Request many paths for many hosts\n\n"

		h += "Usage:\n"
		h += "  meg [path|pathsFile] [hostsFile] [outputDir]\n\n"

		h += "Options:\n"
		h += "  -b,  --body <val>           Set the request body\n"
		h += "  -c,  --concurrency <val>    Set the concurrency level (default: 20)\n"
		h += "  -d,  --delay <millis>       Milliseconds between requests to the same host (default: 5000)\n"
		h += "  -dr, --discresp <string>    Discard responses containing specific string\n"
		h += "  -H,  --header <header>      Send a custom HTTP header\n"
		h += "  -j,  --jsonOut              Output in JSON\n\n"
		h += "  -L,  --location             Follow redirects / location header\n"
		h += "  -r,  --rawhttp              Use the rawhttp library for requests (experimental)\n"
		h += "  -s,  --savestatus <status>  Save only responses with specific status code\n"
		h += "  -sr, --saveresp <string>    Save only responses containing specific string\n"
		h += "  -t,  --timeout <millis>     Set the HTTP timeout (default: 10000)\n"
		h += "  -v,  --verbose              Verbose mode\n"
		h += "  -X,  --method <method>      HTTP method (default: GET)\n\n"
		h += "       --no-headers           Don't output headers\n"

		h += "Defaults:\n"
		h += "  pathsFile: ./paths\n"
		h += "  hostsFile: ./hosts\n"
		h += "  outputDir:  ./out\n\n"

		h += "Paths file format:\n"
		h += "  /robots.txt\n"
		h += "  /package.json\n"
		h += "  /security.txt\n\n"

		h += "Hosts file format:\n"
		h += "  http://example.com\n"
		h += "  https://example.edu\n"
		h += "  https://example.net\n\n"

		h += "Examples:\n"
		h += "  meg /robots.txt\n"
		h += "  meg paths.txt hosts.txt output\n"

		fmt.Fprintf(os.Stderr, h)
	}
}
