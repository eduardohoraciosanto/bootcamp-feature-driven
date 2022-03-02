package cache_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/cache"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/go-redis/redismock/v8"
)

var testLogger = logger.NewLogger(false, "redis mock")

func TestSetOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	b, _ := json.Marshal("test")
	mock.ExpectSet("testKey", string(b), 0).SetVal("test")
	c := cache.NewRedisCache(testLogger, 0, db)

	if c.Set(context.TODO(), "testKey", "test") != nil {
		t.Fatalf("Error was not expected")
	}
}

func TestSetUnmarshallError(t *testing.T) {
	db, _ := redismock.NewClientMock()
	c := cache.NewRedisCache(testLogger, 0, db)

	if c.Set(context.TODO(), "testKey", make(chan int)) == nil {
		t.Fatalf("Error was expected")
	}
}

func TestSetCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	b, _ := json.Marshal("test")
	mock.ExpectSet("testKey", string(b), 0).SetErr(fmt.Errorf("mocked error"))
	c := cache.NewRedisCache(testLogger, 0, db)

	if c.Set(context.TODO(), "testKey", "test") == nil {
		t.Fatalf("Error was expected")
	}
}

func TestGetOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	b, _ := json.Marshal("test")
	mock.ExpectGet("testKey").SetVal(string(b))
	c := cache.NewRedisCache(testLogger, 0, db)
	str := ""
	if c.Get(context.TODO(), "testKey", &str) != nil {
		t.Fatalf("Error was not expected")
	}
	if str != "test" {
		t.Fatalf("Wrong Value fetched")
	}
}

func TestGetCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("testKey").SetErr(fmt.Errorf("cache Error"))
	c := cache.NewRedisCache(testLogger, 0, db)
	str := ""
	if c.Get(context.TODO(), "testKey", &str) == nil {
		t.Fatalf("Error was expected")
	}
}

func TestGetUnmarshalFailure(t *testing.T) {
	db, mock := redismock.NewClientMock()
	b, _ := json.Marshal("test")
	mock.ExpectGet("testKey").SetVal(string(b))
	c := cache.NewRedisCache(testLogger, 0, db)
	hereImpossible := make(chan int)
	if c.Get(context.TODO(), "testKey", &hereImpossible) == nil {
		t.Fatalf("Error was expected")
	}
}

func TestDeleteOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectDel("testKey").SetVal(1)
	c := cache.NewRedisCache(testLogger, 0, db)
	if c.Del(context.TODO(), "testKey") != nil {
		t.Fatalf("Error was not expected")
	}
}

func TestDeleteKeyNotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectDel("testKey").SetVal(0)
	c := cache.NewRedisCache(testLogger, 0, db)
	if c.Del(context.TODO(), "testKey") == nil {
		t.Fatalf("Error was expected")
	}
}

func TestDeleteCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectDel("testKey").SetErr(fmt.Errorf("cache Error"))
	c := cache.NewRedisCache(testLogger, 0, db)
	if c.Del(context.TODO(), "testKey") == nil {
		t.Fatalf("Error was expected")
	}
}

func TestPingOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectPing().SetVal("ok")
	c := cache.NewRedisCache(testLogger, 0, db)
	if c.Alive(context.TODO()) != true {
		t.Fatalf("true was expected")
	}
}

func TestPingError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectPing().SetErr(fmt.Errorf("Cache not ready"))
	c := cache.NewRedisCache(testLogger, 0, db)
	if c.Alive(context.TODO()) == true {
		t.Fatalf("true was not expected")
	}
}
