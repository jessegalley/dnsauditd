package main

import (
	// "fmt"
	// "log"
	"database/sql"
	"fmt"
	// "fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
  semVer = "0.2.0"
  progName = "dnsauditd"
)

func init() {
  setupCliArgs()
  setupLogger()
}

func main() { 
  // setup the ticker for the daemon
  delay := time.Duration(viper.GetInt("tickrate") * int(time.Millisecond))
  ticker := time.NewTicker(delay)
  defer ticker.Stop()

  // try to connect to the DB and exit if we cannot
  var db *sql.DB
  db, err := connectToDb()
  if err != nil {
    slog.Error("couldn't connect to db", err)
    os.Exit(5) 
  }

  // main daemon loop
  for {
    select {
    case <-ticker.C:
      domains, err := getDomainsFromDb(db)
      if err != nil {
        slog.Error("err from getDomainsFromDb")
        panic(err)
      }
      slog.Info("domains to process", "count", len(domains))


      scanjobs, err := getScanJobsFromDb(db)
      if err != nil {
        slog.Error("couldn't get scanjobs", err)
        panic(err)
      }

      slog.Info("scanjobs to process", "count", len(scanjobs))

      for _, sj := range scanjobs {
        fmt.Println(sj) 
        slog.Debug("scanjob", "sj:", sj)
        scanrec, err := processScanJob(sj)
        if err != nil {
          slog.Error("process scan job failed", err)
          panic(err)
        }
        slog.Debug("scan result", "record", scanrec)
      }


      slog.Info("hello world")
      slog.Debug("hello Debug")
    }
  }
}
