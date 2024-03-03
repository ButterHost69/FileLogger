package filelogger

import (
	"errors"
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

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// log.Printf("file does not exists...\nCreating the file %s\n", filePath)
		file, err := os.Create(filePath)
		if err != nil {
			log.Panic(err)
			return err
		}
		file.Close()
	} else if err != nil {
		log.Println(err)
		return err
	}
	for _, opt := range opts {
		opt(fl)
	}

	return nil
}

func (fl *FileLogger) writeToFile(data string) error {
	file, err := os.OpenFile(fl.LogFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data + "\n")
	return err
}

func (fl *FileLogger) Log(param interface{}) error {
	var data string
	switch inp := param.(type) {
	case error:
		if fl.DateAndTime {
			data = time.Now().Format(time.RFC3339) + " " + inp.Error()
		} else {
			data = inp.Error()
		}
	case string:
		if fl.DateAndTime {
			data = time.Now().Format(time.RFC3339) + " " + inp
		} else {
			data = inp
		}
	default:
		data = "Unsupported type"
	}

	err := fl.writeToFile(data)
	if err != nil {
		return errors.New("Error writing to log file:" + err.Error())
	}

	return nil
}