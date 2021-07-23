package conflogger

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-courier/metax"
	"github.com/sirupsen/logrus"
)

type Log struct {
	ReportCaller bool
	Name         string
	Level        string `env:""`
	Format       string
	init         bool
}

func (log *Log) SetDefaults() {
	if log.Name == "" {
		log.Name = os.Getenv("PROJECT_NAME")

		version := os.Getenv("PROJECT_VERSION")
		if version != "" {
			log.Name += "@" + version
		}
	}

	if os.Getenv("GOENV") == "DEV" {
		log.ReportCaller = true
	}

	if log.Level == "" {
		log.Level = "DEBUG"
	}

	if log.Format == "" {
		log.Format = "json"
	}
}

func (log *Log) Init() {
	if !log.init {
		log.Create()
		log.init = true
	}
}

func (log *Log) Create() {

	if log.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			//PrettyPrint:      true,
			CallerPrettyfier: CallerPrettyfier,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			CallerPrettyfier: CallerPrettyfier,
		})
	}

	logrus.SetLevel(getLogLevel(log.Level))
	logrus.SetReportCaller(log.ReportCaller)
	logrus.AddHook(NewProjectAndMetaHook(log.Name))

	logrus.SetOutput(os.Stdout)
}

func getLogLevel(l string) logrus.Level {
	level, err := logrus.ParseLevel(strings.ToLower(l))
	if err == nil {
		return level
	}
	return logrus.InfoLevel
}

func NewProjectAndMetaHook(name string) *ProjectAndMetaHook {
	return &ProjectAndMetaHook{
		Name: name,
	}
}

type ProjectAndMetaHook struct {
	Name string
}

func (hook *ProjectAndMetaHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}

	meta := metax.MetaFromContext(ctx)
	entry.Data["project"] = hook.Name
	for k, v := range meta {
		entry.Data["meta."+k] = v
	}

	return nil
}

func (hook *ProjectAndMetaHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func CallerPrettyfier(f *runtime.Frame) (function string, file string) {
	return f.Function + " line:" + strconv.FormatInt(int64(f.Line), 10), ""
}
