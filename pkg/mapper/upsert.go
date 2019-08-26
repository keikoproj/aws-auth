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

// Upsert update or inserts by rolearn
func (b *AuthMapper) Upsert(args *UpsertArguments) error {
	args.Validate()

	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	if args.MapRoles {
		var roleResource = NewRolesAuthMap(args.RoleARN, args.Username, args.Groups)

		newMap, ok := upsertRole(authData.MapRoles, roleResource)
		if ok {
			log.Printf("role %v has been updated\n", roleResource.RoleARN)
		} else {
			log.Printf("no updates needed to %v\n", roleResource.RoleARN)
		}
		authData.SetMapRoles(newMap)
	}

	if args.MapUsers {
		var userResource = NewUsersAuthMap(args.UserARN, args.Username, args.Groups)

		newMap, ok := upsertUser(authData.MapUsers, userResource)
		if ok {
			log.Printf("role %v has been updated\n", userResource.UserARN)
		} else {
			log.Printf("no updates needed to %v\n", userResource.UserARN)
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

func upsertRole(authMaps []*RolesAuthMap, resource *RolesAuthMap) ([]*RolesAuthMap, bool) {
	var match bool
	var updated bool
	for _, existing := range authMaps {
		// Update
		if existing.RoleARN == resource.RoleARN {
			match = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				existing.SetGroups(resource.Groups)
				updated = true
			}
			if existing.Username != resource.Username {
				existing.SetUsername(resource.Username)
				updated = true
			}
		}
	}

	// Insert
	if !match {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}

func upsertUser(authMaps []*UsersAuthMap, resource *UsersAuthMap) ([]*UsersAuthMap, bool) {
	var match bool
	var updated bool
	for _, existing := range authMaps {
		// Update
		if existing.UserARN == resource.UserARN {
			match = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				existing.SetGroups(resource.Groups)
				updated = true
			}
			if existing.Username != resource.Username {
				existing.SetUsername(resource.Username)
				updated = true
			}
		}
	}

	// Insert
	if !match {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}
