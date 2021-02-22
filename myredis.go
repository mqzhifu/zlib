package zlib

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type MyRedis struct {
	option RedisOption
	Conn redis.Conn
	connPool *redis.Pool
}


type RedisOption struct {
	Host string
	Port string
	Ps	string
	Log *Log
}

func   NewRedisConn(redisOption RedisOption)(*MyRedis,error){
	myRedis := new(MyRedis)
	redisOption.Log.Info("NewRedisConn : ",redisOption.Host,redisOption.Port)
	conn,error := redis.Dial("tcp",redisOption.Host+":"+redisOption.Port)
	if error != nil{
		return nil, error
	}
	myRedis.option = redisOption
	myRedis.Conn = conn
	return myRedis,nil
}

func NewRedisConnPool(redisOption RedisOption)(*MyRedis,error){
	myRedis := new(MyRedis)
	redisOption.Log.Info("NewRedisConn : ",redisOption.Host,redisOption.Port)
	myRedisPool  := &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     2,
		MaxActive:   20,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisOption.Host+":"+redisOption.Port)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", 0)
			return c, nil
		},
	}
	myRedis.option = redisOption
	myRedis.connPool = myRedisPool
	return myRedis,nil
}



func  (myRedis *MyRedis)RedisDo(commandName string, args ...interface{})(reply interface{}, error error){
	myRedis.option.Log.Debug("[redis]redisDo init:",commandName,args)
	conn := myRedis.connPool.Get()
	res,error := conn.Do(commandName,args... )
	defer conn.Close()
	//res,error := myRedis.Conn.Do(commandName,args... )
	if error != nil{
		myRedis.option.Log.Notice("redis err :",error.Error())
		return nil, error
	}
	//reflect.ValueOf(res).IsNil(),reflect.ValueOf(res).Kind(),reflect.TypeOf(res)
	//zlib.MyPrint("redisDo exec ,res : ",res," err :",err)
	return res,error
}

func  (myRedis *MyRedis)RedisDelAllByPrefix(prefix string){
	myRedis.option.Log.Notice(" action redisDelAllByPrefix : ",prefix)
	res,err := redis.Strings(  myRedis.RedisDo("keys",prefix))
	if err != nil{
		ExitPrint("redis keys err :",err.Error())
	}
	myRedis.option.Log.Debug("del element will num :",len(res))
	if len(res) <= 0 {
		myRedis.option.Log.Notice(" keys is null,no need del...")
		return
	}
	for _,v := range res{
		res,_ := myRedis.RedisDo("del",v)
		myRedis.option.Log.Debug("del key ",v , " ,  rs : ",res)
	}
}

func  (myRedis *MyRedis) redisDelAll(redisPrefix string){
	myRedis.RedisDelAllByPrefix( redisPrefix)
}