package main

import (
	"database/sql"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)


func getScanJobsFromDb(db *sql.DB) ([]ScanJob, error) {
  var scanjobs []ScanJob

  rows, err := db.Query("select d.id as domainid, d.domain, s.scanhost as host, s.scantype as type from domains d join scanrules s on s.scope = '*' left join records r on d.id = r.domainid and r.host = s.scanhost and r.type = s.scantype limit 100")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  for rows.Next() {
    var s ScanJob 
    if err := rows.Scan(&s.DomainID, &s.Domain, &s.Host, &s.Type); err != nil {
      return nil, err 
    }
    scanjobs = append(scanjobs, s)
  }

  return scanjobs, nil
}

func getDomainsFromDb(db *sql.DB) ([]Domain, error) { 
  var domains []Domain
  // pingerr := db.Ping()
  // if pingerr  != nil {
  //   return nil, pingerr 
  // }
  //
  // slog.Debug("db connection ping'd")
  //
  rows, err := db.Query("select id, domain, audit, lastaudit from domains where audit = 'yes'  and lastaudit = 0")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  for rows.Next() {
    var d Domain 
    if err := rows.Scan(&d.id, &d.domain, &d.audit, &d.lastaudit); err != nil {
      return nil, err 
    }
    domains = append(domains, d)
  }

  if err := rows.Err(); err != nil {
    return nil, err 
  }

  return domains, nil
}

func connectToDb() (*sql.DB, error)  {
  // println(viper.GetString("database.pass"))
  var db *sql.DB
  pass := viper.GetString("database.pass")
  user := viper.GetString("database.user")
  host := viper.GetString("database.host")
  name := viper.GetString("database.name")


  slog.Debug("dbcreds", "user", user, "pass", pass)
  dbconfig := mysql.Config{
    User: user,
    Passwd: pass,
    Net: "tcp",
    Addr: host,
    DBName: name,
    AllowNativePasswords: true,
  }

  db, err := sql.Open("mysql", dbconfig.FormatDSN())
  if err != nil {
    return nil, err
  }

  slog.Debug("db connection opened")
  pingerr := db.Ping()
  if pingerr  != nil {
    return nil, pingerr 
  }

  slog.Debug("db connection ping'd")

  return db, nil
}

