package watcher

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/foldsh/fold/logging"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	dir      string
	logger   logging.Logger
	watcher  *fsnotify.Watcher
	stop     chan struct{}
	onChange func()
}

func NewWatcher(logger logging.Logger, dir string, onChange func()) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.New("failed to initialise watcher")
	}
	return &Watcher{dir, logger, watcher, make(chan struct{}), onChange}, nil
}

func (w *Watcher) Close() {
	w.stop <- struct{}{}
	w.watcher.Close()
}

func (w *Watcher) Watch() error {
	if err := filepath.Walk(w.dir, w.watchDir); err != nil {
		w.logger.Debugf("Error walking directory marked for watching %v", err)
		return err
	}
	go func() {
		for {
			select {
			case <-w.watcher.Events:
				w.onChange()
			case err := <-w.watcher.Errors:
				w.logger.Debugf("Watcher encountered error %v", err)
			case <-w.stop:
				return
			}
		}
	}()
	return nil
}

func (w *Watcher) watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return w.watcher.Add(path)
	}
	return nil
}
