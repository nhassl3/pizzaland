package sl

import (
	"fmt"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ErrUpLevel(handleName, err string) error {
	return fmt.Errorf(handleName + ": " + err)
}
