package generator
import (
	"reflect"
	"log"
)

var typeReg map[string] reflect.Type

func init(){
	typeReg=make(map[string] reflect.Type)
}


func ShowReg(){
	for k,v := range typeReg{
		log.Println("-->> show Registered types <<----")
		log.Println(k,v)
	}

}

func Reg(s string ,t reflect.Type){
	typeReg[s]=t
}
//checkExists
func HasReg(s string) bool {
	if len(s)==0 {
		return false
	}
	_,exists:=typeReg[s]
	return exists
}

//gen a new Instance of s type
func Gen(s string ) interface{} {
	var ret interface{}
	t,ok:=typeReg[s]
	if !ok{
		return nil
	}
	ret=reflect.New(t).Interface()
	log.Println("gen",t)
	return ret
}

