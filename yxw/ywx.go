package yxw
import (
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"github.com/qydysky/part"
)

func Web(web *http.ServeMux){
	if part.Limit(500,1,3).TO() {return}

	web.HandleFunc("/yxw/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./src/html/"+r.URL.Path)
	})
	web.HandleFunc("/yxw/api/", func(w http.ResponseWriter, r *http.Request) {
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

	var _ReqfVal = part.ReqfVal{
		Url:"http://ocg.resource.m2v.cn/"+tmp[0]+".jpg",
	}

	f,_,err:=part.Reqf(_ReqfVal);

	if err==nil&&!strings.Contains(string(f),"Error") {return "\"http://ocg.resource.m2v.cn/"+tmp[0]+".jpg\""}

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

	var _ReqfVal = part.ReqfVal{
		Url:"https://www.ourocg.cn/search"+tmp[1],
	}

	_l,_,_:=part.Reqf(_ReqfVal);

	l:=string(_l)
	
	l=l[strings.Index(l,"window.__STORE__ = ")+19:strings.LastIndex(l,";")]

	return l
}
