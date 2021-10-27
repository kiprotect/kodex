// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/helpers"
)

const (
	ApiVersion = "v0.1.0"
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

func Router(controller api.Controller, decorator gin.HandlerFunc) (*gin.Engine, error) {

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

	group, err := InitializeRouterGroup(g, controller)

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

func RegisterPlugins(controller api.Controller) error {
	pluginSettings, err := controller.Settings().Get("plugins")

	if err == nil {
		pluginsList, ok := pluginSettings.([]interface{})
		if ok {
			for _, pluginName := range pluginsList {
				pluginNameStr, ok := pluginName.(string)
				if !ok {
					return fmt.Errorf("expected a string")
				}
				if definition, ok := controller.Definitions().PluginDefinitions[pluginNameStr]; ok {
					plugin, err := definition.Maker(nil)
					if err != nil {
						return err
					}
					apiPlugin, ok := plugin.(api.APIPlugin)
					if ok {
						if err := controller.RegisterAPIPlugin(apiPlugin); err != nil {
							return err
						} else {
							kodex.Log.Infof("Successfully registered plugin '%s'", pluginName)
						}
					}
				} else {
					kodex.Log.Errorf("plugin '%s' not found", pluginName)
				}
			}
		}
	}
	return nil
}

func RunApi(controller api.Controller, addr string, wg *sync.WaitGroup) (*http.Server, *gin.Engine, error) {

	if err := RegisterPlugins(controller); err != nil {
		return nil, nil, err
	}

	g, err := Router(controller, nil)

	if err != nil {
		return nil, nil, err
	}

	kodex.Log.Info("Started API - listening on http://" + addr)

	srv := &http.Server{Addr: addr, Handler: g}

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
