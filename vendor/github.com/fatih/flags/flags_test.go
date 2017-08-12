package flags

import (
	"reflect"
	"testing"
)

func TestHasFlag(t *testing.T) {
	var flags = []struct {
		flag    string
		args    []string
		hasFlag bool
	}{
		{args: []string{"--foo"}, flag: "foo", hasFlag: true},
		{args: []string{"--foo=bar"}, flag: "foo", hasFlag: true},
		{args: []string{"--foo", "bar"}, flag: "foo", hasFlag: true},
		{args: []string{"-foo"}, flag: "foo", hasFlag: true},
		{args: []string{"-foo", "--bar"}, flag: "foo", hasFlag: true},
		{args: []string{"-foo", "--bar", "deneme"}, flag: "foo", hasFlag: true},
		{args: []string{"--bar", "bar", "-foo"}, flag: "foo", hasFlag: true},
		{args: []string{"--bar", "val", "-foo", "--bar"}, flag: "foo", hasFlag: true},
		{args: []string{"--foo=bar"}, flag: "foo", hasFlag: true},
		{args: []string{"--foo"}, flag: "bar", hasFlag: false},
		{args: []string{"--foo"}, flag: "", hasFlag: false},
	}

	for _, f := range flags {
		has := Has(f.flag, f.args)
		if has != f.hasFlag {
			t.Errorf("hasFlag: arg: %v flag: %v \n\twant: %v\n\tgot : %v\n", f.flag, f.args, f.hasFlag, has)
		}
	}
}

func TestIsFlag(t *testing.T) {
	var flags = []struct {
		flag   string
		isFlag bool
	}{
		{flag: "--foo", isFlag: true},
		{flag: "--foo=bar", isFlag: true},
		{flag: "-foo", isFlag: true},
		{flag: "-foo=bar", isFlag: true},
		{flag: "-f=bar", isFlag: true},
		{flag: "", isFlag: false},
	}

	for _, f := range flags {
		is := Valid(f.flag)
		if is != f.isFlag {
			t.Errorf("flag: %s\n\twant: %t\n\tgot : %t\n", f.flag, f.isFlag, is)
		}
	}
}

func TestParseFlag(t *testing.T) {
	var flags = []struct {
		flag string
		name string
	}{
		{name: "foo", flag: "--foo"},
		{name: "foo", flag: "-foo"},
		{name: "foo=bar", flag: "-foo=bar"},
		{name: "foo=", flag: "-foo="},
		{name: "foo=b", flag: "-foo=b"},
		{name: "", flag: "---f"},
		{name: "", flag: "f"},
		{name: "", flag: "--"},
		{name: "", flag: "-"},
	}

	for _, f := range flags {
		name, _ := Parse(f.flag)
		if name != f.name {
			t.Errorf("flag: %s\n\twant: %s\n\tgot : %s\n", f.flag, f.name, name)
		}
	}

}

func TestParseSingleFlagValue(t *testing.T) {
	var flags = []struct {
		flag  string
		name  string
		value string
	}{
		{flag: "foo=bar", name: "foo", value: "bar"},
		{flag: "foo=b", name: "foo", value: "b"},
		{flag: "", name: "", value: ""},
	}

	for _, f := range flags {
		name, value := parseSingleFlagValue(f.flag)
		if value != f.value {
			t.Errorf("parsing value from flag: %s\n\twant: %s\n\tgot : %s\n",
				f.flag, f.value, value)
		}

		if name != f.name {
			t.Errorf("parsing name from flag: %s\n\twant: %s\n\tgot : %s\n",
				f.flag, f.name, name)
		}
	}

}

