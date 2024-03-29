package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var version = "0.1.0"

var flagVersion bool
var flagVerbose bool
var flagDebug bool

func init() {
  // set up all commandline flags
  flag.BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")  
  flag.BoolVarP(&flagVersion, "version", "V", false, "print version")  
  flag.BoolVarP(&flagDebug, "debug", "D", false, "debug output")  
  flag.Parse()

  // if -v/--version is given, print version info and exit
  if flagVersion {
    fmt.Println("v",version)
    os.Exit(1)
  }

  // make sure that an incorrect number of args wasn't provided
  expectedArgs := 0
  if len(flag.Args()) != expectedArgs {
    flag.Usage()
    os.Exit(2)
  }

}
func main() { 

  fmt.Println("hello")
}
