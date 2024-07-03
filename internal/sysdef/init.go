package sysdef

// -----------------------------------------------------------------------------------
var (
	AppName      = "VistaDataEngine"
	AppVersion   = "1.0"
	ConfFilename = "setting.yml"
	Log          tLog

	DbMaxOpenCnt = -1
	DbMaxIdleCnt = -1
)

type tLog struct {
	Path            string
	Name            string
	NameCollector   string
	NameVmCollector string
	MaxAge          int
	MaxSizeMB       int
	MaxBackup       int
}

// -----------------------------------------------------------------------------------

func init() {
	Log.Path = "./"
	Log.Name = "Main.log"
	Log.Name = "Main.log"
	Log.NameCollector = "Collector.log"
	Log.NameVmCollector = "CollectorVm.log"
	Log.MaxAge = 31
	Log.MaxSizeMB = 10
	Log.MaxBackup = 10
}
