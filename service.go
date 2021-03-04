package zlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	list map[string][]string	//可用服务列表
	etcd 	*MyEtcd
	option ServiceOption
}

type ServiceOption struct {
	Etcd 	*MyEtcd
	Log 	*Log
	Prefix	string
}

func NewService(serviceOption ServiceOption)*Service {
	service := new(Service)

	service.etcd = serviceOption.Etcd
	service.option = serviceOption
	service.RegThird()
	return service
}
//从配置中心读取：3方可用服务列表,注册到内存中
func (service *Service)RegThird( ){
	//从etcd 中读取，已注册的服务
	allServiceList := service.etcd.GetListByPrefix(service.option.Prefix)
	if len(allServiceList) == 0{
		service.option.Log.Notice( " allServiceList is empty !")
		return
	}

	serviceListMap := make(map[string][]string)
	for k,_ := range allServiceList{
		str := strings.Replace(k,service.option.Prefix,"",-1)
		//MyPrint(str,k)
		serviceArr := strings.Split(str,"/")
		serviceListMap[serviceArr[1]] = append(serviceListMap[serviceArr[1]], serviceArr[2])
	}

	service.list = serviceListMap
	service.option.Log.Debug(serviceListMap)
}
//注册自己的服务
func (service *Service)RegOne(serviceName string,ipPort string){
	now := GetNowTimeSecondToInt()
	putResponse,err := service.etcd.PutOne( service.option.Prefix +"/"+serviceName +"/"+ipPort , strconv.Itoa(now))
	if err != nil{
		ExitPrint("service.etcd.PutOne err ",err.Error())
	}
	service.option.Log.Info("etcd put one ",putResponse.Header)
}

func (service *Service)balanceHost(list []string)string{
	return list[0]
}


func (service *Service)HttpPost(serviceName string,uri string,data interface{}) (responseMsgST ResponseMsgST,errs error){
	serviceIpList ,ok := service.list[serviceName]
	if !ok {
		return responseMsgST,errors.New(serviceName + " 不存在 map 中 ")
	}
	serviceHost := service.balanceHost(serviceIpList)
	url := "http://"+serviceHost + "/" + uri
	//url := "http://192.168.31.46:8080/"+uri
	//ExitPrint(url)
	service.option.Log.Debug("HttpPost",serviceName,serviceHost,uri,url)
	jsonStr, _ := json.Marshal(data)
	service.option.Log.Debug("jsonStr:",jsonStr)
	//ExitPrint(1111)
	req, errs := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", "application/json")
	defer req.Body.Close()
	//
	if errs != nil {
		return responseMsgST,errors.New("NewRequest err")
	}
	//5秒超时
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	service.option.Log.Debug(resp,error)
	if error != nil {
		return responseMsgST,errors.New("client.Do  err"+error.Error())
	}

	if resp.StatusCode != 200{
		return responseMsgST,errors.New("http response code != 200")
	}

	if resp.ContentLength == 0{
		return responseMsgST,errors.New("http response content = 0")
	}
	contentJsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return responseMsgST,errors.New("ioutil.ReadAll err : "+err.Error() )
	}

	errs = json.Unmarshal(contentJsonStr,&responseMsgST)
	if errs != nil{
		return responseMsgST,errors.New(" json.Unmarshal html content err : "+err.Error() )
	}

	service.option.Log.Debug("responseMsgST : ",responseMsgST)
	return responseMsgST,nil
}