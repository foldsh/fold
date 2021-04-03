package handler

import (
	"net/http"
	"os"

	"github.com/foldsh/fold/logging"
)

func NewHTTP(logger logging.Logger, server http.Handler, port string) *HTTPHandler {
	return &HTTPHandler{logger, server, port}
}

type HTTPHandler struct {
	logger logging.Logger
	server http.Handler
	port   string
}

func (hh *HTTPHandler) Serve() {
	hh.logger.Debugf("Listening on port %s", hh.port)
	if err := http.ListenAndServe(hh.port, hh.server); err != nil {
		hh.logger.Errorf(err.Error())
		os.Exit(1)
	}
}

// func (r *Runtime) setupSignalHandler() error {
// 	r.logger.Debugf("Registering signal handler")
// 	if err := signal.Notify(r.signals, syscall.SIGINT, syscall.SIGTERM); err != nil {
// 		return err
// 	}
// 	go func() {
// 		signal := <-r.signals
// 		// We try to exit according to the signal, but if either of these goes wrong we will
// 		// just exit with a non 0 status.
// 		if err := r.supervisor.Signal(signal); err != nil {
// 			os.Exit(1)
// 		}
// 		if err := s.Wait(); err != nil {
// 			os.Exit(1)
// 		}
// 		// Great, lets emit a STOP so that the rest of the runtime can do whatever clean up
// 		// it needs to.
// 		r.Emit(STOP)
// 	}()
// 	return nil
// }
