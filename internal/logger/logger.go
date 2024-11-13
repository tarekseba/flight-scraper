package logger

import (
	"log"
	"os"
)

const flags = log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile

var InfoLogger = log.New(os.Stdout, "\033[32mINFO  \033[0m", flags)
var ErrorLogger = log.New(os.Stdout, "\033[31mERROR \033[0m", flags)
var WarnLogger = log.New(os.Stdout, "\033[33mWARN  \033[0m", flags)
