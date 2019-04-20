package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const tgEndpoint = "https://api.telegram.org/bot<token>/"

func main() {
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGKILL)
	var lastUpd uint
	if resp, err := http.Get(tgEndpoint + "getUpdates"); err == nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			p := make([]APIResponse, 10)
			json.NewDecoder(resp.Body).Decode(&p)
			resp.Body.Close()
			for i := 0; i < len(p); i++ {
				/*if p[i].Result > lastUpd {
					lastUpd
				}*/
			}
		}
	} else {
		panic(err.Error())
	}
	for {
		select {
		case _ = <-sigchan:
			break
		default:
			if resp, err := http.Get(tgEndpoint + "getUpdates"); err == nil {
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					p := APIResponse{}
					json.NewDecoder(resp.Body).Decode(&p)
					if p.Ok {
						fmt.Println(p)
					}
					resp.Body.Close()
				}
			}
		}
	}
}
