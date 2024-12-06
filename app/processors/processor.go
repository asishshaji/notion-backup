package processors

import "github.com/asishshaji/notion-backup/app/actions"

type Processor interface {
	Process() error            // function that triggers the actions
	Actions() []actions.Action // returns the sequence of actions used by the processor
	initialiseActions()        // sets the actions for the process
}
