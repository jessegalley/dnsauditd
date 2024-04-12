package main

import (
	// "fmt"
	"net"
	"strings"
)


type ScanJob struct {
  DomainID int
  Domain string
  Host string
  Type string 
}

type Record struct {
  ID  int
  DomainID int
  Host string 
  Type string 
  Data string 
  Updated string 
  Expires string 
}

func processScanJob(sj ScanJob) (Record, error) {
  var out Record
  out.DomainID = sj.DomainID
  out.Type = sj.Type
  out.Host = sj.Host

  data := ""
  // fmt.Println("t")
  switch sj.Type {
  case "a":
    //query a records
    ips, err := queryDNS(sj.Domain)
    if err != nil {
      return out, err
    }
    // fmt.Println(ips)
    for _, ip := range ips {
      ipv4 := ip.To4()
      if ipv4 == nil {
        continue
      }
      data += ipv4.String() + " "
    }
  case "mx":
    //query mx records 
  }

  out.Data = strings.TrimSpace(data) 

  return out, nil
}

func queryMX(domain string) ([]*net.MX, error) {
  mxs, err := net.LookupMX(domain)
  if err != nil {
    return nil, err
  }

  return mxs, nil
}

func queryDNS(domain string) ([]net.IP, error) {
  ips, err := net.LookupIP(domain)

  if err != nil {
    return nil, err 
  }

  return ips, nil
  
}
