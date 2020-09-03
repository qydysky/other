package wt
import (
	"net/http"
	"fmt"
	"strings"
	"mime/multipart"  
	"strconv"
	"io/ioutil"
	xpart "github.com/qydysky/part"
)

type wt struct{
	WebRoot string
	SavePath string
}

func WT() *wt{
	return &wt{}
}

var web_total = xpart.Limit(10,1000,3000)//every 1000ms accept 10 request and other wait 3000ms
func (t *wt)Web(web *http.ServeMux){

	web.HandleFunc("/wt", func(w http.ResponseWriter, r *http.Request) {
        if web_total.TO() {return}
        w.WriteHeader(404);return;
    })
	web.HandleFunc("/wt/", func(w http.ResponseWriter, r *http.Request) {
		if web_total.TO() {return}
		http.ServeFile(w, r, t.WebRoot + r.URL.Path[3:])
	})
	web.HandleFunc("/wt/api/", func(w http.ResponseWriter, r *http.Request) {
		if web_total.TO() {return}
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
        w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
        w.Header().Set("SN", "wt")
		w.Header().Set("Cache-Control", "max-age=0")
		returnVal := t.webapi(r)
		switch returnVal[0].(int) {
			case 0:{//"json"
				w.Header().Set("content-type", "application/json");
				break;
			}
			case 1:{//"html"
				w.Header().Set("content-type", "text/html; charset=utf-8");
				break;
            }
            case 2:{//"code"
                w.WriteHeader(returnVal[1].(int));
                return
            }
			default:;
		}
        fmt.Fprintf(w,returnVal[1].(string))
	})
}

func (t *wt)webapi(httpRq *http.Request) []interface{} {
    url:=httpRq.URL
	path:=url.Path
	
	if strings.Contains(path,"/list") {return t.list(httpRq)}
	if strings.Contains(path,"/upload") {return t.post(httpRq)}
	if strings.Contains(path,"/copy") {return t.copy(httpRq)}

    return []interface{}{2,404}
}

func (t *wt)copy(httpRq *http.Request) []interface{} {
	if httpRq.Method != "POST" {
		return []interface{}{2,503}
	}

	reqBody,err := ioutil.ReadAll(httpRq.Body)
	if err != nil {
		return []interface{}{0,err.Error()}
	}
	var i string
	if json := xpart.Json().GetValFrom(string(reqBody),"c");json != nil {
		i = json.(string)
	}else{
		return []interface{}{2,503}
	}

	file:=xpart.File()

	file.F.File = t.WebRoot + t.SavePath + "Clipboard.html"
	file.F.Write = true
	file.F.Context = []interface{}{i}

	file.FileWR(file.F)

	return []interface{}{0,"{\"error\":\"\"}"}
}

func (t *wt)post(httpRq *http.Request) []interface{} {
	if httpRq.Method != "POST" {
		return []interface{}{2,503}
	}

	if e := httpRq.ParseMultipartForm(1 << 10);e != nil {
		fmt.Println(e.Error())
		return []interface{}{0,"{\"error\":\""+e.Error()+"\"}"}
	}

	for _, fheaders := range httpRq.MultipartForm.File {
		for _, hdr := range fheaders {

			var (
				infile multipart.File
				e error
			)
			if infile, e = hdr.Open(); nil != e {
					return []interface{}{0,"{\"error\":\""+e.Error()+"\"}"}
			}

			file:=xpart.File()

			file.F.File = t.WebRoot + t.SavePath + hdr.Filename
			file.F.Write = true
			file.F.Context = []interface{}{infile}
		
			file.FileWR(file.F)
		}
	}

	return []interface{}{0,"{\"error\":\"\"}"}
}

func (t *wt)list(httpRq *http.Request) []interface{} {
	l,_,_ := xpart.Checkfile().GetAllFile(t.WebRoot + t.SavePath)
	var returnVal string = "{\"list\":["
	len := len(t.WebRoot + t.SavePath)
	for _,v := range l {
		returnVal += "\""+t.SavePath + v[len:] + "\","
	}
	returnVal = strings.TrimRight(returnVal, ",")
	returnVal += "]}"
	return []interface{}{0,returnVal}
}

func main() {
	web :=  http.NewServeMux()

	wt := WT()
	
	wt.WebRoot = "./html/"
	wt.SavePath = "save/"

	wt.Web(web)

	webAddr := "0.0.0.0"
	xpart.Port().Set("wt",8089)

	server := &http.Server{
		Addr:         webAddr+":"+strconv.Itoa(xpart.Port().Get("wt")),
		Handler:      web,
	}

	xpart.Logf().I("start:",server.Addr)
	xpart.Logf().I("open",server.Addr+"/wt/ to upload")
	server.ListenAndServe()

}