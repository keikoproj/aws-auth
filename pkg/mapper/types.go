/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mapper

import (
	"fmt"
	"reflect"

	"k8s.io/client-go/kubernetes"
)

type Bootstrapper struct {
	KubernetesClient kubernetes.Interface
}

func New(client kubernetes.Interface) *Bootstrapper {
	var bootstrapper = &Bootstrapper{}
	bootstrapper.KubernetesClient = client
	return bootstrapper
}

// AwsAuthData represents the data of the aws-auth configmap
type AwsAuthData struct {
	MapRoles []*AuthMap `yaml:"mapRoles"`
	MapUsers []*AuthMap `yaml:"mapUsers"`
}

// AddUniqueMapRole adds a unique AuthMap into MapRoles
func (m *AwsAuthData) AddUniqueMapRole(authMap *AuthMap) {
	for _, existingMap := range m.MapRoles {
		if reflect.DeepEqual(existingMap, authMap) {
			return
		}
	}
	if authMap.RoleARN == "" || authMap.Username == "" || len(authMap.Groups) == 0 {
		return
	}
	m.MapRoles = append(m.MapRoles, authMap)
}

// AddUniqueMapUser adds a unique AuthMap into MapUsers
func (m *AwsAuthData) AddUniqueMapUser(authMap *AuthMap) {
	for _, existingMap := range m.MapUsers {
		if reflect.DeepEqual(existingMap, authMap) {
			return
		}
	}
	if authMap.RoleARN == "" || authMap.Username == "" || len(authMap.Groups) == 0 {
		return
	}
	m.MapUsers = append(m.MapUsers, authMap)
}

// RemoveMapRole removes an auth map from MapRoles
func (m *AwsAuthData) RemoveMapRole(authMap *AuthMap) {
	var newMap []*AuthMap
	var match bool

	for _, existingMap := range m.MapRoles {
		if existingMap.RoleARN == authMap.RoleARN {
			match = true
			if len(authMap.Groups) != 0 {
				if reflect.DeepEqual(existingMap.Groups, authMap.Groups) {
					match = true
				} else {
					match = false
				}
			}
			if authMap.Username != "" {
				if authMap.Username == existingMap.Username {
					match = true
				} else {
					match = false
				}
			}
		}
		if !match {
			newMap = append(newMap, existingMap)
		}
	}

	if len(m.MapRoles) == len(newMap) {
		fmt.Printf("failed to remove %v, could not find exact match\n", authMap.RoleARN)
	} else {
		fmt.Printf("removed %v from aws-auth\n", authMap.RoleARN)
	}

	m.MapRoles = newMap
}

// RemoveMapUser removes an auth map from Mapusers
func (m *AwsAuthData) RemoveMapUser(authMap *AuthMap) {
	var newMap []*AuthMap
	var match bool

	for _, existingMap := range m.MapUsers {
		if existingMap.RoleARN == authMap.RoleARN {
			match = true
			if len(authMap.Groups) != 0 {
				if reflect.DeepEqual(existingMap.Groups, authMap.Groups) {
					match = true
				} else {
					match = false
				}
			}
			if authMap.Username != "" {
				if authMap.Username == existingMap.Username {
					match = true
				} else {
					match = false
				}
			}
		}
		if !match {
			newMap = append(newMap, existingMap)
		}
	}
	if len(m.MapUsers) == len(newMap) {
		fmt.Printf("failed to remove %v, could not find exact match\n", authMap.RoleARN)
	} else {
		fmt.Printf("removed %v from aws-auth\n", authMap.RoleARN)
	}
	m.MapUsers = newMap
}

// RemoveArguments are the arguments for removing a mapRole or mapUsers
type RemoveArguments struct {
	KubeconfigPath string
	MapRoles       bool
	MapUsers       bool
	Username       string
	RoleARN        string
	Groups         []string
}

// UpsertArguments are the arguments for upserting a mapRole or mapUsers
type UpsertArguments struct {
	KubeconfigPath string
	MapRoles       bool
	MapUsers       bool
	Username       string
	RoleARN        string
	Groups         []string
}

// AuthMap is the basic structure of an authentication object
type AuthMap struct {
	RoleARN  string   `yaml:"rolearn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

// NewAuthMap returns a new AuthMap
func NewAuthMap(rolearn, username string, groups []string) *AuthMap {
	return &AuthMap{
		RoleARN:  rolearn,
		Username: username,
		Groups:   groups,
	}
}

// SetRoleARN sets the RoleARN value
func (s *AuthMap) SetRoleARN(v string) *AuthMap {
	s.RoleARN = v
	return s
}

// SetUsername sets the Username value
func (s *AuthMap) SetUsername(v string) *AuthMap {
	s.Username = v
	return s
}

// SetGroups sets the Groups value
func (s *AuthMap) SetGroups(g []string) *AuthMap {
	s.Groups = g
	return s
}
