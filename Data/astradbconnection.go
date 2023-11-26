package data

import (
	"fmt"
	"os"
	"time"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
)

func connectToAstraDB() (*gocql.Session, error) {
	cluster, err := gocqlastra.NewClusterFromURL(gocqlastra.AstraAPIURL, os.Getenv("ASTRA_DB_ID"), os.Getenv("ASTRA_DB_TOKEN"), 10*time.Second)

	if err != nil {
		return nil, fmt.Errorf("unable to load cluster %s from astra: %v", os.Getenv("ASTRA_DB_TOKEN"), err)
	}

	cluster.Timeout = 30 * time.Second
	session, err := gocql.NewSession(*cluster)

	if err != nil {
		return nil, fmt.Errorf("unable to connect session: %v", err)
	}

	return session, nil
}
