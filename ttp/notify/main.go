package notify

import (
	"fmt"

	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A new server has been added at "

func NotifyViolation[T any](serv *ent.TAServer, ttp *core.TTP[T]) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, serv.Domain)
	return ttp.Notify.Notify([]byte(msg), serv, ttp)
}

func NotifyUpdate[T any](serv *ent.TAServer, ttp *core.TTP[T]) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, serv.Domain)
	return ttp.Notify.Notify([]byte(msg), serv, ttp)
}
