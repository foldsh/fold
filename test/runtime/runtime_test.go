package test

import (
	"os"
	"testing"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
	"github.com/foldsh/fold/runtime/supervisor"
)

func TestBasicService(t *testing.T) {
	logger := logging.NewTestLogger()
	cmd := "go"
	args := []string{"run", "./testdata/basic/"}
	done := make(chan struct{})
	// sout := &bytes.Buffer{}
	// sett := &bytes.Buffer{}
	sout := os.Stdout
	serr := os.Stdout
	rt := runtime.NewRuntime(
		logger,
		cmd,
		args,
		done,
		runtime.WithSupervisor(supervisor.NewSupervisor(logger, cmd, args, sout, serr)),
	)

	rt.Start()
	// if err != nil {
	// 	t.Fatalf("%+v", err)
	// }

	<-done

}
