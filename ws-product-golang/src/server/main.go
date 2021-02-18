package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type counters struct {
	sync.Mutex
	view  int
	click int
}
type countersTwo struct {
	view  int
	click int
}

//similar to counters, without sync.mutex to avoid issues with copying sync.mutex value

var (
	c = counters{}

	content = []string{"sports", "entertainment", "business", "education"}
)

//holds new counters within each 5 second interval
var countersMap = make(map[string]countersTwo)

//the mock store in memory
var mockStore = make(map[string]countersTwo)


func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]
	currTime := time.Now().Format("2006-01-02 15:04")//gets time to support querying to counters by time
	fmt.Println(data)
	fmt.Println(currTime)

	key := data + ":" + currTime

	c.Lock()
	c.view++
	c.Unlock()

	countersMap[key] = countersTwo{view: c.view, click: c.click}//updates view value

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data, currTime)
	}
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string, currTime string) error {
	key := data + ":" + currTime

	c.Lock()
	c.click++
	c.Unlock()
	
	countersMap[key] = countersTwo{view: c.view, click: c.click}//updates click value

	return nil
}

var statsHandlerStart = time.Now()
var statsNumCalls = 0

//implementing fixed window technique with maximum of 5 requests for every 10 seconds
func statsHandler(w http.ResponseWriter, r *http.Request) {
	statsHandlerCurr := time.Now()

	diff := statsHandlerCurr.Sub(statsHandlerStart)//time since start of current interval

	if diff >= time.Second*10{//if we have to start a new 10 second interval
		statsNumCalls = 1
		statsHandlerStart = time.Now()

		if !isAllowed() {
			w.WriteHeader(429)
			return
		}
	}else{//if we are still in a 10 second interval
		if statsNumCalls == 5{//if we reached the max number of requests in the interval
			fmt.Fprint(w, "Too many requests")
		}else{//if we can still process requests
			statsNumCalls++
			if !isAllowed() {
				w.WriteHeader(429)
				return
			}
		}
	}
}

func isAllowed() bool {
	return true
}

func uploadCounters() error {
	return nil
}

func main() {

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)

	go addToStore()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//uploads counters to mock store every 5 seconds
func addToStore(){

	tk := time.NewTicker(5 * time.Second)

	for range tk.C{
		for key, value := range countersMap{
			mockStore[key] = countersTwo{view: value.view, click: value.click}
			delete(countersMap, key)
		}
		fmt.Println(mockStore)
	}
	
}	
