// Package argv implements a simple GNU-style command-line argument parser.
//
// There are several other packages providing similar functionality, but most
// of them require a setup phase where a specification for the expected
// arguments is provided. In contrast, argv simply goes through a list of
// command-line arguments and classifies them as either positional arguments,
// flags or option-value pairs.
//
// Simply do:
//     parsed := New(os.Args[1:]).Parse()
//
// This simplicity has some limitations, for example not knowing whether an
// argument is a boolean flag or a list of strings before parsing leads to
// unexpected results. However, this can still be used to quickly parse and
// process a very simple list of arguments.
package argv

import "strings"

// Argv wraps a list of command-line arguments and facilitates parsing them.
type Argv struct {
	args []string // unparsed args

	PosArgs []string    // parsed positional arguments
	OptArgs [][2]string // parsed option:value pairs
}

// New returns a new Argv object.
func New(args []string) *Argv {
	if args == nil {
		args = []string{}
	}
	return &Argv{args, []string{}, [][2]string{}}
}

// Parse parses the command-line arguments contained in Argv.
func (g *Argv) Parse() *Argv {

	// g.Parse is a thin wrapper around g.parse just so it can be chained, e.g.
	//     g := New(os.Args).Parse()
	// etc.

	g.parse()
	return g
}

// parse contains logic for parsing the arguments in g.args.
//
// Each argument is picked off and determined to be either
// (1) an option,
// (2) a combination of short flags,
// (3) an option=value pair,
// (4) a value for an option immediately preceding it OR
// (5) a positional argument.
//
// The parsed arguments are stored in g.PosArgs and g.OptArgs.
func (g *Argv) parse() {
	lastopt := ""
	for arg := g.next(); arg != "--"; arg = g.next() {
		switch {
		case len(arg) == 0:

		// arg is an option (-o, -ovalue, -o=value, --option, --option=value etc.)
		case arg[0] == '-':
			// last option is expecting value
			if len(lastopt) != 0 {
				lastopt, g.OptArgs = "", append(g.OptArgs, [2]string{lastopt, ""})

				// arg is "-", no need to parse further
				if len(arg) == 1 {
					continue
				}
			}

			// --option or --option=value
			if arg[1] == '-' {
				if optval := strings.SplitN(arg, "=", 2); len(optval) == 2 {
					g.OptArgs = append(g.OptArgs, [2]string{optval[0], optval[1]})
				} else {
					lastopt = arg
				}

				continue
			}

			// -o or -o=val or -opqr=val
			i, j := 1, strings.IndexByte(arg, '=')
			if j < 0 {
				j = len(arg)
			}

			// -opqr=val is treated like -o -p -q -r=val
			optb := [2]byte{'-', 0}
			for ; i < j-1; i++ {
				optb[1] = arg[i]
				g.OptArgs = append(g.OptArgs, [2]string{string(optb[:]), ""})
			}

			optb[1] = arg[j-1]
			if j < len(arg)-1 {
				g.OptArgs = append(g.OptArgs, [2]string{string(optb[:]), arg[j+1:]})
			} else {
				lastopt = string(optb[:])
			}

		// arg is value for last option
		case len(lastopt) != 0:
			// last option was a short-option combo ending in - (e.g. -pqr-)
			if lastopt == "--" {
				lastopt = ""
				g.PosArgs = append(g.PosArgs, arg)
			} else {
				// last option expecting value
				lastopt, g.OptArgs = "", append(g.OptArgs, [2]string{lastopt, arg})
			}

		// arg is a positional argument
		default:
			g.PosArgs = append(g.PosArgs, arg)
		}
	}

	// last --option not stored yet
	if len(lastopt) != 0 {
		g.OptArgs = append(g.OptArgs, [2]string{lastopt, ""})
	}
}

// next slices off the first (next to be processed) member from g.args and
// returns it. If g.args is empty, next returns "--".
func (g *Argv) next() (v string) {
	if len(g.args) == 0 {
		return "--"
	}

	v, g.args = g.args[0], g.args[1:]
	return
}
