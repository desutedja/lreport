package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)

	// Set calling method as a field
	logrus.SetReportCaller(true)
}
