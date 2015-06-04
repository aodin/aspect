package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// GooseConfig is a representation of database configurations for the
// goose migration tool: https://bitbucket.org/liamstask/goose
type GooseConfig map[string]struct {
	Driver string
	Open   string
}

// ParseGooseDatabase will parse a specific database name in the goose
// configuration file.
func ParseGooseDatabase(path, name string) (c Database, err error) {
	configs, err := ParseGooseYAML(path)
	if err != nil {
		return
	}
	var exists bool
	if c, exists = configs[name]; !exists {
		err = fmt.Errorf("config: no database named '%s'", name)
		return
	}
	return
}

// ParseGooseYAML will parse the entire goose database configuration file.
func ParseGooseYAML(path string) (conf map[string]Database, err error) {
	conf = make(map[string]Database)

	// TODO unmarshal directly from the file?
	f, err := os.Open(path)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	var goose GooseConfig
	if err = yaml.Unmarshal(b, &goose); err != nil {
		return
	}

	// Parse each config
	for name, db := range goose {
		c := Database{Driver: db.Driver}

		// Split the open string
		// TODO common operation for doing this
		attrs := strings.Split(db.Open, " ")

		// Where's my dynamic programming?
		m := make(map[string]string)
		for _, attr := range attrs {
			parts := strings.SplitN(attr, "=", 2)
			// Valid attrs will always have 2 parts
			if len(parts) != 2 {
				// TODO error?
				continue
			}
			m[parts[0]] = parts[1]
		}

		// Yup
		c.Host = m["host"]
		if c.Port, err = strconv.ParseInt(m["port"], 10, 64); err != nil {
			return
		}
		c.Name = m["dbname"]
		c.User = m["user"]
		c.Password = m["password"]
		c.SSLMode = m["sslmode"]
		conf[name] = c
	}
	return
}
