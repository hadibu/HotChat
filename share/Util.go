package share

import (
	"log"
	"runtime"
	"encoding/json"
	"strings"
)

const (
	STACK_MAX = 100
)

func RecoverPanicStack() {
	if reco := recover(); reco != nil {
		log.Printf("%v", reco)
		PrintStack()
	}
}

func PrintStack() {
	for i := 1; i < STACK_MAX; i++ {
		funcName, fileName, line, ok := runtime.Caller(i)
		if ok {
			log.Printf("[func:%v,file:%v,line:%v]\n", runtime.FuncForPC(funcName).Name(), fileName, line)
		} else {
			break
		}
	}
}

func HasError(err error) bool {
	if err != nil {
		log.Println("Error = ", err)
		PrintStack()
		return true
	}
	return false
}

/**
* obj → string
*/
func JsonEncode( obj interface {}) (string,error) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Println("JsonEncode error:", err)
		PrintStack()
	}

	return string(b),err
}

/**
* string → obj
*/
func JsonDecode(str string,obj interface {}) (error) {
	dataBuf := string(str)
	dec := json.NewDecoder(strings.NewReader(dataBuf))
	err := dec.Decode(obj)
	if err != nil {
		log.Println("JsonDecode error:", err)
		PrintStack()
	}
	return err
}
