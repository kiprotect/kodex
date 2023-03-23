// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package gin

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/helpers"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func Router(controller api.Controller, prefix string, decorator gin.HandlerFunc) (*gin.Engine, error) {

	debug, _ := controller.Settings().Bool("debug")

	//we enable release mode until explicitly in debug mode
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()

	if decorator != nil {
		g.Use(decorator)
	}

	var meter kodex.Meter
	var err error

	if meter, err = helpers.Meter(controller.Settings()); err != nil {
		return nil, err
	}

	group, err := InitializeRouterGroup(g, prefix, controller)

	if err != nil {
		return nil, err
	}

	for _, routesProvider := range controller.APIDefinitions().Routes {
		if err := routesProvider(group, controller, meter); err != nil {
			return nil, err
		}
	}

	return g, nil

}

func RunApi(controller api.Controller, addr string, prefix string, handlerMaker func(api.Controller, http.Handler) (http.Handler, error), wg *sync.WaitGroup) (*http.Server, *gin.Engine, error) {

	g, err := Router(controller, prefix, nil)

	if err != nil {
		return nil, nil, err
	}

	var handler http.Handler

	kodex.Log.Info("Started API - listening on http://" + addr)

	if handlerMaker != nil {
		if handler, err = handlerMaker(controller, g); err != nil {
			return nil, nil, err
		}
	} else {
		handler = g
	}

	srv := &http.Server{Addr: addr, Handler: handler}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			kodex.Log.Error("HTTP Server Error - ", err)
		}
	}()

	return srv, g, nil

}
