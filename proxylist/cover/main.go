// package main
package proxylistCov

import (
	"flag"
	"os"
	"io"
	// "fmt"
	"bufio"
	"strings"

	"github.com/qydysky/part"
)
type ProxylistCov_item struct {
	Url string
	Opath string
	Diff string
}

var (
	url_ = flag.String("url", "", "list zip")
	opath_ = flag.String("opath", "", "output path")
	diff_ = flag.String("diff","","diff")
)
func main(){
	flag.Parse()
	Main(ProxylistCov_item {
		Url:*url_,
		Opath:*opath_,
		Diff:*diff_,
	})
}

func Main(p ProxylistCov_item){
	

	url := p.Url
	opath := p.Opath
	diff := p.Diff


	if opath == "" {return}

	var buff,diff_buff,source_buff,_diff []string
	var e error

	if url!=""{
		var _ReqfVal = part.Rval{
			Url:url,
		}
		req := part.Req()
		if e:=req.Reqf(_ReqfVal);e!=nil {return}

		var u = part.Filel{
			File:"remote.zip",
			Write:true,
			Loc:0,
			Context:[]interface{}{req.Respon},
		} 
		part.File().FileWR(u)
		if e=part.Zip().UnZip("remote.zip",opath);e!=nil{return}
		os.Remove("remote.zip")
	}
	if e=NewPath(opath+"data/giaOut");e!=nil{return}
	if _,e=read_data_proxy(opath+"data/giaOut",&buff);e!=nil{return}

	if diff != "" {
		if e:=NewPath(opath+"diff/giaOut");e!=nil{return}
		if _,e=read_data_proxy(opath+"diff/giaOut",&diff_buff);e!=nil{return}
	}
	
	if _diff,e=read_data_proxy(opath+"giaOut",&buff);e!=nil{return}
	if diff != "" {
		for _,v:= range _diff {
			var have bool = false
			for _,l:=range diff_buff{
				if v==l {have=true;break}
			}
			if !have {diff_buff=append(diff_buff,v)}
		}
	}

	if e=write_proxy(opath+"data/giaOut",buff);e!=nil{return}
	if diff != "" {
		if e=write_proxy(opath+"diff/giaOut",diff_buff);e!=nil{return}
	}

	buff=source_buff
	diff_buff=source_buff
	_diff=source_buff

	if e=NewPath(opath+"data/fastOut");e!=nil{return}
	if _,e=read_data_proxy(opath+"data/fastOut",&buff);e!=nil{return}

	if diff != "" {
		if e=NewPath(opath+"diff/fastOut");e!=nil{return}
		if _,e=read_data_proxy(opath+"diff/fastOut",&diff_buff);e!=nil{return}
	}

	if _diff,e=read_data_proxy(opath+"fastOut",&buff);e!=nil{return}
	if diff != "" {
		for _,v:= range _diff {
			var have bool = false
			for _,l:=range diff_buff{
				if v==l {have=true;break}
			}
			if !have {diff_buff=append(diff_buff,v)}
		}
	}

	if e=write_proxy(opath+"data/fastOut",buff);e!=nil{return}
	if diff != "" {
		if e=write_proxy(opath+"diff/fastOut",diff_buff);e!=nil{return}
	}
	
	os.Remove(opath+"giaOut")
	os.Remove(opath+"fastOut")
}
func write_proxy(filename string, content []string) error {

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
	defer fd.Close()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fd)
	for k,v:= range content {
		_, err2 := w.WriteString("full:"+v)
		if err2 != nil {
			return err2
		}
		if k<len(content)-1{w.WriteString("\n")}
	}
	w.Flush()
	fd.Sync()
	return nil
}

func read_data_proxy(filename string,link_buff *[]string) ([]string,error) {
	var diff_buff []string
    fi, err := os.Open(filename)
    if err != nil {
        return diff_buff,err
    }
    defer fi.Close()

	br := bufio.NewReader(fi)
	var have_full bool
	
    for {
		a, _, c := br.ReadLine()
        if c == io.EOF||len(a)<1 {
            break
		}
		la := string(a)
		if !have_full && strings.Index(la,"full:") != -1 {have_full=true}
		if have_full {la=la[5:]}
		had := false
		for _,v:= range *link_buff {
			if v == la {had=true;break}
		}
		if !had {
			*link_buff=append(*link_buff,la)
			diff_buff=append(diff_buff,la)
		}
	}
	return diff_buff,nil
}

func NewPath(filename string) error {
	var newpath func(string) error = func (filename string)error{
		/*
			如果filename路径不存在，就新建它
		*/	
		var exist func(string) bool = func (s string) bool {
			_, err := os.Stat(s)
			return err == nil || os.IsExist(err)
		}
	
		for i:=0;true;{
			a := strings.Index(filename[i:],"/")
			if a == -1 {break}
			if a == 0 {a = 1}//bug fix 当绝对路径时开头的/导致问题
			i=i+a+1
			if !exist(filename[:i-1]) {
				err := os.Mkdir(filename[:i-1], os.ModePerm)
				if err != nil {return err}
			}
		}
		
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fd,err:=os.Create(filename)
			if err != nil {
				return err
			}else{
				fd.Close()
			}
		}
		return nil
	}
	return newpath(filename)

}