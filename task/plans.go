package task

import (
	"encoding/json"
	"github.com/sleagon/health-checker/base"
	"log"
	"sync"
	"time"
)

const cacheSize = 1024

type plan struct {
	Name     string `json:"-"`
	URL      string `json:"url"`
	Method   string `json:"method"`
	Body     string `json:"body"`
	Interval int64  `json:"interval"`
	Callback string `json:"callback"`
	Mail     string `json:"mail"`
}

var plans map[string]*plan
var gw sync.WaitGroup

// Dispatch dispatch jobs based on plans
func Dispatch() {
	for _, v := range plans {
		gw.Add(1)
		go v.Dispatch()
	}
	gw.Wait()

}

func (p *plan) Dispatch() {
	for {
		time.Sleep(time.Second * time.Duration(p.Interval))
		jobs <- job{p.Name, p.URL, p.Method, p.Body, time.Now()}
	}

}

func init() {
	plans = make(map[string]*plan)
	jobs = make(chan job, cacheSize)
	pf := base.GetConfig().Plans

	// raw error map
	var rem map[string]*json.RawMessage

	// load plan from json file
	if err := json.Unmarshal(*pf, &rem); err != nil {
		log.Panic("Failed to parse plan file.")
	}
	for k, v := range rem {
		plans[k] = new(plan)
		if json.Unmarshal(*v, plans[k]) != nil {
			log.Panic("Plan file is broken.")
		}
		plans[k].Name = k
	}
	log.Printf(">>>>>>>> %d plans loaded successfully.\n", len(plans))
}
