package logutil

import (
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log := NewLogger(DEBUG, "", os.Stdout)
	log.Debug("debug:%s", time.Now().String())
	log.Info("info:%s", time.Now().String())
	log.Notice("notice:%s", time.Now().String())
	log.Warning("warning:%s", time.Now().String())
	log.Error("error:%s", time.Now().String())
	log.Critical("critical:%s", time.Now().String())
}
