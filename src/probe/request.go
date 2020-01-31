package probe

import (
	"log"
	"os/exec"
	"encoding/json"
	"bytes"
)

// Request represents abstract request
type Request struct {
	target string
}

// NewRequest create a Request Object
func NewRequest(target string)(* Request, error)  {
	r := Request{
		target: target,
	}
	return &r, nil
}

// Run response request result
func (r *Request) Run()(Response, error) {
	log.Println("curl ", r.target)
	resp := Response{}
	fm := `{"NameLookUpTime":"%{time_namelookup}","ConnectTime":"%{time_connect}","AppConnectTime":"%{time_appconnect}","RedirectTime":"%{time_redirect}","PretransferTime":"%{time_pretransfer}","StarttransferTime":"%{time_starttransfer}","TotalTime":"%{time_total}","HTTPCode":"%{http_code}"}`
	cmd := exec.Command("curl", "-I", "-s", "-w", fm, r.target, "-o", "/dev/null")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return resp, err
	}
	
	if err := json.Unmarshal([]byte(stdout.String()), &resp); err != nil {
		return resp, err
	}
	return resp, nil
}