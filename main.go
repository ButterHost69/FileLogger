package filelogger

import (
	"log"
	"os"
	"time"
)

type FileLogger struct {
	LogFilePath string
	DateAndTime bool
}

type LoggerOptions func(fp *FileLogger)

func AddDateAndTime() LoggerOptions {
	return func(fp *FileLogger) {
		fp.DateAndTime = true
	}
}

func (fl *FileLogger) Init(filePath string, opts ...LoggerOptions) error {
	fl.LogFilePath = filePath

	_, err := os.Stat(filePath)
	if err != nil {
		if !os.IsExist(err) {
			log.Println("file does not exists...\nCreating the file ", filePath)
			os.Create(filePath)
		} else {
			log.Println(err)
		}
		return err
	}
	for _, opt := range opts {
		opt(fl)
	}

	return nil
}

func (fl *FileLogger) writeToFile(data string) error {
	return os.WriteFile(fl.LogFilePath, []byte(data), 0777)
}

func (fl *FileLogger) Log(param interface{}) {
	var data string
	switch inp := param.(type) {
	case error:
		if fl.DateAndTime {
			data = time.Now().String() + " " + inp.Error()
		} else {
			data = inp.Error()
		}
		fl.writeToFile(data)
	case string:
		if fl.DateAndTime {
			data = time.Now().String() + " " + inp
		} else {
			data = inp
		}
		err := fl.writeToFile(data)
		if err != nil {
			log.Println(err.Error())
		}

	}
}
