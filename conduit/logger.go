package conduit

import (
	"fmt"
	"log"
	"os"
)

var Logger = NewLogger("conduit-app")

func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags|log.Lmsgprefix)
}
