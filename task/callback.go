package task

import (
	"fmt"
	"github.com/sleagon/health-checker/base"
	"log"
	"net/http"
	"net/smtp"
	"time"
)

const mailInteval = 60

type callback struct {
	Name string
	URL  string
	Mail string
	Time time.Time
}

// record last time sending mail
var lastMail map[string]time.Time

// callback channel
var callbacks chan callback

// Fire fire the triggers
func (cb callback) Fire() {

	// fire callback
	http.Get(cb.URL)

	// make body of email
	text := fmt.Sprintf("Service %s is down at %s, trigger %s has been called. "+
		"If you got this message more than one time, you may need to fix the service manually.",
		cb.Name, cb.Time.Format("2006-1-3 06:04:05"), cb.URL)

	// send email async(ly)
	go cb.mail(text)
}

// send email
func (cb callback) mail(text string) {
	// within 5 min
	if lastMail[cb.Name].Add(time.Second * mailInteval).After(cb.Time) {
		return
	}
	log.Printf("Service %s is down, sending email now.\n", cb.Name)
	cfg := base.GetConfig().Mail
	from := cfg.Username
	passwd := cfg.Password
	to := cb.Mail
	svc := cb.Name
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Service %s is offline\n\n%s\n", from, to, svc, text)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		smtp.PlainAuth("", from, passwd, cfg.Host),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s\n", err)
		return
	}
	// update time
	lastMail[cb.Name] = cb.Time
}

// start moniting
func monit() {
	for v := range callbacks {
		go v.Fire()
	}
}
