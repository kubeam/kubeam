// Copyright 2017 The Kubernetes Authors.
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

package services

import apiv1 "k8s.io/api/core/v1"

// KubeServicePort is a pair of port and protocol, e.g. a service endpoint.
type KubeServicePort struct {
	// Positive port number.
	Port int32 `json:"port"`

	// Protocol name, e.g., TCP or UDP.
	Protocol apiv1.Protocol `json:"protocol"`

	// The port on each node on which service is exposed.
	NodePort int32 `json:"nodePort"`
}

// GetServicePorts returns human readable name for the given service ports list.
func GetServicePorts(apiPorts []apiv1.ServicePort) []KubeServicePort {
	var ports []KubeServicePort
	for _, port := range apiPorts {
		ports = append(ports, KubeServicePort{port.Port, port.Protocol, port.NodePort})
	}
	return ports
}
