package statuscode

// NOTE(arthur91f): should be coherent with docs/statusCodeConvention.md

const MODULE_OK = 0            // it's ok
const MODULE_DONE = 0          // it has been runned well
const MODULE_NOTHING_TO_DO = 0 // it's ok because there was nothing to do
const MODULE_NO_DRIFT = 0      // the plan return that nothing has drift
const MODULE_DRIFT = 2         // the plan have detected a drift
const MODULE_DRIFT_OR_NOT = 3  // the plan don't know if there is a drift or not
const MODULE_ERROR = 4         // the module have fail
const INIT_ERROR = 11          // an error have been encountered during the infra initialisation
const ENRICH_ERROR = 12        // an error have been encountered during the enrichment
const RUN_ERROR = 13           // an error have been encountered during the action flow (not during the module execution)
const FATAL_ERROR = 14         // an error have been encountered during the action flow and interrupt the action flow
const STATUS_CODE_ERROR = 254  // when update failed

func Update(current int, new int) int {
	if current == FATAL_ERROR || new == FATAL_ERROR { // it's the worst error: something went wrong and we don't know what
		return FATAL_ERROR
	}
	if current == MODULE_ERROR || new == MODULE_ERROR { // it's the worst error after FATAL because the module can have fail in the middle of an apply
		return MODULE_ERROR
	}
	if current == 1 || new == 1 { // we don't want to return 1 or it means that we have panic
		return STATUS_CODE_ERROR
	}
	if (current == MODULE_DRIFT && new == MODULE_DRIFT_OR_NOT) ||
		(new == MODULE_DRIFT && current == MODULE_DRIFT_OR_NOT) { // if we know that one brick have drifted so the brick groups have drifted
		return MODULE_DRIFT
	}

	// For the rest ???>RUN>ENRICH>INIT>DRIFT_OR_NOT>DRIFT>OK
	if new > current {
		return new
	}
	return current
}
