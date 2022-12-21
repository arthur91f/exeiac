package arguments

import (
	"fmt"
	"reflect"
	"strings"
)

// A struct matching the arguments available to be provided through the command line.
// Used only to gather the command line arguments once parsed, and to be used as
// a parameter to build a `Configuration` struct.
type Arguments struct {
	Action            string
	BricksNames       []string
	BricksSpecifiers  []string
	NonInteractive    bool
	Interactive       bool
	Format            string
	Modules           map[string]string
	OtherOptions      []string
	Rooms             map[string]string
	ConfigurationFile string
	ShowUsage         bool
	ListBricks        bool
}

func (a Arguments) String() string {
	var sb strings.Builder

	v := reflect.ValueOf(a)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		sb.WriteString("\t")
		sb.WriteString(t.Field(i).Name)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%v", v.Field(i).Interface()))
		sb.WriteString("\n")
	}

	return sb.String()
}

// An array containing all of the supported actions
var actions_list = [...]string{
	"plan", "lay", "remove", "output", "init", "validate_code", "help",
	"show_input", "list_elementary_bricks", "cd",
	"get_brick_path", "get_brick_name"}

// An array containing all of the supported brick's specifiers
var AvailableBricksSpecifiers = [...]string{
	"linked_previous", "all_previous", "lp", "ap",
	"direct_previous", "dp",
	"selected", "s",
	"direct_next", "dn",
	"linked_next", "all_next", "ln", "an"}

var AvailableBricksFormat = [...]string{
	"name", "n",
	"path", "p",
	"all", "a"}
