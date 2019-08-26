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
	var resource = NewAuthMap(args.RoleARN, args.Username, args.Groups)

	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	if args.MapRoles {
		newMap, ok := upsertRole(authData.MapRoles, resource)
		//authData.AddUniqueMapRole(resource)

		if ok {
			log.Printf("role %v has been updated\n", resource.RoleARN)
		} else {
			log.Printf("no updates needed to %v\n", resource.RoleARN)
		}
		authData.SetMapRoles(newMap)
	}

	if args.MapUsers {
		newMap, ok := upsertRole(authData.MapUsers, resource)

		//authData.AddUniqueMapRole(resource)

		if ok {
			log.Printf("role %v has been updated\n", resource.RoleARN)
		} else {
			log.Printf("no updates needed to %v\n", resource.RoleARN)
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

func upsertRole(authMaps []*AuthMap, resource *AuthMap) ([]*AuthMap, bool) {
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
