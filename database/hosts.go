package database

import (
	"fmt"
)

// TryDomain is current just set to hiper.dk since I work at hiper
var TryDomain = "hiper.dk"

func HostEquals(a, b string) bool {
	if a == b {
		return true
	}

	if fmt.Sprintf("%s.%s", a, TryDomain) == b {
		return true
	}

	return false
}

func GetDatabaseEntryConnection(hostQuery string, withSSL bool) *DBConn {
	hosts, err := ListHosts()
	if err != nil {
		panic(err)
	}

	for _, d := range hosts {
		match := HostEquals(hostQuery, d.Hostname)
		if !match {
			continue
		}

		config := MakeDbConfig(d.Hostname, d.Port, d.Username, d.Database)
		return config.Connect(withSSL)
	}

	return nil
}
