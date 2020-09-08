package proxylist

import (
    "bufio"
    "fmt"
    "io"
	"os"
	"strings"
	"errors"
	"github.com/qydysky/part"
)

var limit int = -1

type Main_proxy_type struct {
	Acce_file string
	Sign string
	Filename string
	Discard bool
}

func Main_proxy(A Main_proxy_type){
	var (
		acce_file string = A.Acce_file
		sign string = A.Sign
		filename string = A.Filename
		discard bool = A.Discard
		link_buff []string
	)

	// 写入文件
	// 判断文件是否存在
	var new func(string) error = func (filename string)error{
		
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
	
	if err:=new(filename);err != nil {fmt.Println(err);return}
	read_data_proxy(filename,&link_buff)
	read_proxy(acce_file,sign,&link_buff,proccess_proxy)

	//delete a domain record
	{
		L := len(link_buff)
		if discard && L >= 100 {
			t := int(part.Rand().MixRandom(0, int64(L)))
			link_buff = append(link_buff[:t],link_buff[t+1:]...)
		}
	}

	write_proxy(filename,link_buff)
}

func proccess_proxy(line,sign string,link_buff *[]string){
	if !check_proxy(line,sign) {return}

	link,err := cut_proxy(line)

	if !strings.Contains(link,".")||err!=nil {return}

	for k,v:= range *link_buff {
		if v == link {
			if k == 0 {
				return
			}else{
				*link_buff=append((*link_buff)[:k],(*link_buff)[k+1:]...)
			}
			break
		}
	}
	if limit<=0 || len(*link_buff) <= limit {
		*link_buff=append([]string{link},*link_buff...)
	}else{
		*link_buff=append([]string{link},(*link_buff)[:limit]...)
	}
}

func check_proxy(line,sign string)bool{
	return strings.Contains(line,sign)
}

func cut_proxy(line string)(string,error){
	begin:=strings.Index(line,"d tcp:")+6
	end:=strings.Index(line[begin:],":")+begin
	if begin==-1 || begin >= end  {return "",errors.New("N")}
	return line[begin:end],nil
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

func read_proxy(filename string,sign string,link_buff *[]string,f func(line,sign string,link_buff *[]string)) {

    fi, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    defer fi.Close()

    br := bufio.NewReader(fi)
    for {
        a, _, c := br.ReadLine()
        if c == io.EOF {
            break
		}
		f(string(a),sign,link_buff)
    }
}

func read_data_proxy(filename string,link_buff *[]string) {

    fi, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
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
		if !had {*link_buff=append(*link_buff,la)}
	}
}