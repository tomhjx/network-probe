package resources

import (
	"log"
	"strings"
	"io/ioutil"
	"time"
	"net/http"
)

type Target struct {
	List []string
}

func NewTarget()(* Target, error)  {
	t := &Target{}
	return t, nil
}

func (t *Target) Update(source string) (error) {
	httpctx, err := http.Get(source)
	if httpctx != nil {
		defer httpctx.Body.Close()
	}
	if err != nil {
		log.Println(err)
		return err
	}
	b, err := ioutil.ReadAll(httpctx.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	t.List = strings.Split(string(b), "\n")
	return nil
}

func (t *Target) AutoUpdate(source string)  {
	go func() {
		tick := time.NewTicker(60*time.Second)
		for {
			t.Update(source)
			<-tick.C
		}
	}()
}