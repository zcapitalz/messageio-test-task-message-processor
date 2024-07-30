package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"message-processor/internal/config"
	"os"

	"github.com/pkg/errors"
)

func newTLSConfig(cfg config.KafkaSSLConfig) (*tls.Config, error) {
	caCert, err := os.ReadFile(cfg.CaCertFile)
	if err != nil {
		return nil, errors.Wrap(err, "read CA certificate file")
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("could not parse CA certificate")
	}

	return &tls.Config{
		RootCAs: caCertPool,
	}, nil
}

type logger struct {
	prefix string
}

func newLogger(prefix string) *logger {
	return &logger{prefix: prefix}
}

func (l *logger) Printf(msg string, args ...any) {
	slog.Info(fmt.Sprintf(l.prefix+msg, args...))
}

type errorLogger struct {
	prefix string
}

func newErrorLogger(prefix string) *errorLogger {
	return &errorLogger{prefix: prefix}
}

func (l *errorLogger) Printf(msg string, args ...any) {
	slog.Error(fmt.Sprintf(l.prefix+msg, args...))
}
