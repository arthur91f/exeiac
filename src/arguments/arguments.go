package arguments

import (
	"fmt"
	"reflect"
	"strings"
)

type Arguments struct {
	Action            string
	BricksNames       []string
	BricksSpecifiers  []string
	NonInteractive    bool
	Format            string
	Modules           map[string]string
	OtherOptions      []string
	Rooms             map[string]string
	ConfigurationFile string
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

var actions_list = []string{
	"plan", "lay", "remove", "output", "init", "validate_code", "help",
	"show_input", "list_elementary_bricks", "cd",
	"get_brick_path", "get_brick_name"}

var AvailableBricksSpecifiers = []string{
	"linked_previous", "all_previous", "lp", "ap",
	"direct_previous", "dp",
	"selected", "s",
	"direct_next", "dn",
	"linked_next", "all_next", "ln", "an"}
