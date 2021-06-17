package yxw
import (
	"net/http"
	"fmt"
	"strings"
	"net/url"
	limit "github.com/qydysky/part/limit"
	reqf "github.com/qydysky/part/reqf"
)

func Yxw(){}

var web_total *limit.Limit 
func init(){
	web_total = limit.New(5,10,3000)
}

func Web(pattern string,web *http.ServeMux){
	web.HandleFunc(pattern+"/yxw", func(w http.ResponseWriter, r *http.Request) {
        if web_total.TO() {return}
        w.WriteHeader(404);return;
    })
	web.HandleFunc(pattern+"/yxw/", func(w http.ResponseWriter, r *http.Request) {
		if web_total.TO() {return}
		http.ServeFile(w, r, "./src/html/yxw/"+r.URL.Path[len(pattern+"/yxw/"):])
	})
	web.HandleFunc(pattern+"/yxw/api/", func(w http.ResponseWriter, r *http.Request) {
		if web_total.TO() {return}
        w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
        w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
        w.Header().Set("content-type", "application/json")
        fmt.Fprintf(w,webapi(r.URL))
	})
}

func webapi(MURL *url.URL) string{
	url:=MURL.EscapedPath()
	if strings.Count(url,"/pic/")==1 {return pic(url)}
	if strings.Count(url,"/seach/")==1 {return seachf(url)}
	if strings.Count(url,"/en2jp/")==1 {return en2jp(url)}
    return "404"
}

func pic(url string) string{
	tmp:=strings.Split(url,"/pic/")
	tmp=strings.Split(tmp[1],"&")

	var _ReqfVal = reqf.Rval{
		Url:"http://ocg.resource.m2v.cn/"+tmp[0]+".jpg",
	}
	req := reqf.New()
	if err:=req.Reqf(_ReqfVal);err==nil&&!strings.Contains(string(req.Respon),"Error") {return "\"http://ocg.resource.m2v.cn/"+tmp[0]+".jpg\""}

	return "\"http://ocg.resource.m2v.cn/ygopro/pics/"+tmp[1]+".jpg\""

}

func en2jp(url string) string{
	tmp:=strings.Split(url,"/en2jp/")
	
	if tmp[1]=="null"{return "err"}

	jp:=GetJanNameById(tmp[1])

	return jp
}

func seachf(url string) string{
	tmp:=strings.Split(url,"/seach")

	var _ReqfVal = reqf.Rval{
		Url:"https://www.ourocg.cn/search"+tmp[1],
	}
	req := reqf.New()
	if e:=req.Reqf(_ReqfVal);e != nil{}

	l:=string(req.Respon)
	
	l=l[strings.Index(l,"window.__STORE__ = ")+19:strings.LastIndex(l,";")]

	return l
}
