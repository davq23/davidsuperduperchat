package utils

import "log"

// Logger manage concurrent logs
type Logger struct {
	LogChan chan interface{}
}

// Logs go function for logging
func (l *Logger) Logs() {
	for {
		select {
		case message := <-l.LogChan:
			log.Printf("%v", message)
		}
	}
}
