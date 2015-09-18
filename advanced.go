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
	Upstreams   map[string][]Backend   `json:"upstreams"`
	ServerZones map[string]interface{} `json:"server_zones"`
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

func SendStatsD(gs string, gt string, gv int64) {

	// connect to statsd
	var sd string = "127.0.0.1:8125"
	var gk string = "nginx"

	client, err := statsd.NewClient(sd, gk)
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	// send metrics to statsd
	var fi float32 = 10.0
	var mt string = "status.demo_nginx_com"

	client.Inc(mt+"."+gs+"."+gt, gv, fi)
}

func main() {

	for {

		// grab nginx plus status json (delayed)
		nt, err := NginxStatus()
		if err != nil {
			log.Println(err)
		}

		// sleep x seconds - time.Millisecond * 1000 = 1sec
		time.Sleep(time.Millisecond * 1000)

		// grab nginx plus status json
		nr, err := NginxStatus()
		if err != nil {
			log.Println(err)
		}

		// send nginx plus connection metrics
		SendStatsD("connections", "accepted_", nr.Connections.Accepted-nt.Connections.Accepted)
		SendStatsD("connections", "dropped", nr.Connections.Dropped-nt.Connections.Dropped)
		SendStatsD("connections", "active", nr.Connections.Active)
		SendStatsD("connections", "idle", nr.Connections.Idle)

		// testing loop of server zones
		for _, zone := range nr.ServerZones {

			fmt.Println(zone)
		}

		// loop through upstream
		for _, backend := range nr.Upstreams {

			fmt.Println(backend)

			for _, server := range backend {

				fmt.Println(server.Server)

			}
		}

	}
}
