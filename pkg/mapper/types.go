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
	"io/ioutil"
	"log"

	"k8s.io/client-go/kubernetes"
)

func init() {
	log.SetFlags(0)

}

type AuthMapper struct {
	KubernetesClient kubernetes.Interface
	LoggingEnabled   bool
}

func New(client kubernetes.Interface, isCommandline bool) *AuthMapper {
	var mapper = &AuthMapper{}
	mapper.KubernetesClient = client

	if !isCommandline {
		log.SetOutput(ioutil.Discard)
	}
	return mapper
}

// AwsAuthData represents the data of the aws-auth configmap
type AwsAuthData struct {
	MapRoles []*RolesAuthMap `yaml:"mapRoles"`
	MapUsers []*UsersAuthMap `yaml:"mapUsers"`
}

// SetMapRoles sets the MapRoles element
func (m *AwsAuthData) SetMapRoles(authMap []*RolesAuthMap) {
	m.MapRoles = authMap
}

// SetMapUsers sets the MapUsers element
func (m *AwsAuthData) SetMapUsers(authMap []*UsersAuthMap) {
	m.MapUsers = authMap
}

// RemoveArguments are the arguments for removing a mapRole or mapUsers
type RemoveArguments struct {
	KubeconfigPath string
	MapRoles       bool
	MapUsers       bool
	Username       string
	RoleARN        string
	UserARN        string
	Groups         []string
}

func (args *RemoveArguments) Validate() {
	if args.RoleARN == "" && args.MapRoles {
		log.Fatal("error: --rolearn not provided")
	}

	if args.UserARN == "" && args.MapUsers {
		log.Fatal("error: --userarn not provided")
	}

	if args.MapUsers && args.MapRoles {
		log.Fatal("error: --mapusers and --maproles are mutually exclusive")
	}

	if !args.MapUsers && !args.MapRoles {
		log.Fatal("error: must select --mapusers or --maproles")
	}
}

// UpsertArguments are the arguments for upserting a mapRole or mapUsers
type UpsertArguments struct {
	KubeconfigPath string
	MapRoles       bool
	MapUsers       bool
	Username       string
	RoleARN        string
	UserARN        string
	Groups         []string
}

func (args *UpsertArguments) Validate() {
	if args.RoleARN == "" && args.MapRoles {
		log.Fatal("error: --rolearn not provided")
	}

	if args.UserARN == "" && args.MapUsers {
		log.Fatal("error: --userarn not provided")
	}

	if len(args.Groups) == 0 {
		log.Fatal("error: --groups not provided")
	}

	if args.Username == "" {
		log.Fatal("error: --username not provided")
	}

	if args.MapUsers && args.MapRoles {
		log.Fatal("error: --mapusers and --maproles are mutually exclusive")
	}

	if !args.MapUsers && !args.MapRoles {
		log.Fatal("error: must select --mapusers or --maproles")
	}
}

// RolesAuthMap is the basic structure of a mapRoles authentication object
type RolesAuthMap struct {
	RoleARN  string   `yaml:"rolearn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

// UsersAuthMap is the basic structure of a mapUsers authentication object
type UsersAuthMap struct {
	UserARN  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

// NewRolesAuthMap returns a new NewRolesAuthMap
func NewRolesAuthMap(rolearn, username string, groups []string) *RolesAuthMap {
	return &RolesAuthMap{
		RoleARN:  rolearn,
		Username: username,
		Groups:   groups,
	}
}

// NewUsersAuthMap returns a new NewUsersAuthMap
func NewUsersAuthMap(userarn, username string, groups []string) *UsersAuthMap {
	return &UsersAuthMap{
		UserARN:  userarn,
		Username: username,
		Groups:   groups,
	}
}

// SetUserARN sets the RoleARN value
func (s *UsersAuthMap) SetUserARN(v string) *UsersAuthMap {
	s.UserARN = v
	return s
}

// SetUsername sets the Username value
func (s *UsersAuthMap) SetUsername(v string) *UsersAuthMap {
	s.Username = v
	return s
}

// SetGroups sets the Groups value
func (s *UsersAuthMap) SetGroups(g []string) *UsersAuthMap {
	s.Groups = g
	return s
}

// SetRoleARN sets the RoleARN value
func (s *RolesAuthMap) SetRoleARN(v string) *RolesAuthMap {
	s.RoleARN = v
	return s
}

// SetUsername sets the Username value
func (s *RolesAuthMap) SetUsername(v string) *RolesAuthMap {
	s.Username = v
	return s
}

// SetGroups sets the Groups value
func (s *RolesAuthMap) SetGroups(g []string) *RolesAuthMap {
	s.Groups = g
	return s
}
