package main

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/integrii/flaggy"
	"github.com/samber/lo"
	"github.com/warl0ck/lazycontainer/pkg/app"
	"github.com/warl0ck/lazycontainer/pkg/config"
)

const DEFAULT_VERSION = "unversioned"

var (
	commit      string
	version     = DEFAULT_VERSION
	date        string
	buildSource = "unknown"

	debuggingFlag = false
)

func main() {
	updateBuildInfo()

	info := fmt.Sprintf(
		"%s\nDate: %s\nBuildSource: %s\nCommit: %s\nOS: %s\nArch: %s",
		version,
		date,
		buildSource,
		commit,
		runtime.GOOS,
		runtime.GOARCH,
	)

	flaggy.SetName("lazycontainer")
	flaggy.SetDescription("The lazier way to manage Apple containers")
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/warl0ck/lazycontainer"

	flaggy.Bool(&debuggingFlag, "d", "debug", "Enable debug mode")
	flaggy.SetVersion(info)

	flaggy.Parse()

	appConfig, err := config.NewAppConfig("lazycontainer", version, commit, date, buildSource, debuggingFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	application, err := app.NewApp(appConfig)
	if err == nil {
		err = application.Run()
	}

	if application != nil {
		application.Close()
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func updateBuildInfo() {
	if version == DEFAULT_VERSION {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			revision, ok := lo.Find(buildInfo.Settings, func(setting debug.BuildSetting) bool {
				return setting.Key == "vcs.revision"
			})
			if ok {
				commit = revision.Value
				version = safeTruncate(revision.Value, 7)
			}

			time, ok := lo.Find(buildInfo.Settings, func(setting debug.BuildSetting) bool {
				return setting.Key == "vcs.time"
			})
			if ok {
				date = time.Value
			}
		}
	}
}

func safeTruncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
