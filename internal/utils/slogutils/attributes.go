package slogutils

import (
	"log/slog"

	"github.com/pkg/errors"
)

func ErrorAttr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ErrorAttrWrap(err error, message string) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(errors.Wrap(err, message).Error()),
	}
}

func PanicAttr(v any) slog.Attr {
	return slog.Attr{
		Key:   "panic",
		Value: slog.AnyValue(v),
	}
}
