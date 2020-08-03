package yxw

import (
	// "fmt"
	"strings"
	"github.com/qydysky/part"
	"github.com/thedevsaddam/gojsonq"
)

func GetJanNameById(id string)string{
	code:=GetCodeById(id)
	
	if code == "" {return "op"}

	jp:=Getinf(
		GetJanPageUrl(code),
		"<title>",
		" |",
		7,
	);
	return jp
}

func Id(id string)interface {}{
	jq := gojsonq.New().File("src/ref/yxw.json").From("data")
	res := jq.Where("id", "=", id).Get()
    return res
}

func GetCodeById(id string)string{
	if len(Id(id).([]interface{})) == 0 {return ""}
	md, _ := Id(id).([]interface{})[0].(map[string]interface{})
	return md["card_sets"].([]interface{})[0].(map[string]interface{})["set_code"].(string)
}

func GetEngPageUrl(code string)string{
	return "https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=1&stype=4&request_locale=en&keyword="+code
}

func GetJanPageUrl(code string)string{
	return "https://www.db.yugioh-card.com"+Getinf(GetEngPageUrl(code),"link_value\" value=\"","\">",19)+"&request_locale=ja"
}

func Getinf(url,op,ed string,op_len int) string {
	var _ReqfVal = part.Rval{
		Url:url,
		Timeout:10,
		Retry:2,
	}

	_l,_,_:=part.Req().Reqf(_ReqfVal);
	
	l:=string(_l)
	
	var oop int

	oop=strings.Index(l,op);
	
	if oop==-1 {return "op"}

	l=l[oop+op_len:];
	
	oop=strings.Index(l,ed);

	if oop==-1 {return "ed"}

	return l[:oop];
}