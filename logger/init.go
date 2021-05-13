package logger

import "path/filepath"

// Logger logger single
var Logger Helper

// Init logger
func Init(options *Options) (e error) {
	if options.Filename == "" {
		options.Filename = filepath.Join("logs", "master.log")
	}
	Logger.Attach(NewHelper(options))
	return
}
