package service

import (
	"log"
	"time"
	"os"
	"encoding/json"
	"github.com/tomhjx/network-probe/probe"
	"github.com/tomhjx/network-probe/resources"
)

type ProbePacket struct {
	target string
	response probe.Response
}

type Processor struct {
	target *resources.Target
	probePacket chan ProbePacket
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
		log.Printf("create dir %s", dirpath)
	}
	f, err := os.OpenFile(dirpath+"/"+t.Format("20060102")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(c))
	if _, err := f.WriteString(string(c)+"\n"); err != nil {
		log.Println(err)
		return
	}
}

func NewProcessor()(* Processor)  {
	return &Processor{}
}

func (proc *Processor) probe()  {
	for _, url := range proc.target.List {
		go func(url string) {
			r, err := probe.NewRequest(url)
			if err != nil {
				log.Println(err)
				return
			}
			resp, err := r.Run()
			if err != nil {
				log.Printf("request failed [%s], error %s", url, err)
				return
			}
	
			proc.probePacket <- ProbePacket{
				target: url,
				response: resp,
			}
			log.Printf("probePacket <- %s", url)
		}(url)
	}
}

func (proc *Processor) Run()(error)  {

	proc.probePacket = make(chan ProbePacket)

	// 定时更新探测目标
	var err error
	proc.target, err = resources.NewTarget()
	if err != nil {
		return err
	}
	proc.target.AutoUpdate(os.Getenv("TARGET_SOURCE_URL"))
	// 执行探测
	go func() {
		for {
			proc.probe()
			time.Sleep(1 * time.Second)
		}
	}()

	// 获取探测结果
	for {
		p := <- proc.probePacket
		preserve(p)
	}
}