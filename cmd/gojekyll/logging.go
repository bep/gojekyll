package main

import (
	"fmt"

	"github.com/osteele/gojekyll/helpers"
)

type bannerLogger struct{ labelWidth int }

var logger = bannerLogger{}

func (l *bannerLogger) Info(a ...interface{}) {
	fmt.Println(a...)
}

func (l *bannerLogger) label(label string, msg string, a ...interface{}) {
	if len(label) > l.labelWidth {
		l.labelWidth = len(label)
	}
	if !quiet {
		fmt.Printf("%s %s\n", helpers.LeftPad(label, l.labelWidth), fmt.Sprintf(msg, a...))
	}
}

func (l *bannerLogger) path(label string, filename string) {
	l.label(label, helpers.MustAbs(filename))
}
