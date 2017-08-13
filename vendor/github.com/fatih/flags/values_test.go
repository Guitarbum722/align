package flags

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestStringList(t *testing.T) {
	regions := StringSlice(nil, "to", "Regions to be used")

	os.Args = []string{"cmd", "-to", "us-east-1,eu-west-2"}
	flag.Parse()

	want := []string{"us-east-1", "eu-west-2"}
	if !reflect.DeepEqual(*regions, want) {
		t.Errorf("Regions = %q, want %q", regions, want)
	}
}

func TestStringListVar(t *testing.T) {
	var regions []string
	StringSliceVar(&regions, nil, "tos", "Regions to be used")

	os.Args = []string{"cmd", "-tos", "us-east-1,eu-west-2"}
	flag.Parse()

	want := []string{"us-east-1", "eu-west-2"}
	if !reflect.DeepEqual(regions, want) {
		t.Errorf("Regions = %q, want %q", regions, want)
	}
}

func TestStringListFlagSet(t *testing.T) {
	f := flag.NewFlagSet("TestTags", flag.PanicOnError)

	var regions []string
	f.Var(NewStringSlice(nil, &regions), "to", "Regions to be used")
	f.Parse([]string{"-to", "us-east-1,eu-west-2"})

	want := []string{"us-east-1", "eu-west-2"}
	if !reflect.DeepEqual(regions, want) {
		t.Errorf("Regions = %q, want %q", regions, want)
	}
}

func TestIntList(t *testing.T) {
	ids := IntSlice(nil, "ids", "Servers to be used")

	os.Args = []string{"cmd", "-ids", "123,456"}
	flag.Parse()

	want := []int{123, 456}
	if !reflect.DeepEqual(*ids, want) {
		t.Errorf("Ids = %q, want %q", *ids, want)
	}
}

func TestIntListVar(t *testing.T) {
	var ids []int
	IntSliceVar(&ids, nil, "idss", "Servers to be used")

	os.Args = []string{"cmd", "-idss", "123,456"}
	flag.Parse()

	want := []int{123, 456}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("Ids = %q, want %q", ids, want)
	}
}

func TestIntListFlagSet(t *testing.T) {
	f := flag.NewFlagSet("TestTags", flag.PanicOnError)

	var ids []int
	f.Var(NewIntSlice(nil, &ids), "ids", "Ids to be used")
	f.Parse([]string{"-ids", "123,456"})

	want := []int{123, 456}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("Ids = %q, want %q", ids, want)
	}
}
