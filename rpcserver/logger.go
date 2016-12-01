package rpcserver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func createOrOpenFile(filename string) *os.File {
	var f *os.File
	f, _ = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	return f
}

func NewLogger(dir, fileName, perfix string) (*log.Logger, error) {
	if len(perfix) > 0 {
		perfix = fmt.Sprintf("[%s] ", perfix)
	}
	f := filepath.Join(dir, fileName)
	myLog := log.New(createOrOpenFile(f), perfix, log.Lmicroseconds|log.Lshortfile)
	return myLog, nil
}
