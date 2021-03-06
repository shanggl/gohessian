package gohessian

import (
  "bytes"
  "strings"
  "log"
  "runtime"
  "time"
  "unicode/utf8"
  "reflect"
  "github.com/shanggl/gohessian/util"
)

/*
对 基本数据进行 Hessian 编码
支持:
int8 int16 int32 int64
float64
time.Time
[]byte
[]interface{}
map[interface{}]interface{}
nil
bool
*/
/* shanggl
add marshal /unmashal struct object support
add Array support
*/

type Encoder struct {
}

const (
  CHUNK_SIZE    = 0x8000
  ENCODER_DEBUG = true
)

func init() {
  _, filename, _, _ := runtime.Caller(0)
  if ENCODER_DEBUG {
    log.SetPrefix(filename + "\n")
  }
}

func Encode(v interface{}) (b []byte, err error) {

  switch v.(type) {

  case []byte:
    b, err = encode_binary(v.([]byte))

  case bool:
    b, err = encode_bool(v.(bool))

  case time.Time:
    b, err = encode_time(v.(time.Time))

  case float64:
    b, err = encode_float64(v.(float64))

  case int:
    if v.(int) >= -2147483648 && v.(int) <= 2147483647 {
      b, err = encode_int32(int32(v.(int)))
    } else {
      b, err = encode_int64(int64(v.(int)))
    }

  case int32:
    b, err = encode_int32(v.(int32))

  case int64:
    b, err = encode_int64(v.(int64))

  case string:
    b, err = encode_string(v.(string))

  case nil:
    b, err = encode_null(v)

  case [] Any:
    b, err = encode_list(v.([]Any))

  case map[Any]Any:
    b, err = encode_map(v.(map[Any]Any))

  default:
    t:=reflect.TypeOf(v)
    if reflect.Ptr==t.Kind(){
	tmp:=reflect.ValueOf(v).Elem()
	t=reflect.TypeOf(tmp)
    }
    switch k:=t.Kind();k{
    case reflect.Struct:
	b,err=marshal(v)
    case reflect.Slice,reflect.Array:
	b,err=encode_list(v.([] Any))
    case reflect.Map:
	b,err=encode_map(v.(map[Any]Any))
    default:
	log.Printf("type not Support! %s",t.Kind().String())
    	panic("unknow type")
    }
  }
  if ENCODER_DEBUG {
    log.Println(util.SprintHex(b))
  }
  return
}

//=====================================
//对各种数据类型的编码
//=====================================

// binary
func encode_binary(v []byte) (b []byte, err error) {
  var (
    tag   byte
    len_b []byte
    len_n int
  )

  if len(v) == 0 {
    if len_b, err = util.PackUint16(0); err != nil {
      b = nil
      return
    }
    b = append(b, 'B')
    b = append(b, len_b...)
    return
  }

  r_buf := *bytes.NewBuffer(v)

  for r_buf.Len() > 0 {
    if r_buf.Len() > CHUNK_SIZE {
      tag = 'b'
      if len_b, err = util.PackUint16(uint16(CHUNK_SIZE)); err != nil {
        b = nil
        return
      }
      len_n = CHUNK_SIZE
    } else {
      tag = 'B'
      if len_b, err = util.PackUint16(uint16(r_buf.Len())); err != nil {
        b = nil
        return
      }
      len_n = r_buf.Len()
    }
    b = append(b, tag)
    b = append(b, len_b...)
    b = append(b, r_buf.Next(len_n)...)
  }
  return
}

// boolean
func encode_bool(v bool) (b []byte, err error) {
  if v == true {
    b = append(b, 'T')
  } else {
    b = append(b, 'F')
  }
  return
}

// date
func encode_time(v time.Time) (b []byte, err error) {
  var tmp_v []byte
  b = append(b, 'd')
  if tmp_v, err = util.PackInt64(v.UnixNano() / 1000000); err != nil {
    b = nil
    return
  }
  b = append(b, tmp_v...)
  return
}

// double
func encode_float64(v float64) (b []byte, err error) {
  var tmp_v []byte
  if tmp_v, err = util.PackFloat64(v); err != nil {
    b = nil
    return
  }
  b = append(b, 'D')
  b = append(b, tmp_v...)
  return
}

// int
func encode_int32(v int32) (b []byte, err error) {
  var tmp_v []byte
  if tmp_v, err = util.PackInt32(v); err != nil {
    b = nil
    return
  }
  b = append(b, 'I')
  b = append(b, tmp_v...)
  return
}

// long
func encode_int64(v int64) (b []byte, err error) {
  var tmp_v []byte
  if tmp_v, err = util.PackInt64(v); err != nil {
    b = nil
    return
  }
  b = append(b, 'L')
  b = append(b, tmp_v...)
  return

}

