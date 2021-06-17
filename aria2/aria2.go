package aria2

import (
	reqf "github.com/qydysky/part/reqf"
	"github.com/qydysky/part"
	log "github.com/qydysky/part/log"
	"os/exec"
	"fmt"
	"strings"
)

var (
	prog *exec.Cmd
	aria2_log = log.New(log.Config{
        File:part.Sys().Pdir(part.Sys().Cdir())+`/log/aria2.log`,
        Stdout:true,
        Prefix_string:map[string]struct{}{`T: `:log.On,`I: `:log.On,`W: `:log.On,`E: `:log.On},
    }).Base(`aria2`)
)

func Aria2(){}

func Run(Replace map[string]string)*exec.Cmd{
	first(Replace)
	go func(){
		main()
	}()
	return prog
}

func main(){
	aria2_log.L(`I: `,"aria2 start")

	for check_and_close() {}

	part.Exec().Start(prog)

	err:=prog.Wait();

	if err == nil {
		aria2_log.L(`I: `,"fin with no error")
	}else{
		aria2_log.L(`E: `,err.Error())
	}
}

func check_and_close() bool {
	if part.Sys().CheckProgram(`aria2`)[0] > 0 {
		aria2_log.L(`I: `,"closeing aria2")
		req := reqf.New()
		if e:=req.Reqf(reqf.Rval{
			Url:`http://127.0.0.1:6800/jsonrpc?method=aria2.shutdown&id=op`,
		});e != nil {
			aria2_log.L(`W: `,e.Error())
		}
		part.Sys().Timeoutf(2)
		return true
	}
	return false
}

func first(Replace map[string]string){
	var (
		runFile=""
		rundir=""
	)

	if part.Checkfile().IsExist(part.Sys().Cdir()+"/ref/aria2") {
		rundir+=part.Sys().Cdir()+"/ref/aria2/"
	}else if part.Checkfile().IsExist(part.Sys().Cdir()+"/other/aria2/main") {
		rundir+=part.Sys().Cdir()+"/other/aria2/main/"
	}else if part.Checkfile().IsExist("main") {
		rundir+="main/"
	}

	if part.Sys().GetSys("windows") {
		runFile=rundir+"aria2c.exe"
	}else{
		runFile=rundir+"aria2c.run"
	}

	var u = part.Filel {
		File:rundir+"aria2.conf",
		Loc:0,
		ReadNum:0,
	}

	Replace["{dir}"] = rundir
	conf := part.File().FileWR(u)
	for k,v := range Replace{
		conf = strings.Replace(conf, k, v, -1 )
	}
	u.Context=[]interface{}{conf}
	u.File=rundir+"aria2.tmp.conf"

	part.File().FileWR(u)
	prog=exec.Command(runFile,fmt.Sprintf("--conf-path="+rundir+"aria2.tmp.conf"))
}