package utilities

import (
	"github.com/sirupsen/logrus"
	"io"
)

func SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		logrus.WithError(err).Warn("Close operation failed")
	}
}
