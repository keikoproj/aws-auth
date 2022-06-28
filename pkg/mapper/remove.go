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
	"errors"
	"fmt"
	"log"
	"reflect"
)

// Remove removes by match of provided arguments
func (b *AuthMapper) Remove(args *MapperArguments) error {
	args.Validate()

	if args.WithRetries {
		_, err := WithRetry(func() (interface{}, error) {
			return nil, b.removeAuth(args)
		}, args)
		return err
	}
	return b.removeAuth(args)
}

// RemoveByUsername removes all map roles and map users that match provided username
func (b *AuthMapper) RemoveByUsername(args *MapperArguments) error {
	args.IsGlobal = true
	args.Validate()
	if args.WithRetries {
		_, err := WithRetry(func() (interface{}, error) {
			return nil, b.removeAuthByUser(args)
		}, args)
		return err
	}
	return b.removeAuthByUser(args)
}

func (b *AuthMapper) removeAuthByUser(args *MapperArguments) error {
	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}
	removed := false

	var newRolesAuthMap []*RolesAuthMap

	for _, mapRole := range authData.MapRoles {
		// Add all other members except the matched
		if args.Username != mapRole.Username {
			newRolesAuthMap = append(newRolesAuthMap, mapRole)
		} else {
			removed = true
		}
	}

	var newUsersAuthMap []*UsersAuthMap

	for _, mapUser := range authData.MapUsers {
		// Add all other members except the matched
		if args.Username != mapUser.Username {
			newUsersAuthMap = append(newUsersAuthMap, mapUser)
		} else {
			removed = true
		}
	}

	if !removed {
		msg := fmt.Sprintf("failed to remove based on username %v, found zero matches\n", args.Username)
		log.Printf(msg)
		if args.Force {
			return nil
		}
		return errors.New(msg)
	}

	authData.SetMapRoles(newRolesAuthMap)
	authData.SetMapUsers(newUsersAuthMap)

	return UpdateAuthMap(b.KubernetesClient, authData, configMap)
}

func (b *AuthMapper) removeAuth(args *MapperArguments) error {
	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	if args.MapRoles {
		var rolesResource = NewRolesAuthMap(args.RoleARN, args.Username, args.Groups)
		newMap, ok := removeRole(authData.MapRoles, rolesResource)

		if !ok {
			log.Printf("failed to remove %v, could not find exact match\n", rolesResource.RoleARN)
			if args.Force {
				return nil
			}
			return errors.New("could not find rolemap")
		}
		log.Printf("removed %v from aws-auth\n", rolesResource.RoleARN)
		authData.SetMapRoles(newMap)
	}

	if args.MapUsers {
		var usersResource = NewUsersAuthMap(args.UserARN, args.Username, args.Groups)
		newMap, ok := removeUser(authData.MapUsers, usersResource)

		if !ok {
			log.Printf("failed to remove %v, could not find exact match\n", usersResource.UserARN)
			if args.Force {
				return nil
			}
			return errors.New("could not find rolemap")
		}
		log.Printf("removed %v from aws-auth\n", usersResource.UserARN)
		authData.SetMapUsers(newMap)
	}

	return UpdateAuthMap(b.KubernetesClient, authData, configMap)
}

func removeRole(authMaps []*RolesAuthMap, targetMap *RolesAuthMap) ([]*RolesAuthMap, bool) {
	var newMap []*RolesAuthMap
	var match bool
	var removed bool

	for _, existingMap := range authMaps {
		match = false
		if existingMap.RoleARN == targetMap.RoleARN {
			match = true
			if len(targetMap.Groups) != 0 {
				if reflect.DeepEqual(existingMap.Groups, targetMap.Groups) {
					match = true
				} else {
					match = false
				}
			}
			if targetMap.Username != "" {
				if existingMap.Username == targetMap.Username {
					match = true
				} else {
					match = false
				}
			}
		}
		if match {
			removed = true
		} else {
			newMap = append(newMap, existingMap)
		}
	}
	return newMap, removed
}

func removeUser(authMaps []*UsersAuthMap, targetMap *UsersAuthMap) ([]*UsersAuthMap, bool) {
	var newMap []*UsersAuthMap
	var match bool
	var removed bool

	for _, existingMap := range authMaps {
		match = false
		if existingMap.UserARN == targetMap.UserARN {
			match = true
			if len(targetMap.Groups) != 0 {
				if reflect.DeepEqual(existingMap.Groups, targetMap.Groups) {
					match = true
				} else {
					match = false
				}
			}
			if targetMap.Username != "" {
				if existingMap.Username == targetMap.Username {
					match = true
				} else {
					match = false
				}
			}
		}
		if match {
			removed = true
		} else {
			newMap = append(newMap, existingMap)
		}
	}
	return newMap, removed
}
