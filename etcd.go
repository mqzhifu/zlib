package zlib

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)
type ResponseMsgST struct {
	Code	int
	Msg 	interface{}
	//Code 	int `json:"code"`
	//Data 	interface{} `json:"data"`
}

type MyEtcd struct {
	cli *clientv3.Client
	option EtcdOption
	AppConflist map[string]string
}

type EtcdOption struct {
	AppName string
	AppENV	string
	FindEtcdUrl		string
	LinkAddressList	[]string
	Log *Log
}
type EtcdHttpResp struct {
	Code 	int `json:"code"`
	Data 	Etcdconfig `json:"data"`
}
type Etcdconfig struct {
	Username	string `json:"username"`
	Password	string	`json:"password"`
	Hosts		[]string `json:"hosts"`
}
func NewMyEtcdSdk(etcdOption EtcdOption)(myEtcd *MyEtcd,errs error){
	myEtcd = new (MyEtcd)
	htmlContentJson ,errs := getEtcdHostPort(etcdOption)
	if errs != nil {
		return nil,errors.New("http request err :" + errs.Error())
	}

	if len(htmlContentJson) == 0{
		return nil,errors.New("http request content empty! :" + errs.Error())
	}
	//MyPrint(string(htmlContentJson))
	jsonStruct :=  EtcdHttpResp{}
	errs = json.Unmarshal(htmlContentJson,&jsonStruct)
	if errs != nil {
		return nil,errors.New("http request err : Unmarshal " + errs.Error())
	}
	//etcdConfig := strings.Split(jsonStruct.Msg.(string),",")
	if len(jsonStruct.Data.Hosts) == 0 {
		return nil,errors.New("http request err : etcdConfig is empty ")
	}
	etcdOption.Log.Info("etcdConfig ip list : ", jsonStruct.Data.Hosts)
	etcdOption.LinkAddressList = jsonStruct.Data.Hosts

	cli, errs := clientv3.New(clientv3.Config{
		Endpoints:  jsonStruct.Data.Hosts,
		DialTimeout: 5 * time.Second,
		Username: jsonStruct.Data.Username,
		Password: jsonStruct.Data.Password,
	})
	//etcdOption.Log.Info("link etcd :",etcdConfig)
	if errs != nil {
		return nil,errors.New("clientv3.New error :  " + errs.Error())
	}
	myEtcd.cli = cli
	myEtcd.option = etcdOption

	myEtcd.iniAppConf()
	return myEtcd,nil
}
//寻找etcd host ip 列表
func getEtcdHostPort(etcdOption EtcdOption)( []byte,error){
	//url := "http://39.106.65.76:1234/system/etcd/cluster1/list/"
	etcdOption.Log.Info("find etcd host:port  : ",etcdOption.FindEtcdUrl)
	resp, errs := http.Get(etcdOption.FindEtcdUrl)
	if errs != nil{
		return nil,errs
	}
	htmlContentJson,_ := ioutil.ReadAll(resp.Body)
	return htmlContentJson,errs
}
//申请一个X秒TTL的租约
func (myEtcd *MyEtcd)NewLeaseGrand(ctx context.Context ,ttl int64,autoKeepAlive int)(clientv3.LeaseID,error){
	//创建一个租约实体
	lease :=  clientv3.NewLease(myEtcd.cli)
	//申请一个60秒的 租约 实体
	leaseGrant, err := lease.Grant(ctx, ttl)
	if  err != nil {
		myEtcd.option.Log.Error("lease.Grant err :",err.Error())
		return 0,err
	}
	if autoKeepAlive == 1{
		leaseKeepAliveResponse,err :=lease.KeepAlive(ctx,leaseGrant.ID)
		if err !=nil{
			myEtcd.option.Log.Error("lease.KeepAlive err :",err.Error(),leaseKeepAliveResponse)
			return 0,err
		}
	}
	myEtcd.option.Log.Info("create New Lease and Grand ,  ttl :",ttl, " id : ",leaseGrant.ID)
	return leaseGrant.ID,nil
}
//往一个租约里写入内容
func (myEtcd *MyEtcd)putLease(ctx context.Context,leaseId clientv3.LeaseID,k string,v string)(putResponse *clientv3.PutResponse,err error){
	//创建一个KV 容器
	kv := clientv3.KV(myEtcd.cli)
	myEtcd.option.Log.Info("putLease k:",k," v:",v)
	putResponse, err = kv.Put(ctx, k, v, clientv3.WithLease(leaseId))
	//myEtcd.option.Log.Info("putLease (",leaseId,"): ",putResponse, err)
	if err != nil{
		return putResponse,err
	}

	return putResponse,nil
}

