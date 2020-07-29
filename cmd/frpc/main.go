// Copyright 2016 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/utils/log"
	"github.com/gorilla/mux"
)

const (
	httpReadTimeout  = 10 * time.Second
	httpWriteTimeout = 10 * time.Second
	addr             = "127.0.0.1"
	port             = 7300
)

type GeneralResponse struct {
	Code int
	Msg  string
}

func main() {
	_ = runLifecycleServer()
}

func runLifecycleServer() (err error) {
	// url router
	router := mux.NewRouter()

	// api
	router.HandleFunc("/api/start", apiStart).Methods("GET")
	router.HandleFunc("/api/restart", apiRestart).Methods("GET")
	router.HandleFunc("/api/stop", apiStop).Methods("GET")

	address := fmt.Sprintf("%s:%d", addr, port)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
	}
	if address == "" {
		address = ":http"
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	server.Serve(ln)
	return
}

// GET api/start
func apiStart(w http.ResponseWriter, r *http.Request) {
	resp := GeneralResponse{Code: 200}

	log.Info("Http request [/api/start]")
	defer func() {
		log.Info("Http response [/api/start], code [%d]", resp.Code)
		w.WriteHeader(resp.Code)
		if len(resp.Msg) > 0 {
			w.Write([]byte(resp.Msg))
		}
	}()

	err := sub.Start()
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
	} else {
		resp.Msg = "OK"
	}

}

// GET api/restart
func apiRestart(w http.ResponseWriter, r *http.Request) {
	resp := GeneralResponse{Code: 200}

	log.Info("Http request [/api/restart]")
	defer func() {
		log.Info("Http response [/api/restart], code [%d]", resp.Code)
		w.WriteHeader(resp.Code)
		if len(resp.Msg) > 0 {
			w.Write([]byte(resp.Msg))
		}
	}()

	err := sub.Restart()

	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
	} else {
		resp.Msg = "OK"
	}
}

// GET api/stop
func apiStop(w http.ResponseWriter, r *http.Request) {
	resp := GeneralResponse{Code: 200}

	log.Info("Http request [/api/stop]")
	defer func() {
		log.Info("Http response [/api/stop], code [%d]", resp.Code)
		w.WriteHeader(resp.Code)
		if len(resp.Msg) > 0 {
			w.Write([]byte(resp.Msg))
		}
	}()

	sub.Stop()
}
