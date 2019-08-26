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
	"log"
	"reflect"
)

// Remove removes by match of provided arguments
func (b *AuthMapper) Remove(args *RemoveArguments) error {
	args.Validate()
	var resource = NewAuthMap(args.RoleARN, args.Username, args.Groups)

	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	if args.MapRoles {
		newMap := removeRole(authData.MapRoles, resource)

		if len(authData.MapRoles) == len(newMap) {
			log.Printf("failed to remove %v, could not find exact match\n", resource.RoleARN)
		} else {
			log.Printf("removed %v from aws-auth\n", resource.RoleARN)
		}
		authData.SetMapRoles(newMap)

	}

	if args.MapUsers {
		newMap := removeRole(authData.MapUsers, resource)

		if len(authData.MapUsers) == len(newMap) {
			log.Printf("failed to remove %v, could not find exact match\n", resource.RoleARN)
		} else {
			log.Printf("removed %v from aws-auth\n", resource.RoleARN)
		}
		authData.SetMapUsers(newMap)

	}

	// Update the config map and return an AuthMap
	err = UpdateAuthMap(b.KubernetesClient, authData, configMap)
	if err != nil {
		return err
	}

	return nil
}

func removeRole(authMaps []*AuthMap, targetMap *AuthMap) []*AuthMap {
	var newMap []*AuthMap
	var match bool

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
		if !match {
			newMap = append(newMap, existingMap)
		}
	}
	return newMap
}
