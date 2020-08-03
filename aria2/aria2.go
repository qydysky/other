package aria2

import (
	"github.com/qydysky/part"
	"os/exec"
	"fmt"
	"strings"
)

var prog *exec.Cmd

func Aria2(){}

func Run()*exec.Cmd{
	first()
	go func(){
		main()
	}()
	return prog
}

func main(){
	
	part.Exec().Startf(prog)

	err:=prog.Wait();

	if err == nil {
		part.Logf().I("aria2 fin with no error")
	}else{
		part.Logf().I("[error]aria2:"+err.Error())
	}
}

func first(){
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
		Write:false,
		Loc:0,
		ReadNum:0,
	}
	u.Context=strings.Replace(part.File().FileWR(u), "{dir}", rundir, -1 )
	u.Write=true
	u.File=rundir+"aria2.tmp.conf"

	part.File().FileWR(u)
	prog=exec.Command(runFile,fmt.Sprintf("--conf-path="+rundir+"aria2.tmp.conf"))
}