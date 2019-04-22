package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type safeStore struct {
	mutex    sync.Mutex
	lastUpds [10]uint
}

const (
	tgEndpoint = "https://api.telegram.org/bot<token>/"
	yourID     = <your ID>
)

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGKILL)
	pool := sync.WaitGroup{}
	lastUpd := 0
mainLoop:
	for {
		select {
		case <-sigchan:
			break mainLoop
		default:
			if resp, err := http.Get(
				tgEndpoint + `getUpdates?timeout=60` +
					`&offset=` + strconv.Itoa(lastUpd) +
					`&allowed_updates=["message"]`); err == nil {
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					result := APIResponse{}
					json.NewDecoder(resp.Body).Decode(&result)
					resp.Body.Close()
					if result.Ok {
						updates := make([]Update, 0)
						if err := json.Unmarshal(result.Result, &updates); err != nil {
							fmt.Println(err.Error())
						} else {
							for _, upd := range updates {
								if upd.Message.Chat.ID == yourID {
									lastUpd = upd.UpdateID + 1
									pool.Add(1)
									go func(upd Update, pool *sync.WaitGroup) {
										defer pool.Done()
										subpr := exec.Command("bash", "-c", `"`+upd.Message.Text+`"`)
										fmt.Println("Replying to update number", upd.UpdateID)
										/*if pipe, err := subpr.StdinPipe(); err != nil {
											fmt.Println(err)
											return
										} else {
											io.WriteString(pipe, upd.Message.Text)
										}*/
										var out string
										if outbytes, err := subpr.CombinedOutput(); err == nil {
											out = string(outbytes)
										} else {
											out = err.Error()
										}
										http.Get(
											tgEndpoint +
												"sendMessage?chat_id=" + strconv.Itoa(yourID) +
												"&text=" + url.QueryEscape(out) +
												"&reply_to_message_id=" + strconv.Itoa(upd.Message.MessageID))
									}(upd, &pool)
								}
							}
						}
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	}
	fmt.Println("Exiting...")
	pool.Wait()
}