func TestValueFromFlag(t *testing.T) {
	var arguments = []struct {
		args  []string
		value string
	}{
		{args: []string{"--provider=aws", "foo"}, value: "aws"},
		{args: []string{"-provider=aws", "foo", "bar"}, value: "aws"},
		{args: []string{"-provider=aws,do"}, value: "aws,do"},
		{args: []string{"--provider", "aws"}, value: "aws"},
		{args: []string{"-provider", "aws"}, value: "aws"},
		{args: []string{"-provider", "--foo"}, value: ""},
		{args: []string{"--provider", "--foo"}, value: ""},
		{args: []string{"--foo"}, value: ""},
		{args: []string{"--foo", "--provider"}, value: ""},
	}

	for _, args := range arguments {
		before := make([]string, len(args.args))
		copy(before, args.args)

		value, _ := Value("provider", args.args)

		if value != args.value {
			t.Errorf("parsing args value: %v\n\twant: %s\n\tgot : %s\n",
				args.args, args.value, value)
		}

		if !reflect.DeepEqual(before, args.args) {
			t.Errorf("parsing args modified the underlying slice for value: %v\n\twant: %s\n\tgot : %s\n",
				value, before, args.args)
		}
	}

}

func TestValueFromDashFlag(t *testing.T) {
	var arguments = []struct {
		args  []string
		value string
	}{
		{args: []string{"--access-key=aws", "foo"}, value: "aws"},
		{args: []string{"-access-key=aws", "foo", "bar"}, value: "aws"},
		{args: []string{"-access-key=aws,do"}, value: "aws,do"},
		{args: []string{"--access-key", "aws"}, value: "aws"},
		{args: []string{"-access-key", "aws"}, value: "aws"},
		{args: []string{"-access-key", "--foo"}, value: ""},
		{args: []string{"--access-key", "--foo"}, value: ""},
		{args: []string{"--asdasd"}, value: ""},
	}

	for _, args := range arguments {
		value, _ := Value("access-key", args.args)

		if value != args.value {
			t.Errorf("parsing dash args value: %v\n\twant: %s\n\tgot : %s\n",
				args.args, args.value, value)
		}

	}
}

func TestExcludeFlag(t *testing.T) {
	var arguments = []struct {
		args    []string
		remArgs []string
	}{
		{args: []string{}, remArgs: []string{}},
		{args: []string{"-provider"}, remArgs: []string{}},
		{args: []string{"--provider"}, remArgs: []string{}},
		{args: []string{"--provider=aws", "foo"}, remArgs: []string{"foo"}},
		{args: []string{"-provider=aws", "foo", "bar"}, remArgs: []string{"foo", "bar"}},
		{args: []string{"-provider=aws,do"}, remArgs: []string{}},
		{args: []string{"--test", "foo", "--provider=aws", "foo"}, remArgs: []string{"--test", "foo", "foo"}},
		{args: []string{"--example", "foo"}, remArgs: []string{"--example", "foo"}},
		{args: []string{"--test", "--provider", "aws"}, remArgs: []string{"--test"}},
		{args: []string{"--test", "--provider", "aws", "--test2"}, remArgs: []string{"--test", "--test2"}},
		{args: []string{"--test", "bar", "--provider", "aws"}, remArgs: []string{"--test", "bar"}},
		{args: []string{"--provider", "aws"}, remArgs: []string{}},
		{args: []string{"--provider", "aws", "--test"}, remArgs: []string{"--test"}},
		{args: []string{"--provider", "--test"}, remArgs: []string{"--test"}},
		{args: []string{"--test", "--provider"}, remArgs: []string{"--test"}},
		{args: []string{"--test", "bar", "--foo", "--provider"}, remArgs: []string{"--test", "bar", "--foo"}},
		{args: []string{"--test", "--provider", "--test2", "aws"}, remArgs: []string{"--test", "--test2", "aws"}},
	}

	for _, args := range arguments {
		before := make([]string, len(args.args))
		copy(before, args.args)

		remainingArgs := Exclude("provider", args.args)

		if !reflect.DeepEqual(remainingArgs, args.remArgs) {
			t.Errorf("parsing and returning rem args: %v\n\twant: %s\n\tgot : %s\n",
				args.args, args.remArgs, remainingArgs)
		}

		if !reflect.DeepEqual(before, args.args) {
			t.Errorf("parsing args modified the underlying slice: \n\twant: %s\n\tgot : %s\n",
				before, args.args)
		}
	}
}
