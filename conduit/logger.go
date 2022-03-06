package conduit

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ashalfarhan/realworld/config"
)

var Logger = log.New(os.Stdout, "[ConduitApp] ", log.LstdFlags|log.Lmsgprefix)

func NewLogger(prefix string) *log.Logger {
	var out io.Writer = os.Stdout
	if config.Co.Env == "test" {
		Logger.SetOutput(io.Discard)
		out = io.Discard
	}
	return log.New(out, fmt.Sprintf("[%s] ", prefix), log.LstdFlags|log.Lmsgprefix|log.Llongfile)
}
