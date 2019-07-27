package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

func (err MyError) Error() string {
	return err.Message
}

//lowlevel
type LowlevelErr struct {
	error
}

func isGloballExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowlevelErr{wrapError(err, err.Error())}
	}
	return info.Mode().Perm()&100 == 0100, nil
}

//intermediate
type InterMediateErr struct {
	error
}

func runJob(id string) error {
	const jobBibPath = "/bad/job/binary"
	isExectable, err := isGloballExec(jobBibPath)
	if err != nil {
		return InterMediateErr{wrapError(
			err,
			"cannot runjob %q: requisite binaries not available",
			id,
		)}
	} else if isExectable == false {
		return wrapError(
			nil,
			"cannot runjob %q: requisite binaries are not executable",
		)
	}
	return exec.Command(jobBibPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: \n", key))
	log.Printf("%#v\n", err)
	fmt.Printf("[%v] %v\n", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(InterMediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
