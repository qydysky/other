package v2rayf

//go:generate errorgen

import (
	"flag"
	// "fmt"
	"io/ioutil"
	// "log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	// "strconv"

	// "./part"
	// "../session"
	"v2ray.com/core"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/errors"
	_ "v2ray.com/core/main/distro/all"

	"github.com/qydysky/part"
)

type errPathObjHolder struct{}

var (
	configFiles cmdarg.Arg // "Config file for V2Ray.", the option is customed type, parse in main
	configDir   string
	version     = flag.Bool("version", false, "Show current version of V2Ray.")
	test        = flag.Bool("test", false, "Test config file only, without launching V2Ray server.")
	format      = flag.String("format", "json", "Format of input file.")

	/* We have to do this here because Golang's Test will also need to parse flag, before
	 * main func in this file is run.
	 */
	_ = func() error {

		flag.Var(&configFiles, "config", "Config file for V2Ray. Multiple assign is accepted (only json). Latter ones overrides the former ones.")
		flag.Var(&configFiles, "c", "Short alias of -config")
		flag.StringVar(&configDir, "confdir", "", "A dir with multiple json config")

		return nil
	}()
)

var (
	osSignals chan os.Signal = make(chan os.Signal, 1)
	isopen int
)


func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func dirExists(file string) bool {
	if file == "" {
		return false
	}
	info, err := os.Stat(file)
	return err == nil && info.IsDir()
}

func readConfDir(dirPath string) {
	confs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		part.Logf().E(err.Error())
	}
	for _, f := range confs {
		if strings.HasSuffix(f.Name(), ".json") {
			configFiles.Set(path.Join(dirPath, f.Name()))
		}
	}
}

func getConfigFilePath() (cmdarg.Arg, error) {
	if dirExists(configDir) {
		part.Logf().I("Using confdir from arg:", configDir)
		readConfDir(configDir)
	} else {
		if envConfDir := platform.GetConfDirPath(); dirExists(envConfDir) {
			part.Logf().I("Using confdir from env:", envConfDir)
			readConfDir(envConfDir)
		}
	}

	if len(configFiles) > 0 {
		return configFiles, nil
	}

	if workingDir, err := os.Getwd(); err == nil {
		configFile := filepath.Join(workingDir, "config.json")
		if fileExists(configFile) {
			part.Logf().I("Using default config: ", configFile)
			return cmdarg.Arg{configFile}, nil
		}
	}

	if configFile := platform.GetConfigurationPath(); fileExists(configFile) {
		part.Logf().I("Using config from env: ", configFile)
		return cmdarg.Arg{configFile}, nil
	}

	part.Logf().I("Using config from STDIN")
	return cmdarg.Arg{"stdin:"}, nil
}

func GetConfigFormat() string {
	switch strings.ToLower(*format) {
	case "pb", "protobuf":
		return "protobuf"
	default:
		return "json"
	}
}

func startV2Ray() (core.Server, error) {
	configFiles, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	config, err := core.LoadConfig(GetConfigFormat(), configFiles[0], configFiles)
	if err != nil {
		return nil, newError("failed to read config files: [", configFiles.String(), "]").Base(err)
	}

	server, err := core.New(config)
	if err != nil {
		return nil, newError("failed to create server").Base(err)
	}

	return server, nil
}

func printVersion() {
	version := core.VersionStatement()
	for _, s := range version {
		part.Logf().I(s)
	}
}

func Start() (returnVal error) {

	defer func() {
		configFiles = cmdarg.Arg{}
		if err := recover(); err != nil {
			part.Logf().E("[v2]panic：", err.(error).Error())
			isopen = 2
            returnVal = err.(error)
		}
    }()

	part.Logf().New("log/v2.log")
	
	if b,_,e := part.Reqf(part.ReqfVal{
		Url:"http://127.0.0.1:8080/api/session?k=config",
	});e != nil {
		part.Logf().E("[v2]get config session：", e.Error())
	}else{
		configFiles.Set("http://127.0.0.1:8080/api/config?v="+string(b))
	}

	flag.Parse()

	printVersion()

	if *version {
		return nil
	}

	server, err := startV2Ray()
	if err != nil {
		part.Logf().E(err.Error())
		// Configuration error. Exit with a special value to prevent systemd from restarting.
		// os.Exit(23)
		return err
	}

	if *test {
		part.Logf().I("Configuration OK.")
		// os.Exit(0)
		return nil
	}

	if err := server.Start(); err != nil {
		part.Logf().E("Failed to start", err.Error())
		// os.Exit(-1)
		return err
	}
	defer server.Close()

	// Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()

	{
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		isopen = 1
		<-osSignals
		isopen = 0
	}

	part.Logf().I("v2ray exit with no error")

	return nil
}

func Stop(){
	osSignals <- syscall.SIGTERM
}

func IsOpen() int {
	return isopen
}

func newError(values ...interface{}) *errors.Error {
	return errors.New(values...).WithPathObj(errPathObjHolder{})
}