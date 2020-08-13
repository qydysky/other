package main

import (
	"github.com/qydysky/other/aria2"
	"github.com/qydysky/other/proxylist"
	"github.com/qydysky/other/yxw"
	_ "github.com/qydysky/other/run"
)

func main(){
	aria2.Aria2()
	proxylist.Proxylist()
	yxw.Yxw()
}