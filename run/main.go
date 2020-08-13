package run

import (
	"github.com/qydysky/part"
	"os"
)

func init(){
	s := part.Sys()
	path := s.Cdir()

	for s.CheckProgram("RUN")[0] != 0 {
		part.Logf().I("wait Run stop")
		part.Sys().Timeoutf(2)
	}

	if part.Checkfile().IsExist(path + "/RUN.run") {
		part.Logf().I("Run.run update")
		os.Rename(path + "/RUN.run", s.Pdir(path) + "/RUN.run")
	}
	if part.Checkfile().IsExist(path + "/RUN.exe") {
		part.Logf().I("Run.exe update")
		os.Rename(path + "/RUN.exe", s.Pdir(path) + "/RUN.exe")
	}
}