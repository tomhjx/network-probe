package main

import (
	"log"
	"strings"
	"io/ioutil"
	"time"
	"net/http"
	"os"
	"encoding/json"
	"github.com/tomhjx/network-probe/probe"
)

type ProbePacket struct {
	target string
	response probe.Response
}


func upgres() ([]string, error)  {
	url := "http://test-mib-config.hk.ufileos.com/network-probe-targets"
	targets, err := http.Get(url)
	defer targets.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(targets.Body)
	if err != nil {
		log.Fatal(err)
	}
	rlist := strings.Split(string(b), "\n")
	return rlist, nil
}

func preserve(p ProbePacket)  {
	t := time.Now()

	type outrow struct {
		Target string `json:"tg"`
		NameLookUpTime string `json:"nlut"`
		ConnectTime string `json:"cnt"`
		AppConnectTime string `json:"acnt"`
		RedirectTime string `json:"rt"`
		PretransferTime string `json:"pt"`
		StarttransferTime string `json:"st"`
		TotalTime string `json:"tt"`
		HTTPCode string `json:"c"`
		UnixTime int64 `json:"t"`
	}

	out := &outrow{
		Target: p.target,
		NameLookUpTime: p.response.NameLookUpTime,
		ConnectTime: p.response.ConnectTime,
		AppConnectTime: p.response.AppConnectTime,
		RedirectTime: p.response.RedirectTime,
		PretransferTime: p.response.PretransferTime,
		StarttransferTime: p.response.StarttransferTime,
		TotalTime: p.response.TotalTime,
		HTTPCode: p.response.HTTPCode,
		UnixTime: t.Unix(),
	}

	c, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	dirpath := "/app/log/probe"
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		os.MkdirAll(dirpath, 0777)
	}
	f, err := os.OpenFile(dirpath+"/"+t.Format("20060102")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := f.WriteString(string(c)+"\n"); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	targets := []string{}
	// 定时更新探测目标
	go func() {
		var err error
		t := time.NewTicker(60*time.Second)
		for {
			targets, err = upgres()
			if err != nil {
				log.Fatal(err)
			}
			<-t.C
		}
	}()

	// 执行探测
	probePacket := make(chan ProbePacket)
	go func() {
		for {
			for _, url := range targets {
				r, err := probe.NewRequest(url)
				if err != nil {
					log.Fatal(err)
					return
				}

				resp, err := r.Run()
				if err != nil {
					log.Printf("error %s", err)
					continue
				}

				probePacket <- ProbePacket{
					target: url,
					response: resp,
				}
			}
	
			time.Sleep(1 * time.Second)
			log.Println("wait..")
		}
	}()

	// 获取探测结果
	for {
		p := <- probePacket
		preserve(p)	
	}
}