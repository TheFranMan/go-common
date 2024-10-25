package testing

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/redis/go-redis/v9"
)

func GetDockerMysql(pool *dockertest.Pool, runOptions dockertest.RunOptions, hostConfig func(*docker.HostConfig), expire uint) (*dockertest.Resource, *sqlx.DB, error) {
	var db *sqlx.DB

	mysql, err := pool.RunWithOptions(&runOptions, hostConfig)
	if nil != err {
		return nil, nil, fmt.Errorf("cannot create dockertest mysql: %w", err)
	}

	mysql.Expire(expire)

	err = pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", mysql.GetPort("3306/tcp")))
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

func GetDockerRedis(pool *dockertest.Pool, runOptions dockertest.RunOptions, hostConfig func(*docker.HostConfig), expire uint) (*dockertest.Resource, *redis.Client, error) {
	var red *redis.Client

	redisResource, err := pool.RunWithOptions(&runOptions, hostConfig)
	if nil != err {
		return nil, nil, fmt.Errorf("cannot create dockertest redis container: %w", err)
	}

	redisResource.Expire(expire)

	err = pool.Retry(func() error {
		var err error
		red = redis.NewClient(&redis.Options{
			Addr:     "localhost:" + redisResource.GetPort("6379/tcp"),
			Password: "",
			DB:       0,
		})

		_, err = red.Ping(context.Background()).Result()
		return err
	})
	if nil != err {
		return nil, nil, fmt.Errorf("cannot open dockertest redis connection")
	}

	return redisResource, red, nil
}
