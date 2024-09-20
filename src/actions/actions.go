package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
	extools "src/exeiac/tools"
)

struct ModulesCommand {
	Priority int
	Name string
	SubCommand    []string
}

const ModulesCommands = []ModulesCommand{
	ModulesCommand{
		Priority: 0,
		Name: "non-mutable"
	},
	ModulesCommand{
		Priority: ,
		Name: "output"
	},
	ModulesCommand{
		Priority: ,
		Name: "plan"
	},
	ModulesCommand{
		Priority: ,
		Name: "remove"
	},
	ModulesCommand{
		Priority: ,
		Name: "partial-lay"
	},
	ModulesCommand{
		Priority: ,
		Name: "lay"
		SubCommand: []{"partial-lay"}
	},
	ModulesCommand{
		Priority: ,
		Name: "post-lay"
	},
}

// actionOrder := []string{
// 	0, "non-mutable", // don't change infra and exeiac's system
// 	1, "output"
// 	101, "system-mutable",
// 	102, "iac-mutable",
// 	103, "non-infra-mutable", // don't change infra but may change the system that runs exeiac as install-python-deps
// 	104, "pre-lay",
// 	299, "plan",
// 	399, "remove",
// 	428, "partial-create",
// 	429, "create",
// 	458, "partial-update",
// 	459, "update-after-create",
// 	498, "destroy-recreate",
// 	499, "full-lay",
// 	599, "post-lay",
// }
