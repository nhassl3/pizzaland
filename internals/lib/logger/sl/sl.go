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

func ErrUpLevel(handleName string, err error) error {
	return fmt.Errorf("%s: %w", handleName, err)
}
