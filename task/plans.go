package task

import (
	"encoding/json"
	"github.com/sleagon/health-checker/base"
	"log"
	"sync"
	"time"
)

const cacheSize = 1024

var _loaded bool

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

func dispatch() {
	for _, v := range plans {
		gw.Add(1)
		go v.Dispatch()
	}
	gw.Wait()
}

// NewAndRun all tasks
func NewAndRun() {
	New()
	Run()
}

// Run run a instance
func Run() {
	go monit()
	go Consume()
	dispatch()
}

func (p *plan) Dispatch() {
	for {
		time.Sleep(time.Second * time.Duration(p.Interval))
		jobs <- job{p.Name, p.URL, p.Method, p.Body, time.Now()}
	}
}

// New init task package
func New() {
	if _loaded {
		return
	}
	defer func() { _loaded = true }()
	// load base package first
	base.New()
	// init basic var(s)
	pf := base.GetConfig().Plans
	plans = make(map[string]*plan)
	jobs = make(chan job, cacheSize)
	lastMail = make(map[string]time.Time)
	callbacks = make(chan callback, cacheSize)

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
