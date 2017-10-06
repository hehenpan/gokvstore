// It includes skill, equipment, card and so on
package main

import (
	//"database/sql"
	"flag"
	"fmt"
	"time"

	"net/http"
	"runtime"

	//"gokvstore/common/goredis"
	"gokvstore/common/libutil"
	"gokvstore/common/logging"

	//_ "github.com/go-sql-driver/mysql"
)

var cfg = struct {
	Log struct {
		File   string
		Level  string
		Name   string
		Suffix string
	}

	Prog struct {
		CPU        int
		Daemon     bool
		HealthPort string
	}

	Server struct {
		//Redis    string
		//Mysql    string
		PortInfo string
	}
}{}



func main() {
	//配置解析
	config := flag.String("c", "conf/config.json", "config file")
	flag.Parse()
	if err := libutil.ParseJSON(*config, &cfg); err != nil {
		fmt.Printf("parse config %s error: %s\n", *config, err.Error())
		return
	}

	//日志
	if err := libutil.TRLogger(cfg.Log.File, cfg.Log.Level, cfg.Log.Name, cfg.Log.Suffix, cfg.Prog.Daemon); err != nil {
		fmt.Printf("init time rotate logger error: %s\n", err.Error())
		return
	}
	if cfg.Prog.CPU == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU()) //配0就用所有核
	} else {
		runtime.GOMAXPROCS(cfg.Prog.CPU)
	}

	//日志
    logging.Debug("server start")


	//Mysql
	//InitMysql()
	//TestMysql()

	//Redis
	//InitRedis()
	//TestRedis()

	libutil.InitSignal()
/*
	go func() {
		err := http.ListenAndServe(cfg.Prog.HealthPort, nil)
		if err != nil {
			logging.Error("ListenAndServe: %s", err.Error())
		}
	}()
*/
	//registerHttpHandle()
	//httpclient
	//TestHttpClient()

	

	go func(){
		http.ListenAndServe(cfg.Server.PortInfo, nil)
	}()

	file, err := libutil.DumpPanic("gsrv")
	if err != nil {
		logging.Error("init dump panic error: %s", err.Error())
	}

	defer func() {
		logging.Info("server stop...:%d", runtime.NumGoroutine())
		time.Sleep(time.Second)
		logging.Info("server stop...,ok")
		if err := libutil.ReviewDumpPanic(file); err != nil {
			logging.Error("review dump panic error: %s", err.Error())
		}

	}()
	<-libutil.ChanRunning

}
