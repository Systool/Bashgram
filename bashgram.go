package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type safeStore struct {
	mutex    sync.Mutex
	lastUpds []uint
}

const tgEndpoint = "https://api.telegram.org/bot711908048:AAGiRadEwO3cG93QtPKCn8ebBn2dj3JFPEU/"

func main() {
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGKILL)
	/*updatesDone := safeStore{mutex: sync.Mutex{}, lastUpds: make([]uint, 0)}
	if resp, err := http.Get(tgEndpoint + "getUpdates"); err == nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp := make([]APIResponse, 0)
			json.NewDecoder(resp.Body).Decode(&resp)
			resp.Body.Close()
			for i := 0; i < len(resp) && resp.Ok; i++ {
				message := Update{}
				resp[i].Result.UnmarshalJSON(resp[i].Result)
				append(updatesDone.lastUpd, resp[i].Result.UpdateID)
			}
		}
	} else {
		panic(err.Error())
	}*/
	pool := sync.WaitGroup{}
	//for {
	select {
	case _ = <-sigchan:
		break
	default:
		if resp, err := http.Get(tgEndpoint + "getUpdates"); err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				result := APIResponse{}
				json.NewDecoder(resp.Body).Decode(&result)
				resp.Body.Close()
				if result.Ok {
					updates := make([]Update, 0)
					if err := json.Unmarshal(result.Result, &updates); err != nil {
						fmt.Println(err.Error())
						//continue
					}
					for _, upd := range updates {
						pool.Add(1)
						go func(upd Update, pool *sync.WaitGroup) {
							fmt.Println(upd.UpdateID)
							pool.Done()
						}(upd, &pool)
					}
				}
			} else {
				fmt.Println(err)
			}
		}
	}
	pool.Wait()
	//}
fmt.Println("Exiting...")
}
