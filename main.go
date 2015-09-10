package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type NginxResponse struct {
	Connections struct {
		Accepted int64 `json:"accepted"`
		Active   int64 `json:"active"`
		Dropped  int64 `json:"dropped"`
		Idle     int64 `json:"idle"`
	} `json:"connections"`
	Upstreams map[string][]Backend `json:"upstreams"`
}

type Backend struct {
	ID     int    `json:"id"`
	Server string `json:"server"`
}

func NginxStatus() (*NginxResponse, error) {
	var nginxStatusServer string = "demo.nginx.com"
	resp, err := http.Get(fmt.Sprintf("http://%s/status", nginxStatusServer))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Non 200 OK")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var er NginxResponse
	if err := json.Unmarshal(data, &er); err != nil {
		return nil, err
	}
	return &er, nil
}

func main() {

	for {
		nt, err := NginxStatus()
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Millisecond * 1000)

		nr, err := NginxStatus()
		if err != nil {
			log.Println(err)
		}

		client, err := statsd.NewClient("127.0.0.1:8125", "nginx")
		if err != nil {
			log.Println(err)
		}

		var interval float32 = 1.0

		client.Inc("status.demo.connections.accepted", nr.Connections.Accepted-nt.Connections.Accepted, interval)
		client.Inc("status.demo.connections.dropped", nr.Connections.Dropped-nt.Connections.Dropped, interval)
		client.Inc("status.demo.connections.active", nr.Connections.Active, interval)
		client.Inc("status.demo.connections.idle", nr.Connections.Idle, interval)

		for _, backend := range nr.Upstreams {
			for _, server := range backend {
				fmt.Println(server.Server)
			}
		}

		// explicit close statsd connections
		client.Close()

	}
}
