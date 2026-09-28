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
)

// Remove removes by match of provided arguments
func (b *AuthMapper) Remove(args *MapperArguments) error {
	if err := args.Validate(); err != nil {
		return err
	}

	if args.WithRetries {
		_, err := b.WithRetry(func() (interface{}, error) {
			return nil, b.removeAuth(args)
		}, args)
		return err
	}
	return b.removeAuth(args)
}

// RemoveByUsername removes all map roles and map users that match provided username
func (b *AuthMapper) RemoveByUsername(args *MapperArguments) error {
	args.IsGlobal = true
	if err := args.Validate(); err != nil {
		return err
	}
	if args.WithRetries {
		_, err := b.WithRetry(func() (interface{}, error) {
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
		if args.Force {
			return nil
		}
		return fmt.Errorf("failed to remove based on username %v, found zero matches", args.Username)
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
			if args.Force {
				return nil
			}
			return fmt.Errorf("failed to remove %v, could not find exact match", rolesResource.RoleARN)
		}
		b.Logger.Printf("removed %v from aws-auth\n", rolesResource.RoleARN)
		authData.SetMapRoles(newMap)
	}

	if args.MapUsers {
		var usersResource = NewUsersAuthMap(args.UserARN, args.Username, args.Groups)
		newMap, ok := removeUser(authData.MapUsers, usersResource)

		if !ok {
			if args.Force {
				return nil
			}
			return fmt.Errorf("failed to remove %v, could not find exact match", usersResource.UserARN)
		}
		b.Logger.Printf("removed %v from aws-auth\n", usersResource.UserARN)
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
			if len(targetMap.Groups) != 0 && !reflect.DeepEqual(existingMap.Groups, targetMap.Groups) {
				match = false
			}
			if match && targetMap.Username != "" && existingMap.Username != targetMap.Username {
				match = false
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
			if len(targetMap.Groups) != 0 && !reflect.DeepEqual(existingMap.Groups, targetMap.Groups) {
				match = false
			}
			if match && targetMap.Username != "" && existingMap.Username != targetMap.Username {
				match = false
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