// null
func encode_null(v interface{}) (b []byte, err error) {
  b = append(b, 'N')
  return
}

// string
func encode_string(v string) (b []byte, err error) {
  var (
    len_b []byte
    s_buf = *bytes.NewBufferString(v)
    r_len = utf8.RuneCountInString(v)

    s_chunk = func(_len int) {
      for i := 0; i < _len; i++ {
        if r, s, err := s_buf.ReadRune(); s > 0 && err == nil {
          b = append(b, []byte(string(r))...)
        }
      }
    }
  )

  if v == "" {
    if len_b, err = util.PackUint16(uint16(r_len)); err != nil {
      b = nil
      return
    }
    b = append(b, 'S')
    b = append(b, len_b...)
    b = append(b, []byte{}...)
    return
  }

  for {
    r_len = utf8.RuneCount(s_buf.Bytes())
    if r_len == 0 {
      break
    }
    if r_len > CHUNK_SIZE {
      if len_b, err = util.PackUint16(uint16(CHUNK_SIZE)); err != nil {
        b = nil
        return
      }
      b = append(b, 's')
      b = append(b, len_b...)
      s_chunk(CHUNK_SIZE)
    } else {
      if len_b, err = util.PackUint16(uint16(r_len)); err != nil {
        b = nil
        return
      }
      b = append(b, 'S')
      b = append(b, len_b...)
      s_chunk(r_len)
    }
  }
  return
}

// list
func encode_list(v []Any) (b []byte, err error) {
  list_len := len(v)
  var (
    len_b []byte
    tmp_v []byte
  )

  b = append(b, 'V')

  if len_b, err = util.PackInt32(int32(list_len)); err != nil {
    b = nil
    return
  }
  b = append(b, 'l')
  b = append(b, len_b...)

  for _, a := range v {
    if tmp_v, err = Encode(a); err != nil {
      b = nil
      return
    }
    b = append(b, tmp_v...)
  }
  b = append(b, 'z')
  return
}

// map
func encode_map(v map[Any]Any) (b []byte, err error) {
  var (
    tmp_k []byte
    tmp_v []byte
  )
  b = append(b, 'M')
  for k, v := range v {
    if tmp_k, err = Encode(k); err != nil {
    }
    if tmp_v, err = Encode(v); err != nil {
      b = nil
      return
    }
    b = append(b, tmp_k...)
    b = append(b, tmp_v...)
  }
  b = append(b, 'z')
  return
}

/* 
	see example 
*/

// struct marshal to map

func marshal (v Any) (b []byte,err error) {

	var tmp_v []byte
	var s=reflect.ValueOf(v)
	var t=reflect.TypeOf(v)

//check Type exists
	mt:= s.MethodByName("GetType")	// mast contains Type Field to convert to object
	if !mt.IsValid(){ 
		log.Printf("Dos't contains GetType !")
		return
	}

	b=append(b,'M')
//encode type Name 
	b=append(b,'t')
	typeName:=mt.Call([] reflect.Value{})[0]	//call return [string,]
	var s_buf = *bytes.NewBufferString(typeName.String())
	typeNameLen:=utf8.RuneCount(s_buf.Bytes())
	if tmp_v, err = util.PackUint16(uint16(typeNameLen));tmp_v!= nil {
		b = append(b, tmp_v...)
	}
	for i := 0; i < typeNameLen; i++ {
		if r, s, err := s_buf.ReadRune(); s > 0 && err == nil {
		   b = append(b, []byte(string(r))...)
		}
	}
//encode the Fields
	for i:=0 ;i < t.NumMethod();i++ {
		var tmp_r [] reflect.Value
		var mi reflect.Method
		var tmp_s string

		mi=t.Method(i)

		if !strings.HasPrefix(mi.Name,"Get"){
			continue
		}
		if strings.EqualFold(mi.Name,"GetType"){
			continue  //jump type Field
		}
	//name change GetXaa to xaa
		if mi.Name[3]<'a' {
			tmp_s=string(mi.Name[3]+32)
		}else{
			tmp_s=string(mi.Name[3])
		}
		tmp_v,err=encode_string(tmp_s+mi.Name[4:])
		if err==nil && tmp_v!=nil{
			b=append(b,tmp_v...)
		}else{
			return nil,err
		}
	//value
		tmp_r=s.Method(i).Call([] reflect.Value{})//return [] reflect.Value
		tmp_v,err=Encode(tmp_r[0].Interface()) //GetXXX returns [string,]
		if err==nil && tmp_v!=nil{
			b=append(b,tmp_v...)
		}else {
			return nil,err
		}


	}//end of for 

	b=append(b,'z')
	return 
}

