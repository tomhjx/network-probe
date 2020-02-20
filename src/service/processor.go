package service

import (
	"encoding/json"
	"log"
	"log/syslog"
	"os"
	"strconv"
	"time"

	"github.com/tomhjx/network-probe/probe"
	"github.com/tomhjx/network-probe/resources"
)

type ProbePacket struct {
	target   string
	response probe.Response
}

type Processor struct {
	target      *resources.Target
	probePacket chan ProbePacket
}

func preserve(p ProbePacket) {
	t := time.Now()

	type outrow struct {
		Target            string `json:"tg"`
		NameLookUpTime    string `json:"nlut"`
		ConnectTime       string `json:"cnt"`
		AppConnectTime    string `json:"acnt"`
		RedirectTime      string `json:"rt"`
		PretransferTime   string `json:"pt"`
		StarttransferTime string `json:"st"`
		TotalTime         string `json:"tt"`
		HTTPCode          string `json:"c"`
		UnixTime          int64  `json:"t"`
	}

	out := &outrow{
		Target:            p.target,
		NameLookUpTime:    p.response.NameLookUpTime,
		ConnectTime:       p.response.ConnectTime,
		AppConnectTime:    p.response.AppConnectTime,
		RedirectTime:      p.response.RedirectTime,
		PretransferTime:   p.response.PretransferTime,
		StarttransferTime: p.response.StarttransferTime,
		TotalTime:         p.response.TotalTime,
		HTTPCode:          p.response.HTTPCode,
		UnixTime:          t.Unix(),
	}

	c, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(c))
	// sysLog, err := syslog.Dial("unixgram", "/dev/syslog/log.sock", syslog.LOG_LOCAL1, "netprobe")
	sysLog, err := syslog.Dial("tcp", "rsyslog-agent:514", syslog.LOG_LOCAL1, "netprobe")

	if err != nil {
		log.Println(err)
		return
	}
	sysLog.Info(string(c))
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (proc *Processor) probe() {
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
				target:   url,
				response: resp,
			}
			log.Printf("probePacket <- %s", url)
		}(url)
	}
}

func (proc *Processor) Run() error {
	probeInterval, _ := strconv.ParseInt(os.Getenv("PROBE_INTERVAL_SECOND"), 10, 32)
	if probeInterval <= 0 {
		probeInterval = 10
	}
	proc.probePacket = make(chan ProbePacket)

	var err error
	// 定时更新探测目标
	proc.target, err = resources.NewTarget()
	if err != nil {
		return err
	}
	proc.target.AutoUpdate(os.Getenv("TARGET_SOURCE_URL"))
	// 执行探测
	go func(interval int64) {
		for {
			proc.probe()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}(probeInterval)

	// 获取探测结果
	for {
		p := <-proc.probePacket
		preserve(p)
	}
}
