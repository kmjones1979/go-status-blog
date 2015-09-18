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
	Address     string `json:"address"`
	Connections struct {
		Accepted int64 `json:"accepted"`
		Active   int64 `json:"active"`
		Dropped  int64 `json:"dropped"`
		Idle     int64 `json:"idle"`
	} `json:"connections"`
	Upstreams struct {
		Demoupstreams struct {
			Keepalive int64 `json:"keepalive"`
			Peers     []struct {
				Active       int64 `json:"active"`
				Backup       bool  `json:"backup"`
				Downstart    int64 `json:"downstart"`
				Downtime     int64 `json:"downtime"`
				Fails        int64 `json:"fails"`
				HealthChecks struct {
					Checks     int64 `json:"checks"`
					Fails      int64 `json:"fails"`
					LastPassed bool  `json:"last_passed"`
					Unhealthy  int64 `json:"unhealthy"`
				} `json:"health_checks"`
				ID        int64 `json:"id"`
				MaxConns  int64 `json:"max_conns"`
				Received  int64 `json:"received"`
				Requests  int64 `json:"requests"`
				Responses struct {
					OneXx   int64 `json:"1xx"`
					TwoXx   int64 `json:"2xx"`
					ThreeXx int64 `json:"3xx"`
					FourXx  int64 `json:"4xx"`
					FiveXx  int64 `json:"5xx"`
					Total   int64 `json:"total"`
				} `json:"responses"`
				Selected int64  `json:"selected"`
				Sent     int64  `json:"sent"`
				Server   string `json:"server"`
				State    string `json:"state"`
				Unavail  int64  `json:"unavail"`
				Weight   int64  `json:"weight"`
			} `json:"peers"`
		} `json:"demoupstreams"`
	} `json:"upstreams"`
}

func NginxStatus() (*NginxResponse, error) {

	// assign variable for nging plus server
	var nginxStatusServer string = "demo.nginx.com"

	//
	resp, err := http.Get(fmt.Sprintf("http://%s/status", nginxStatusServer))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Non 200 OK")
	}
	// clean up the connection
	defer resp.Body.Close()

	// read the body of the request into a variable
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshall the JSON data into a variable
	var er NginxResponse
	if err := json.Unmarshal(data, &er); err != nil {
		return nil, err
	}
	return &er, nil
}

func SendStatsD(gs string, gt string, gv int64) {

	// assign variables for statsd server and metric name prefix
	var statsdServer string = "127.0.0.1:8125"
	var gk string = "nginx"

	// connect to statsd
	client, err := statsd.NewClient(statsdServer, gk)
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	// assign variables for metric name and interval
	var fi float32 = 10.0
	var st string = "status"
	var sn string = "demo_nginx_com"

	// send metrics to statsd
	client.Inc(st+"."+sn+"."+gs+"."+gt, gv, fi)
}

func main() {

	for {

		// make a call to the NGINX status API function
		nr, err := NginxStatus()
		if err != nil {
			log.Println(err)
		}

		// print connection statistics to screen
		fmt.Println("NGINX Plus")
		fmt.Println("Connections Accepted:", nr.Connections.Accepted)
		fmt.Println("Connections Dropped:", nr.Connections.Dropped)
		fmt.Println("Connections Active:", nr.Connections.Active)
		fmt.Println("Connections Idle", nr.Connections.Idle)

		// loop upstream peers and print active connections for each upstream
		for _, np := range nr.Upstreams.Demoupstreams.Peers {
			fmt.Println("Server:", np.Server, "Active Connections:", np.Active)
		}

		// send metrics to statsd
		SendStatsD("connections", "accepted", nr.Connections.Accepted)
		SendStatsD("connections", "dropped", nr.Connections.Dropped)
		SendStatsD("connections", "active", nr.Connections.Active)
		SendStatsD("connections", "idle", nr.Connections.Idle)

		// sleep to desired time
		time.Sleep(time.Millisecond * 1000)

	}

}