func (myEtcd *MyEtcd)GetListByPrefix(key string)(list map[string]string){
	//myEtcd.option.Log.Info(" etcd GetListByPrefix , ",key ," : ")
	rootContext := context.Background()
	kvc := clientv3.NewKV(myEtcd.cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	response, err := kvc.Get(ctx, key,clientv3.WithPrefix())
	defer cancelFunc()
	//myEtcd.option.Log.Debug(" ",response, err)
	if err != nil {
		myEtcd.option.Log.Notice("client Get err : ",err.Error())
		return list
	}

	if response.Count == 0{
		return list
	}

	kvs := response.Kvs
	list = make(map[string]string)
	for _,v := range kvs{
		//MyPrint(string(v.Key),string(v.Value))
		list[string(v.Key)] =  string(v.Value)
	}
	//MyPrint(list)
	return list
}

func (myEtcd *MyEtcd)GetListValue(key string)(list []string){
	myEtcd.option.Log.Info(" etcd GetOne , ",key ," : ")
	rootContext := context.Background()
	kvc := clientv3.NewKV(myEtcd.cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	response, err := kvc.Get(ctx, key)
	myEtcd.option.Log.Debug(" ",response, err)
	if err != nil {
		myEtcd.option.Log.Error("Get",err)
	}
	cancelFunc()

	if response.Count == 0{
		return nil
	}

	kvs := response.Kvs

	for _,v := range kvs{
		list = append(list,string(v.Value))
	}
	return list
}

func (myEtcd *MyEtcd)GetOneValue(key string)string{
	myEtcd.option.Log.Info(" etcd GetOne , ",key ," : ")
	rootContext := context.Background()
	kvc := clientv3.NewKV(myEtcd.cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	response, err := kvc.Get(ctx, key)
	myEtcd.option.Log.Debug(" ",response, err)
	if err != nil {
		myEtcd.option.Log.Error("Get",err)
	}
	cancelFunc()

	if response.Count == 0{
		return ""
	}

	kvs := response.Kvs
	value := string( kvs[0].Value )
	return value
}
func (myEtcd *MyEtcd)SetLog(log *Log){
	myEtcd.option.Log = log
}
func (myEtcd *MyEtcd) PutOne(k string, v string)(putResponse *clientv3.PutResponse,errs error){
	myEtcd.option.Log.Info(" etcd PutOne: ",k , v)
	rootContext := context.Background()
	kvc := clientv3.NewKV(myEtcd.cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	defer cancelFunc()
	putResponse, errs = kvc.Put(ctx, k,v)

	if errs != nil {
		myEtcd.option.Log.Error("RegOneService : ",errs.Error())
		switch errs {
		case context.Canceled:
			myEtcd.option.Log.Error("ctx is canceled by another routine: %v", errs.Error())
		case context.DeadlineExceeded:
			myEtcd.option.Log.Error("ctx is attached with a deadline is exceeded: %v", errs.Error())
		//case rpctypes.ErrEmptyKey:
		//	log.Error("client-side error: %v", err)
		default:
			myEtcd.option.Log.Error("bad cluster endpoints, which are not etcd servers: %v", errs.Error())
		}
	}
	myEtcd.option.Log.Info("RegOneService success",putResponse.Header,putResponse.PrevKv)
	return putResponse, errs
}

func  (myEtcd *MyEtcd)Watch(key string) <-chan clientv3.WatchResponse {
	myEtcd.option.Log.Notice("etcd create new watch :",key)
	watchChan  := myEtcd.cli.Watch(context.TODO(),key,clientv3.WithPrefix())
	//MyPrint("return watchChan")
	return watchChan
	//rch := cli.Watch(context.Background(), "/xi")
}

func  (myEtcd *MyEtcd)getConfRootPrefix()string{
	rootPath := "/"+myEtcd.option.AppName + "/"+  myEtcd.option.AppENV + "/"
	return rootPath
}

func  (myEtcd *MyEtcd)iniAppConf() {
	myEtcd.option.Log.Info("etcd iniAppConf : ")
	confListEtcd := myEtcd.GetListByPrefix(myEtcd.getConfRootPrefix())
	if len(confListEtcd) == 0{
		return
	}
	confList := make(map[string]string)
	for k,v := range confListEtcd{
		str := strings.Replace(k,myEtcd.getConfRootPrefix(),"",-1)
		//serviceArr := strings.Split(str,"/")
		myEtcd.option.Log.Info("conf " , str,v)
		confList[str] = v
	}
	myEtcd.AppConflist = confList
}

func  (myEtcd *MyEtcd)GetAppConfByKey(key string)(str string){
	val := myEtcd.AppConflist[key]
	return val
}
