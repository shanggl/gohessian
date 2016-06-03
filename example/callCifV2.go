package main

import (
	"fmt"
	"time"
	"github.com/shanggl/gohessian/client"
	"github.com/shanggl/gohessian"
)

type QueryUserBaseInfoByUserNoRequest struct{
	Type	string	`key:"com.lz.cif.individual.queryinfo.facade.dto.QueryUserBaseInfoByUserNoRequest"`
	ReqSysDate 	time.Time	`key:"reqSysDate"`
	OperationCode 	string	`key:"operationCode"`
	OperationSys 	string	`key:"operationSys"`
	OperBy		string	`key:"operBy"`
	ChannelType 	string	`key:"channelType"`
	FlowNo		string	`key:"flowNo"`
	ChannelFlowNo	string	`key:"channelFlowNo"`
	IpAdress	string	`key:"ipAdress"`
	Imei		string	`key:"imei"`
	MacAdress	string	`key:"macAdress"`
	LoginId 	string	`key:"loginId"`
	UserId		string	`key:"userId"`
}

type UserQueryRequest struct{
	Type	string	`key:"com.rb.owk.wolfs.ebox.Request"`
	Tx	string	`key:"tx"`
	Version	string	`key:"version"`
	Args 	[] gohessian.Any	`key:"args"`
}



func testQueryUserBaseInfo(userId string){
	//basic info
	//url:="http://192.168.9.98:8380/pcif-web/ebox"
	url:="http://192.168.8.62:8180/pcif-web/ebox"
	method:="invoke"

	var user QueryUserBaseInfoByUserNoRequest
	user.FlowNo="10010101010101"//flowNO
	user.OperationSys="" //sysType
	user.OperationCode="103"//syscode
	user.UserId=userId //userNo


	var	req UserQueryRequest
	req.Tx="cif_individual_0026_0001"
	req.Version="1.0"
	req.Args=append(req.Args,user)
	
	rsp,err:=client.Request(url,method,req) 
	if err!=nil {
		fmt.Println("resp:",rsp)

	}else {
		print ("error",err)
	}


}

func main(){
	testQueryUserBaseInfo("110000001840808")
}
