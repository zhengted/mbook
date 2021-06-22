package dynamiccache

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var (
	pool      *redis.Pool = nil
	MaxIdle   int         = 0
	MaxOpen   int         = 0
	ExpireSec int64       = 0
)

func InitCache() {
	addr := beego.AppConfig.String("dynamiccache_addrstr")
	if len(addr) == 0 {
		addr = "127.0.0.1:6379"
	}
	if MaxIdle <= 0 {
		MaxIdle = 256
	}
	if MaxOpen <= 0 {
		MaxOpen = 100
	}
	pass := beego.AppConfig.String("dynamiccache_passwd")
	if len(pass) == 0 {
		pool = &redis.Pool{
			MaxIdle:     MaxIdle,
			MaxActive:   MaxOpen,
			IdleTimeout: time.Duration(120 * time.Second),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp", addr,
					redis.DialReadTimeout(1*time.Second),
					redis.DialWriteTimeout(1*time.Second),
					redis.DialConnectTimeout(1*time.Second))
			},
		}
	} else {
		pool = &redis.Pool{
			MaxIdle:     MaxIdle,
			MaxActive:   MaxOpen,
			IdleTimeout: time.Duration(120 * time.Second),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp", addr,
					redis.DialReadTimeout(1*time.Second),
					redis.DialWriteTimeout(1*time.Second),
					redis.DialConnectTimeout(1*time.Second),
					redis.DialPassword(pass))
			},
		}
	}
}

func rdsdo(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}
	return con.Do(cmd, params...)
}

func WriteString(key string, value string) error {
	_, err := rdsdo("SET", key, value)
	rdsdo("EXPIRE", key, ExpireSec)
	return err
}

func ReadString(key string) (string, error) {
	res, err := rdsdo("GET", key)
	if err == nil {
		str, _ := redis.String(res, err)
		return str, nil
	} else {
		return "", err
	}
}

func WriteStruct(key string, obj interface{}) error {
	objJson, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return WriteString(key, string(objJson))
}

// 引用传递
func ReadStruct(key string, obj interface{}) error {
	reply, err := ReadString(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(reply), obj)
	return err
}

func WriteList(key string, list interface{}, total int) error {
	realKeyList := key + "_list"
	realKeyCount := key + "_count"
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = WriteString(realKeyCount, strconv.Itoa(total))
	if err != nil {
		return err
	}
	err = WriteString(realKeyList, string(data))
	if err != nil {
		rdsdo("del", realKeyCount)
		return err
	}
	return nil
}

func ReadList(key string, list interface{}) (int, error) {
	realKeyList := key + "_list"
	realKeyCount := key + "_count"
	if data, err := ReadString(realKeyList); nil == err {
		totalStr, _ := ReadString(realKeyCount)
		total := 0
		if len(totalStr) > 0 {
			total, _ = strconv.Atoi(totalStr)
		}
		return total, json.Unmarshal([]byte(data), list)
	} else {
		return 0, err
	}
}
