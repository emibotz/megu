package loggers

import "fmt"

type Test struct {
}

func (l *Test) Infof(format string, args ...interface{}) {
	fmt.Printf("INFO: "+format+"\n", args...)
}

func (l *Test) Errorf(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}

func (l *Test) Debugf(format string, args ...interface{}) {
	fmt.Printf("DEBUG: "+format+"\n", args...)
}

func (l *Test) Warnf(format string, args ...interface{}) {
	fmt.Printf("WARN: "+format+"\n", args...)
}

func (l *Test) Info(args ...interface{}) {
	fmt.Println("INFO:", fmt.Sprint(args...))
}

func (l *Test) Error(args ...interface{}) {
	fmt.Println("ERROR:", fmt.Sprint(args...))
}

func (l *Test) Debug(args ...interface{}) {
	fmt.Println("DEBUG:", fmt.Sprint(args...))
}

func (l *Test) Warn(args ...interface{}) {
	fmt.Println("WARN:", fmt.Sprint(args...))
}
