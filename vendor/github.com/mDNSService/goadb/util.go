package adb

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mDNSService/goadb/internal/errors"
)

var (
	whitespaceRegex = regexp.MustCompile(`^\s*$`)
)

func containsWhitespace(str string) bool {
	return strings.ContainsAny(str, " \t\v")
}

func isBlank(str string) bool {
	return whitespaceRegex.MatchString(str)
}

func wrapClientError(err error, client interface{}, operation string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*errors.Err); !ok {
		panic("err is not a *Err: " + err.Error())
	}

	clientType := reflect.TypeOf(client)

	return &errors.Err{
		Code:    err.(*errors.Err).Code,
		Cause:   err,
		Message: fmt.Sprintf("error performing %s on %s", fmt.Sprintf(operation, args...), clientType),
		Details: client,
	}
}
