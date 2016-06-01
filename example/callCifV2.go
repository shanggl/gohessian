package main

import (
	"fmt"
	"github.com/shanggl/gohessian/client"
	"github.com/shanggl/gohessian"
)

type QueryUserBaseInfoByUserNoRequest struct{
	Type	string		`key:"com.lz.cif.individual.queryinfo.facaed.dto.QueryUserBaseInfoByUserNoRequest"`
	ReqSysDate 	int64	`key:"reqSysDate"`
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
	url:="http://192.168.9.86/epbox-pcif"
	method:="invoke"
	var user QueryUserBaseInfoByUserNoRequest
	user.FlowNo="10010101010101"//flowNO
	user.OperationSys="" //sysType
	user.OperationCode="103"//syscode
	user.UserId=userId //userNo


	var	req UserQueryRequest
	req.Tx="cif_individial_0026_005"
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
	testQueryUserBaseInfo("100301323232323")
}
