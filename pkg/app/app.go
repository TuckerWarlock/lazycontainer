package app

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/config"
	"github.com/warl0ck/lazycontainer/pkg/gui"
)

// App represents the main application
type App struct {
	closers          []io.Closer
	Config           *config.AppConfig
	Log              *logrus.Entry
	OSCommand        *commands.OSCommand
	ContainerCommand *commands.ContainerCommand
	Gui              *gui.Gui
}

// NewApp creates a new App instance
func NewApp(appConfig *config.AppConfig) (*App, error) {
	app := &App{
		Config:  appConfig,
		closers: []io.Closer{},
	}

	var err error

	// Setup logging
	app.Log = newLogger(appConfig)

	// Setup OS command executor
	app.OSCommand = commands.NewOSCommand(app.Log)

	// Setup Container command wrapper
	app.ContainerCommand = commands.NewContainerCommand(app.Log, app.OSCommand, appConfig)

	// Setup GUI
	app.Gui, err = gui.NewGui(app.Log, app.ContainerCommand, app.OSCommand, appConfig)
	if err != nil {
		return app, err
	}

	return app, nil
}

// Run starts the application
func (app *App) Run() error {
	return app.Gui.Run()
}

// Close cleans up application resources
func (app *App) Close() {
	for _, closer := range app.closers {
		_ = closer.Close()
	}
}

func newLogger(config *config.AppConfig) *logrus.Entry {
	log := logrus.New()

	if config.Debug {
		log.SetLevel(logrus.DebugLevel)
		logPath := config.GetLogPath()
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		}
	} else {
		log.SetLevel(logrus.ErrorLevel)
		log.SetOutput(io.Discard)
	}

	log.SetFormatter(&logrus.JSONFormatter{})

	return log.WithFields(logrus.Fields{
		"app": config.Name,
	})
}
