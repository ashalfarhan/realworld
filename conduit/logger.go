package conduit

import (
	"fmt"
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[conduit-app] ", log.LstdFlags|log.Lmsgprefix)

func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
}
