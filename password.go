package pgpass

import (
	"io"
	"net/url"
	"strings"
)

// PasswordFrom reads password for given host and user from r, which should
// be in a valid pgpass format. Host should be of the form "hostname:port".
func PasswordFrom(host, user string, r io.Reader) (string, error) {
	hp := strings.SplitN(host, ":", 2) // split to hostname:port
	if len(hp) == 1 {
		// Add default postgresql port
		hp = append(hp, "5432")
	}

	er := NewEntryReader(r)
	for er.Next() {
		e := er.Entry()
		if eq(hp[0], e.Hostname) && eq(hp[1], e.Port) && eq(user, e.Username) {
			return e.Password, nil
		}
	}
	return "", er.Err()
}

func eq(s, pattern string) bool {
	return pattern == "*" || s == pattern
}

// Password reads password for given host and user from a default pgpass file.
// Host should be of the form "hostname:port".
func Password(host, user string) (string, error) {
	f, err := OpenDefault()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	return PasswordFrom(host, user, f)
}

// UpdateURL injects password into URL if not already provided.
// Password will be loaded from the default pgpass file.
func UpdateURL(dbURL string) (string, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return "", err
	}
	if dbUser := u.User; dbUser != nil {
		if _, ok := dbUser.Password(); !ok {
			userName := dbUser.Username()
			pass, err := Password(u.Host, userName)
			if err != nil {
				return "", err
			}
			if pass != "" {
				u.User = url.UserPassword(userName, pass)
			}
		}
	}

	return u.String(), nil
}
