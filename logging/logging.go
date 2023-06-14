package logging

import "log"

type Logfile struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
}

func Check(erro error, loggin Logfile) {
	if erro != nil {
		loggin.ErrorLogger.Fatal(erro)
	}
}
