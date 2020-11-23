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
	"testing"
	"time"

	"github.com/intorch/config"
	"github.com/stretchr/testify/assert"
)

const healthURL string = "http://localhost:3000/health"

var (
	serverUp bool    = false
	status   *Status = nil
)

func TestNew(t *testing.T) {
	h := New()

	assert := assert.New(t)
	assert.NotEqual(h, nil, "The health object should not be nil")
	assert.Equalf(h.CurrentState, http.StatusOK, "Expected %d, but found %d", http.StatusOK, h.CurrentState)
	assert.Empty(h.ErrorMessages, "The ErrorMessages should be empty")
}

func startServer() *Status {
	if !serverUp {

		var jsonBody string = `{
			"docs": [
				{
					"name": "test", 
					"type": "health", 
					"data": {
						"route": "/health", 
						"addr": ":3000"
					}
				}
			]
		}`

		c := config.New([]byte(jsonBody))

		status = New()

		go func() {
			status.Start(c, "test")
		}()

		serverUp = true
		time.Sleep(1 * time.Second)
	}

	return status
}

func TestState_Status200(t *testing.T) {
	assert := assert.New(t)

	startServer()

	//status 200
	res, err := http.Get("http://localhost:3000/health")
	assert.Nil(err, "The get request should not return error")

	var target Status
	err = json.NewDecoder(res.Body).Decode(&target)
	res.Body.Close()

	assert.Nil(err, "The Decode should not return errror")
	assert.Equalf(target.CurrentState, http.StatusOK, "Expected %d, but found %d.", http.StatusOK, target.CurrentState)
}

func TestState_Status500(t *testing.T) {
	assert := assert.New(t)

	h := startServer()
	h.AddError(http.StatusInternalServerError, "Erro de Teste")

	res, err := http.Get("http://localhost:3000/health")
	assert.Nil(err, "The get request should not return error")

	defer res.Body.Close()

	var target Status
	err = json.NewDecoder(res.Body).Decode(&target)
	assert.Nil(err, "The Decode should not return errror")

	assert.Equalf(target.CurrentState, http.StatusInternalServerError, "Expected %d, but found %d.", http.StatusInternalServerError, target.CurrentState)
	assert.Lenf(target.ErrorMessages, 1, "Expected 1, but found %d.", len(h.ErrorMessages))
}
