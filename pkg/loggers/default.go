package loggers

import (
	"fmt"
	"time"
)

type Default struct {
	name string
}

func NewDefault(name string) *Default {
	return &Default{
		name: name,
	}
}

func (l *Default) prefix(p string, s string) string {
	return fmt.Sprintf("%s[%s] [%s]%s", p, time.Now().Format(time.DateTime), l.name, s)
}

func (l *Default) Infof(format string, args ...any) {
	fmt.Printf(l.prefix(" [INFO] ", " ")+format+"\n", args...)
}

func (l *Default) Errorf(format string, args ...any) {
	fmt.Printf(l.prefix("[ERROR] ", " ")+format+"\n", args...)
}

func (l *Default) Debugf(format string, args ...any) {
	fmt.Printf(l.prefix("[DEBUG] ", " ")+format+"\n", args...)
}

func (l *Default) Warnf(format string, args ...any) {
	fmt.Printf(l.prefix(" [WARN] ", " ")+format+"\n", args...)
}

func (l *Default) Info(args ...any) {
	fmt.Println(l.prefix(" [INFO] ", " "), fmt.Sprint(args...))
}

func (l *Default) Error(args ...any) {
	fmt.Println(l.prefix("[ERROR] ", " "), fmt.Sprint(args...))
}

func (l *Default) Debug(args ...any) {
	fmt.Println(l.prefix("[DEBUG] ", " "), fmt.Sprint(args...))
}

func (l *Default) Warn(args ...any) {
	fmt.Println(l.prefix(" [WARN] ", " "), fmt.Sprint(args...))
}
