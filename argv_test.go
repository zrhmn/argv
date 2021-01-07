package argv_test

import (
	"strings"
	"testing"

	"github.com/zrhmn/argv"
)

type _Case struct {
	args    []string
	posargs []string
	optargs [][2]string
}

func (tc _Case) Test(t *testing.T) {
	// t.Parallel()

	g := argv.New(tc.args).Parse()

	// that all expected positional arguments are parsed
	for _, arg := range tc.posargs {
		found := false
		for _, posarg := range g.PosArgs {
			if arg == posarg {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected %q in parsed positional arguments", arg)
		}
	}

	// that all parsed positional arguments were expected
	for _, posarg := range g.PosArgs {
		found := false
		for _, arg := range tc.posargs {
			if arg == posarg {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("did not expect %q in parsed positional arguments", posarg)
		}
	}

	// that all expected option arguments are parsed
	for _, arg := range tc.optargs {
		found := false
		for _, optarg := range g.OptArgs {
			if arg[0] == optarg[0] && arg[1] == optarg[1] {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected (%q = %q) in parsed option arguments", arg[0], arg[1])
		}
	}

	// that all parsed option arguments were expected
	for _, optarg := range g.OptArgs {
		found := false
		for _, arg := range tc.optargs {
			if arg[0] == optarg[0] && arg[1] == optarg[1] {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("did not expect (%q = %q) in parsed option arguments", optarg[0], optarg[1])
		}
	}
}

func TestArgvParse(t *testing.T) {
	cases := []_Case{
		{},
		{
			args:    strings.Split("pos args only", " "),
			posargs: []string{"pos", "args", "only"},
			optargs: nil,
		},
		{
			args:    strings.Split("--long options --only", " "),
			posargs: nil,
			optargs: [][2]string{
				{"--long", "options"},
				{"--only", ""},
			},
		},
		{
			args:    strings.Split("--long options --with - skip --optarg", " "),
			posargs: []string{"skip"},
			optargs: [][2]string{
				{"--long", "options"},
				{"--with", ""},
				{"--optarg"},
			},
		},
		{
			args:    strings.Split("--mixed --long -a nd -s hort --opt ions", " "),
			posargs: nil,
			optargs: [][2]string{
				{"--mixed", ""},
				{"--long", ""},
				{"-a", "nd"},
				{"-s", "hort"},
				{"--opt", "ions"},
			},
		},
		{
			args:    strings.Split("--options -a -n d pos args --mixed", " "),
			posargs: []string{"pos", "args"},
			optargs: [][2]string{
				{"--options", ""},
				{"-a", ""},
				{"-n", "d"},
				{"--mixed", ""},
			},
		},
		{
			args:    strings.Split("--long=args --with=equals and pos --args", " "),
			posargs: []string{"and", "pos"},
			optargs: [][2]string{
				{"--long", "args"},
				{"--with", "equals"},
				{"--args", ""},
			},
		},
		{
			args:    strings.Split("-s=hort -a=rgs with -e=quals and -p=osargs", " "),
			posargs: []string{"with", "and"},
			optargs: [][2]string{
				{"-s", "hort"},
				{"-a", "rgs"},
				{"-e", "quals"},
				{"-p", "osargs"},
			},
		},
		{
			args:    strings.Split("--long=and -s hort --args with -m=ixed equals", " "),
			posargs: []string{"equals"},
			optargs: [][2]string{
				{"--long", "and"},
				{"-s", "hort"},
				{"--args", "with"},
				{"-m", "ixed"},
			},
		},
		{
			args:    strings.Split("-sho rt --flag combo", " "),
			posargs: nil,
			optargs: [][2]string{
				{"-s", ""},
				{"-h", ""},
				{"-o", "rt"},
				{"--flag", "combo"},
			},
		},
		{
			args:    strings.Split("-sho rt --flag combo -w- ith --skip - optarg", " "),
			posargs: []string{"optarg", "ith"},
			optargs: [][2]string{
				{"-s", ""},
				{"-h", ""},
				{"-o", "rt"},
				{"--flag", "combo"},
				{"-w", ""},
				{"--skip", ""},
			},
		},
		{
			args:    strings.Split("-sho=rt --flag combo -w=ith -equal= s", " "),
			posargs: nil,
			optargs: [][2]string{
				{"-s", ""},
				{"-h", ""},
				{"-o", "rt"},
				{"--flag", "combo"},
				{"-w", "ith"},
				{"-e", ""},
				{"-q", ""},
				{"-u", ""},
				{"-a", ""},
				{"-l", "s"},
			},
		},
	}

	for _, _case := range cases {
		t.Run("", _case.Test)
	}
}
