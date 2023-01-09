# Status code convention

## exeIaC status code

You can read the status code to understand why exeiac failed.

- 0: when everything is ok
- 3: only for plan action: when the module can't decide if brick has 
  drift or not. It's a shortcut when it's too difficult to get that information.
- 2: only for plan action: when a drift is detected
- 11: an error have been encounter during the exeiac initialisation
  (parsing error, bad arg, bad conf file...)
- 12: an error have been encounter during the "enrichment" it's usually due to
  the brick.yml, or the non ability to retrieve outut of dependencies bricks
- 13: an error have been encounter during the run of the action flow (but not
  during a module execution). This error type haven't stopped the action flow
- 4: when we encounter an error during the module execution
- 14: an error have been encounter during the the action flow (but not
  during a module execution) and stopped it.
- 254: an internal error have been encounter when choosing a status code
- 1: we try to avoid this status code, if you encounter it it usually means that 
  exeiac has paniced

**Note:** If you run  an action on multiple elementary brick you can encounter
multiple errors but obviously you still have only one status code. In such a case
you will have the last status code on this list

## module status code convention

You can use the entire exeiac convention for your module if you want. But it can
be smarter to simply forward status code of the tool used by your module.

But at least you have to follow that convention for all action except plan:
- 0: when action succeed
- 1-255: when action failed (except for plan)

For plan:
- 0: there is no drift
- 2: there is a drift
- 3: module can't decide if there is a drift or not (but it has displayed the plan)
- 1,4-255: plan has failed

## linked package

src/statuscode
