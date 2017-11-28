package utility
import (
slogger "github.com/myafeier/logger"
"log"
"os"
)

var logger slogger.ILogger

func init() {
	logger = slogger.NewSimpleLogger2(os.Stdout, "[utility]", log.Lshortfile|log.Ldate|log.Lmicroseconds)
}