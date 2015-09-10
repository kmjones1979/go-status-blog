package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type StatusJSON struct {
	Connections struct {
		Accepted int `json:"accepted"`
		Active   int `json:"active"`
		Dropped  int `json:"dropped"`
		Idle     int `json:"idle"`
	} `json:"connections"`
}

type NginxResponse struct {
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

	//	// define the status api url
	//	var status_json string = "http://demo.nginx.com/status"
	//
	//	// request the json from the api using net/http
	//	x, err := http.Get(status_json)
	//
	//	if err != nil {
	//		log.Fatalf("ERROR: %s", err)
	//	}
	//
	//	x_dec := json.NewDecoder(x.Body)
	//
	//	var x_data StatusJSON
	//	if err := x_dec.Decode(&x_data); err == io.EOF {
	//		//break
	//	} else if err != nil {
	//		log.Fatal("ERROR: %s", err)
	//	}

	for {

		nr, err := NginxStatus()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(2 * time.Second))
		}

		for _, backend := range nr.Upstreams {
			for _, server := range backend {
				fmt.Println(server.Server, "\n")
			}
		}

		//for _, backend := range nr.NginxResponse.Backend {
		//	fmt.Println(backend)
		//}

		//		// print some metrics
		//		fmt.Println("NGINX Plus Metric\t", "Value")
		//		fmt.Println("Connections Accepted:\t", x_data.Connections.Accepted)
		//		fmt.Println("Connections Active:\t", x_data.Connections.Active)
		//		fmt.Println("Connections Dropped:\t", x_data.Connections.Dropped)
		//		fmt.Println("Connections Accepted:\t", x_data.Connections.Idle)
		//
		// loop through the zones or upstream

	}
}
