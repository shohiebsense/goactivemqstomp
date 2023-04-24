package main

import (
	"log"

	"github.com/go-stomp/stomp/v3"

	"net/http"
	"time"
)

var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("guest", "guest"),
	stomp.ConnOpt.Host("/"),
}

func toTimeFormat(date string) string {
	return date + " 23:59:59"
}

func getCaseInsensitiveQuery(key string) string {
	return "LOWER(" + key + ") = LOWER(?)"
}

func SecureServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte("Secure Hello World.\n"))
}

func main() {

	conn, err := stomp.Dial("tcp", "localhost:61613", options...)

	if err != nil {
		println(err.Error())
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:

			error := conn.Send("/topic/SampleTopic", "application/json",
				[]byte("{\"browsers\":{\"firefox\":{\"name\":\"Firefox\",\"pref_url\":\"about:config\",\"releases\":{\"1\":{\"release_date\":\"2004-11-09\",\"status\":\"retired\",\"engine\":\"Gecko\",\"engine_version\":\"1.7\"}}}}}"),
				stomp.SendOpt.Header("activemq.subscriptionName", "SampleSubscription"))

			if error != nil {
				log.Panic(error.Error())
				return
			}

		case <-quit:
			ticker.Stop()
			conn.Disconnect()
			return
		}
	}
}
