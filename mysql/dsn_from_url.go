package mysql

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
)

// DsnFromUrl takes a standard connection url and converts it to a dsn that specifically works with github.com/go-sql-driver/mysql
// The most notable oddity is the hostname:port in the dsn being `tcp(ip:port)` instead of `ip:port`
func DsnFromUrl(connUrl string) (*mysql.Config, error) {
	u, err := url.Parse(connUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %w", err)
	}

	config := mysql.NewConfig()
	config.User = u.User.Username()
	config.Passwd, _ = u.User.Password()
	config.DBName = strings.TrimPrefix(u.Path, "/")
	config.Net = "tcp"
	config.Addr = u.Host
	for k, arr := range u.Query() {
		for _, v := range arr {
			config.Params[k] = v
		}
	}
	return config, nil
}
