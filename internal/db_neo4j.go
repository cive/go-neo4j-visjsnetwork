package internal

import (
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Neo4jConfig struct {
	username string
	password string
	uri      string
}

type Connection struct {
	driver  neo4j.Driver
	session neo4j.Session
}

func NewDB() *Neo4jConfig {
	var username string = os.Getenv("NEO4J_USER")
	if username == "" {
		log.Fatal("neo4j username is not found.")
		return nil
	}
	var password string = os.Getenv("NEO4J_PASS")
	if password == "" {
		log.Fatal("neo4j password is not found.")
		return nil
	}
	var uri string = os.Getenv("NEO4J_URI")
	if uri == "" {
		log.Fatal("neo4j uri is not found.")
		return nil
	}
	return &Neo4jConfig{username, password, uri}
}

func (conf Neo4jConfig) Connect(accessMode neo4j.AccessMode) *Connection {
	driver, err := neo4j.NewDriver(conf.uri, neo4j.BasicAuth(conf.username, conf.password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	session, err := driver.Session(accessMode)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &Connection{driver, session}
}

func (conn Connection) Close() error {
	err := conn.session.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = conn.driver.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}
