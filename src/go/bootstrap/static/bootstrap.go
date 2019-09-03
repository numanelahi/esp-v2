// Copyright 2019 Google Cloud Platform Proxy Authors
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

package static

import (
	"fmt"

	bt "cloudesf.googlesource.com/gcpproxy/src/go/bootstrap"
	gen "cloudesf.googlesource.com/gcpproxy/src/go/configgenerator"
	sc "cloudesf.googlesource.com/gcpproxy/src/go/configinfo"
	bootstrappb "github.com/envoyproxy/data-plane-api/api/bootstrap"
	ldspb "github.com/envoyproxy/data-plane-api/api/lds"
	conf "google.golang.org/genproto/googleapis/api/serviceconfig"
)

// ServiceToBootstrapConfig outputs envoy bootstrap config from service config.
// id is the service configuration ID. It is generated when deploying
// service config to ServiceManagement Server, example: 2017-02-13r0.
func ServiceToBootstrapConfig(serviceConfig *conf.Service, id string, options sc.EnvoyConfigOptions) (*bootstrappb.Bootstrap, error) {
	bootstrap := &bootstrappb.Bootstrap{
		Node:  bt.CreateNode(),
		Admin: bt.CreateAdmin(8001),
	}

	serviceInfo, err := sc.NewServiceInfoFromServiceConfig(serviceConfig, id, options)
	if err != nil {
		return nil, fmt.Errorf("fail to initialize ServiceInfo, %s", err)
	}

	listener, err := gen.MakeListener(serviceInfo)
	if err != nil {
		return nil, err
	}
	clusters, err := gen.MakeClusters(serviceInfo)
	if err != nil {
		return nil, err
	}

	bootstrap.StaticResources = &bootstrappb.Bootstrap_StaticResources{
		Listeners: []*ldspb.Listener{
			listener,
		},
		Clusters: clusters,
	}
	return bootstrap, nil
}