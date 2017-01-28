package main

import (
	"github.com/robfig/cron"
	"time"
	"os/exec"
	"strings"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 2s", func(){
		out, _ := exec.Command("/bin/sh", "-c", `docker service create --restart-condition none alpine echo Hello`).Output()
		id := strings.Trim(string(out), "\n")
		println(id)
		time.Sleep(2 * time.Second)
		// TODO: Explore whether to use API requests instead `awk`
		out, _ = exec.Command("/bin/sh", "-c", `docker service ps ` + id + ` | tail -1 | awk '{print $6}'`).Output()
		println(string(out))
		if strings.HasPrefix(string(out), "Complete") {
			println("OK")
		} else {
			println("NOK")
		}
		println("END")
	})
	c.Start()
	time.Sleep(10 * time.Second)
	for _, e := range c.Entries() {
		println(e.Next.String())
	}
	c.Stop()
}
