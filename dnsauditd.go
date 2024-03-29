package main

import (
	"fmt"
	"log"
  "log/slog"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

const (
  semVer = "0.1.0"
  progName = "dnsauditd"
)

var flagVersion bool
var flagVerbose bool
var flagDebug bool
var flagTickrate int
var flagConfig string

var configFile os.FileInfo

// setupConfigFile wraps all of the config file parsing tasks
// it wil check some default common config file locations unless
// -c/--config <file> is set, then it will only check there and
// ignore any other paths 
// it will exit the program if config files are not present
func setupConfigFile() {
  confPaths := []string{
    "/etc/dnsauditd.conf",
    "/etc/dnsauditd/dnsauditd.conf",
    "./dnsauditd.conf",
  }

  configFile = nil

  if flagConfig != "" {
    slog.Info("-c/--config found, checking config", "path", flagConfig)
    file, err := os.Stat(flagConfig)
    if err != nil {
      slog.Info("not found", "path", flagConfig)
      slog.Error("invalid -c/--config file path", "path", flagConfig)
      os.Exit(4)
    }
    loadConfigFile(file)
  }

  for _, path := range confPaths {
    // checking if a configFile has already been loaded (via -c) 
    // if so, break out of this loop before checking the rest 
    if configFile != nil {
      break
    }

    slog.Info("checking for config file", "path", path)
    file, err := os.Stat(path)
    if err != nil {
      slog.Info("not found", "path", path)
      continue
    }

    slog.Info("config file found!", "path", path)
    loadConfigFile(file)
  }

  if configFile == nil {
    slog.Error("No valid config files found!")
    os.Exit(3)
  }
}

// loadConfigFile loads a given file into current running config 
func loadConfigFile(file os.FileInfo) {
  slog.Info("loaded config file!", "path", file.Name())
  configFile = file
}

// setupCliArgs wraps the various commandline arguments and options parsing
// and set up tasks for this program. It will also initiate the argparser 
// and handle basic housekeeping tasks like counting positional arguments 
// and handling arguments such as verson or help
func setupCliArgs () {
  // set up all commandline flags
  flag.BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")  
  flag.BoolVarP(&flagVersion, "version", "V", false, "print version")  
  flag.BoolVarP(&flagDebug, "debug", "D", false, "debug output")  
  flag.IntVarP(&flagTickrate, "tickrate", "T", 1000, "service tickrate in millseconds")  
  flag.StringVarP(&flagConfig, "config", "c", "", "config file to use")  
  flag.Parse()

  // if -v/--version is given, print version info and exit
  if flagVersion {
    fmt.Println("v", semVer)
    os.Exit(1)
  }

  // make sure that an incorrect number of args wasn't provided
  expectedArgs := 0
  if len(flag.Args()) != expectedArgs {
    flag.Usage()
    os.Exit(2)
  }
}

// setupLogger wraps the various logger setup tasks for this program
func setupLogger () {
  if flagDebug {
    slog.SetLogLoggerLevel(slog.LevelDebug)
  }
  log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
  log.SetPrefix(progName+": ")
}

func init() {
  setupCliArgs()
  setupLogger()
  setupConfigFile()
}

// daemonize simply sets up a ticker channel and main process loop
func daemonize() {
  // setup the ticker for the daemon
  delay := time.Duration(flagTickrate * int(time.Millisecond))
  ticker := time.NewTicker(delay)
  defer ticker.Stop()

  // main daemon loop
  for {
    select {
    case <-ticker.C:
      //do work
      slog.Info("hello world")
      slog.Debug("test debug message")
    }
  }
}

func main() { 
  // start the daemon
  daemonize()
}
