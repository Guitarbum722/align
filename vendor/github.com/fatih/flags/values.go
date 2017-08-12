package flags

import (
	"flag"
	"strconv"
	"strings"
)

type StringSliceValue []string

// NewStringSlice returns a new StringSlice which satisfies the flag.Value interface. This is
// useful to be used with flag.FlagSet. For the global flag instance,
// StringSlice() and StringSliceVar() can be used.
func NewStringSlice(val []string, p *[]string) *StringSliceValue {
	*p = val
	return (*StringSliceValue)(p)
}

func (s *StringSliceValue) Set(val string) error {
	// if empty default is used
	if val == "" {
		return nil
	}

	*s = StringSliceValue(strings.Split(val, ","))
	return nil
}

func (s *StringSliceValue) Get() interface{} { return []string(*s) }

func (s *StringSliceValue) String() string { return strings.Join(*s, ",") }

// StringSlice defines a []string flag with specified name, default value, and usage
// string. The return value is the address of a []string variable that stores
// the value of the flag.
func StringSlice(value []string, name, usage string) *[]string {
	p := new([]string)
	flag.Var(NewStringSlice(value, p), name, usage)
	return p
}

// StringSliceVar defines a []string flag with specified name, default value,
// and usage string. The argument p points to a []string variable in which to
// store the value of the flag.
func StringSliceVar(p *[]string, value []string, name, usage string) {
	flag.Var(NewStringSlice(value, p), name, usage)
}

type IntSliceValue []int

// NewIntSlice returns a new IntSlice which satisfies the flag.Value interface.
// This is useful to be used with flag.FlagSet. For the global flag instance,
// IntSlice() and IntSliceVar() can be used.
func NewIntSlice(val []int, p *[]int) *IntSliceValue {
	*p = val
	return (*IntSliceValue)(p)
}

func (i *IntSliceValue) Set(val string) error {
	// if empty default is used
	if val == "" {
		return nil
	}

	var list []int
	for _, in := range strings.Split(val, ",") {
		i, err := strconv.Atoi(in)
		if err != nil {
			return err
		}

		list = append(list, i)
	}

	*i = IntSliceValue(list)
	return nil
}

func (i *IntSliceValue) Get() interface{} { return []int(*i) }

func (i *IntSliceValue) String() string {
	var list []string
	for _, in := range *i {
		list = append(list, strconv.Itoa(in))
	}
	return strings.Join(list, ",")
}

// IntSlice defines a []int flag with specified name, default value, and usage
// string. The return value is the address of a []int variable that stores
// the value of the flag.
func IntSlice(value []int, name, usage string) *[]int {
	p := new([]int)
	flag.Var(NewIntSlice(value, p), name, usage)
	return p
}

// IntSliceVar defines a []int flag with specified name, default value, and
// usage string. The argument p points to a []int variable in which to store
// the value of the flag.
func IntSliceVar(p *[]int, value []int, name, usage string) {
	flag.Var(NewIntSlice(value, p), name, usage)
}
