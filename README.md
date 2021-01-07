[![Go Reference](https://pkg.go.dev/badge/pkg.go.dev/github.com/zrhmn/argv.svg)](https://pkg.go.dev/pkg.go.dev/github.com/zrhmn/argv)

argv is yet another command-line argument parser.
However, this one doesn't expect anything from you.
As a result, all of the responsibility of using this lies on you.

A simple example:

```go
package main

import (
  "os"
  "log"

  "github.com/zrhmn/argv"
)

func main() {
  parsed := argv.New(os.Args[1:]).Parse()
  infile := os.Stdin
  verbosity := 0

  for _, optarg := range parsed.OptArgs {
    switch optarg[0] {
      case "-f", "--file":
        f, err := os.Open(optarg[1])
        if err != nil {
          log.Fatalf("could not open file: %v", err)
        }

        defer f.Close()
        infile = f

      case "-v":
        if len(optarg[1]) != 0 {
          log.Fatal("option -v does not accept a value")
        }

        verbosity += 1

      // ... etc.
      
      default:
        log.Fatalf("option provided but not recognized: %s", optarg[0])
    }
  }

  // ... do something with gathered information.
}
```
