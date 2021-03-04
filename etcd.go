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
	Code 	int
	Msg 	interface{}
}

type MyEtcd struct {
	cli *clientv3.Client
	option EtcdOption
}

type EtcdOption struct {
	FindEtcdUrl		string
	LinkAddressList	[]string
	Log *Log
}


func getEtcdHostPort(etcdOption EtcdOption)( []byte,error){
	//configCenter = configcenter.NewConfiger(10,2,10,"ini")
	//configCenter.StartLoading("/data/www/golang/src/configcenter")
	//systemConfigStr,_ := configCenter.Search("system")

	//systemConfig := make(map[string]map[string]map[string]interface{})
	//json.Unmarshal([]byte(systemConfigStr),&systemConfig)
	//fmt.Printf("%+v",systemConfig)
	//zlib.ExitPrint(systemConfig,errs)
	//etcdConfigStr := systemConfig["system"]["groups"]["list"]

	//url := "http://39.106.65.76:1234/system/etcd/cluster1/list/"
	//etcdOption.Log.Info("getEtcdHostPort : ",etcdOption.FindEtcdUrl))
	resp, errs := http.Get(etcdOption.FindEtcdUrl)
	etcdOption.Log.Info(" get etcd config ip:port list : ",etcdOption.FindEtcdUrl,errs)
	if errs != nil{
		return nil,errs
	}
	htmlContentJson,_ := ioutil.ReadAll(resp.Body)
	return htmlContentJson,errs
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

	jsonStruct :=  ResponseMsgST{}
	errs = json.Unmarshal(htmlContentJson,&jsonStruct)
	if errs != nil {
		return nil,errors.New("http request err : Unmarshal " + errs.Error())
	}
	etcdConfig := strings.Split(jsonStruct.Msg.(string),",")
	if len(etcdConfig) == 0 {
		return nil,errors.New("http request err : etcdConfig is empty ")
	}
	etcdOption.Log.Info("etcdConfig : ", etcdConfig)
	etcdOption.LinkAddressList = etcdConfig

	cli, errs := clientv3.New(clientv3.Config{
		Endpoints:   etcdConfig,
		DialTimeout: 5 * time.Second,
	})
	etcdOption.Log.Info("link etcd :",etcdConfig)
	if errs != nil {
		return nil,errors.New("clientv3.New error :  " + errs.Error())
	}
	myEtcd.cli = cli

	myEtcd.option = etcdOption

	return myEtcd,nil
}

//oneMsg mvccpb.KeyValue
//"github.com/coreos/etcd/mvcc/mvccpb"
//type MyEtcdGetOneMsg struct {
//	Key				string
//	Create_revision	int
//	Mod_revision 	int
//	Version			int
//	Value			string
//}

func (myEtcd *MyEtcd)GetListByPrefix(key string)(list map[string]string){
	myEtcd.option.Log.Info(" etcd GetListByPrefix , ",key ," : ")
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
	//zlib.ExitPrint("aaaa",rs)
	//zlib.MyPrint(rs)
	//fmt.Println(kvs)
	//fmt.Printf("%+v",oneMsg.Key)
	//fmt.Printf("last value is :%s\r\n", string(kvs[0].Value))
	//os.Exit(-333)
	return value
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
	myEtcd.option.Log.Notice("etcd new watch :",key)
	//<-chan WatchResponse
	watchChan  := myEtcd.cli.Watch(context.Background(),key,clientv3.WithPrefix())
	//MyPrint("return watchChan")
	return watchChan
	//rch := cli.Watch(context.Background(), "/xi")
}
