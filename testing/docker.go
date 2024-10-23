package testing

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func GetDockerMysql(runOptions dockertest.RunOptions, hostConfig func(*docker.HostConfig)) (*dockertest.Resource, *sqlx.DB, error) {
	var db *sqlx.DB

	pool, err := dockertest.NewPool("")
	if nil != err {
		return nil, nil, fmt.Errorf("cannot create a new dockertest pool: %w", err)
	}

	err = pool.Client.Ping()
	if nil != err {
		return nil, nil, fmt.Errorf("cannot ping dockertest client: %w", err)
	}

	mysql, err := pool.RunWithOptions(&runOptions, hostConfig)
	if nil != err {
		return nil, nil, fmt.Errorf("cannot create dockertest mysql: %w", err)
	}

	err = pool.Retry(func() error {
		var err error
		db, err := sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", mysql.GetPort("3306/tcp")))
		if err != nil {
			return err
		}

		return db.Ping()
	})
	if nil != err {
		return nil, nil, fmt.Errorf("cannot open dockertest mysql connection")
	}

	return mysql, db, nil
}
