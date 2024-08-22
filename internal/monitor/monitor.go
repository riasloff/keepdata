package monitor

import (
	"fmt"
	"keepdata/internal/logger"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func Monitor(files, dirs []string) error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	errCh := make(chan error)

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					logger.Log.Warn("выполнена запись файла", zap.String("user", "user"), zap.String("action", fmt.Sprint(event.Op)), zap.String("resource", event.Name))
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					logger.Log.Warn("выполнено изменение прав файла", zap.String("user", "user"), zap.String("action", fmt.Sprint(event.Op)), zap.String("resource", event.Name))
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					logger.Log.Warn("выполнено переименование файла", zap.String("user", "user"), zap.String("action", fmt.Sprint(event.Op)), zap.String("resource", event.Name))
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					logger.Log.Warn("выполнено удаление прав файла", zap.String("user", "user"), zap.String("action", fmt.Sprint(event.Op)), zap.String("resource", event.Name))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				errCh <- err
			}
		}
	}()

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			return err
		}
		logger.Log.Info("начали новое слежение", zap.String("file", file))
	}

	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			return err
		}
		logger.Log.Info("начали новое слежение", zap.String("dir", dir))
	}

	<-make(chan struct{})
	return nil
}
