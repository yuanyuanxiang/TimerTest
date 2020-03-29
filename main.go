package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

const PoolTimeout = time.Millisecond

var mutex = sync.Mutex{}

var counter int

var last time.Time

var timers = sync.Pool{
	New: func() interface{} {
		t := time.NewTimer(PoolTimeout)
		t.Stop()
		return t
	},
}

func main() {
	const bind = "54321"
	fmt.Println("Bind port = ", bind)
	http.HandleFunc("/VIID/Faces", callabck)
	if err := http.ListenAndServe(":"+bind, nil); err != nil {
		fmt.Println(err)
	}
}

func callabck(w http.ResponseWriter, r *http.Request) {
	timer := timers.Get().(*time.Timer)
	timer.Reset(PoolTimeout)
	var m map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Unmarshal json failed", http.StatusBadRequest)
		return
	}

	now := time.Now()
	mutex.Lock()
	if now.Unix()-last.Unix() > 0 {
		counter++
		last = now
		fmt.Println(counter, " \t ", now.Format(time.RFC3339))
	}
	mutex.Unlock()

	<-timer.C
	timers.Put(timer)
	w.WriteHeader(http.StatusOK)
}
