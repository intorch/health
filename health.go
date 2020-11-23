// Copyright 2020 intorch.org. All Rights Reserved.
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

package health

import (
	"encoding/json"
	"net/http"

	"github.com/intorch/config"
)

//Conf struct to parse configuration
type Conf struct {
	Route string
	Addr  string
}

const healthComponent config.Component = "health"

//Status data structure to Health
type Status struct {
	CurrentState  int
	ErrorMessages []string
}

//New create new Health State object
func New() *Status {
	return &Status{
		CurrentState: http.StatusOK,
	}
}

func (st *Status) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(st.CurrentState)

	json.NewEncoder(w).Encode(st)
}

//Start init the healthcheck service that enable the health endpoint. The endpoint and
//port, both are provided by service configuration with 'conf.Health.Route' and
//'conf.Health.Addr' respectivelly
func (st *Status) Start(conf *config.Configuration, taskName string) {
	var health Conf

	item := conf.Get(taskName, healthComponent)
	item.Decode(&health)

	http.HandleFunc(health.Route, st.handleRequest)
	err := http.ListenAndServe(health.Addr, nil)

	print(err)
}

//AddError function do add new error to health check object. This error will be
//sent to kubernetes when it calls the configurated endpoint
func (st *Status) AddError(status int, err string) {
	st.CurrentState = status
	st.ErrorMessages = append(st.ErrorMessages, err)
}
