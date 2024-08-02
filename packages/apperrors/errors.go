package apperrors

import (
	"errors"
	"fmt"
	"log/slog"

	cerrors "github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
)

var (
	ErrNotFound = errors.New("item not found")
)

func WithStack(err error) error {
	if err == nil {
		return cerrors.WithStack(err)
	}

	_, ok := err.(interface{ SafeFormatError(errbase.Printer) error })
	if ok {
		return err
	}

	return cerrors.WithStack(err)
}

func LogStackTrace(err error) slog.Attr {
	if err == nil {
		return slog.Any("stacktrace", []any{})
	}

	_, ok := err.(interface{ SafeFormatError(errbase.Printer) error })
	if !ok {
		return slog.String("details", fmt.Sprintf("%+v", err))

	}

	return slog.Any("stacktrace", cerrors.GetReportableStackTrace(err).Frames)
}
