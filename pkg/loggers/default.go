package loggers

import (
	"fmt"
	"strings"
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
	fmt.Printf(l.prefix(" [INFO] ", " ")+strings.Trim(format, " \n")+"\n", args...)
}

func (l *Default) Errorf(format string, args ...any) {
	fmt.Printf(l.prefix("[ERROR] ", " ")+strings.Trim(format, " \n")+"\n", args...)
}

func (l *Default) Debugf(format string, args ...any) {
	fmt.Printf(l.prefix("[DEBUG] ", " ")+strings.Trim(format, " \n")+"\n", args...)
}

func (l *Default) Warnf(format string, args ...any) {
	fmt.Printf(l.prefix(" [WARN] ", " ")+strings.Trim(format, " \n")+"\n", args...)
}

func (l *Default) Info(args ...any) {
	fmt.Println(l.prefix(" [INFO] ", ""), strings.Trim(fmt.Sprint(args...), " \r\n"))
}

func (l *Default) Error(args ...any) {
	fmt.Println(l.prefix("[ERROR] ", ""), strings.Trim(fmt.Sprint(args...), " \r\n"))
}

func (l *Default) Debug(args ...any) {
	fmt.Println(l.prefix("[DEBUG] ", ""), strings.Trim(fmt.Sprint(args...), " \r\n"))
}

func (l *Default) Warn(args ...any) {
	fmt.Println(l.prefix(" [WARN] ", ""), strings.Trim(fmt.Sprint(args...), " \r\n"))
}
