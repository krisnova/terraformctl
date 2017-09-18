package tfmain

import (
	"log"
	"os"
)

func HackedMain(args []string) int {
	// Override global prefix set by go-dynect during init()
	log.SetPrefix("")

	//
	//
	//
	//
	os.Args = args
	//
	//
	//
	//

	return realMain()
}
