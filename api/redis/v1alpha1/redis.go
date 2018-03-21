// Copyright 2018 The redis-operator Authors
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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


// RedisList is a list of redis clusters.
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisSpec   `json:"spec"`
	Status            RedisStatus `json:"status"`
}

// RedisSpec represents a Redis spec
type RedisSpec struct {
	// Redis defines its failover settings
	Redis RedisSettings `json:"redis,omitempty"`

	// Sentinel defines its failover settings
	Sentinel SentinelSettings `json:"sentinel,omitempty"`

}

// RedisSettings defines the specification of the redis system
type RedisSettings struct {
	Replicas  int32                  `json:"replicas,omitempty"`
	Resources RedisResources 		 `json:"resources,omitempty"`
	Exporter  bool                   `json:"exporter,omitempty"`
	Version   string                 `json:"version,omitempty"`
}

// SentinelSettings defines the specification of the sentinel cluster
type SentinelSettings struct {
	Replicas  int32                  `json:"replicas,omitempty"`
	Resources RedisResources 	     `json:"resources,omitempty"`
}

// RedisResources sets the limits and requests for a container
type RedisResources struct {
	Requests CPUAndMem `json:"requests,omitempty"`
	Limits   CPUAndMem `json:"limits,omitempty"`
}

// CPUAndMem defines how many cpu and ram the container will request/limit
type CPUAndMem struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// RedisStatus has the status of the system
type RedisStatus struct {
	Phase      Phase       `json:"phase"`
	Conditions []Condition `json:"conditions"`
	Master     string      `json:"master"`
}

// Phase of the RF status
type Phase string

// Condition saves the state information of the redis system
type Condition struct {
	Type           ConditionType `json:"type"`
	Reason         string        `json:"reason"`
	TransitionTime string        `json:"transitionTime"`
}

// ConditionType defines the condition that the redis can have
type ConditionType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object