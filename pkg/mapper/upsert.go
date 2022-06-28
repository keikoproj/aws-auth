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
func (b *AuthMapper) Upsert(args *MapperArguments) error {
	args.Validate()

	if args.WithRetries {
		_, err := WithRetry(func() (interface{}, error) {
			return nil, b.upsertAuth(args)
		}, args)
		return err
	}

	return b.upsertAuth(args)
}

/**
 *  UpsertMultiple upserts list of mapRoles and mapUsers into the configmap
 *  if no changes are required based on new entries, configmap doesn't get updated
 */
func (b *AuthMapper) UpsertMultiple(newMapRoles []*RolesAuthMap, newMapUsers []*UsersAuthMap) error {
	updated := false
	mapRoles := []*RolesAuthMap{}
	mapUsers := []*UsersAuthMap{}

	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	// Insert all new mapRole entries
	for _, newMember := range newMapRoles {
		found := false
		for _, existing := range authData.MapRoles {

			if existing.RoleARN == newMember.RoleARN {
				found = true
			}
		}

		if !found {
			updated = true
			mapRoles = append(mapRoles, newMember)
		}
	}

	// Upsert existing mapRoles
	for _, existing := range authData.MapRoles {

		for _, newMember := range newMapRoles {

			if existing.RoleARN != newMember.RoleARN {
				continue
			}

			if !reflect.DeepEqual(existing.Groups, newMember.Groups) {
				existing.SetGroups(newMember.Groups)
				updated = true
			}

			if existing.Username != newMember.Username {
				existing.SetUsername(newMember.Username)
				updated = true
			}

		}

		mapRoles = append(mapRoles, existing)
	}

	// Insert all new mapUser entries
	for _, newMember := range newMapUsers {
		found := false
		for _, existing := range authData.MapUsers {

			if existing.UserARN == newMember.UserARN {
				found = true
			}
		}

		if !found {
			updated = true
			mapUsers = append(mapUsers, newMember)
		}
	}

	// Upsert existing mapUsers
	for _, existing := range authData.MapUsers {

		for _, newMember := range newMapUsers {

			if existing.UserARN != newMember.UserARN {
				continue
			}

			if !reflect.DeepEqual(existing.Groups, newMember.Groups) {
				existing.SetGroups(newMember.Groups)
				updated = true
			}

			if existing.Username != newMember.Username {
				existing.SetUsername(newMember.Username)
				updated = true
			}

		}

		mapUsers = append(mapUsers, existing)
	}

	if !updated {
		log.Printf("found zero changes to update, configmap is not changed \n")
		return nil
	}

	authData.SetMapRoles(mapRoles)
	authData.SetMapUsers((mapUsers))

	// Update the config map and return an AuthMap
	err = UpdateAuthMap(b.KubernetesClient, authData, configMap)
	if err != nil {
		return err
	}

	return nil
}

func (b *AuthMapper) upsertAuth(args *MapperArguments) error {
	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	opts := &UpsertOptions{
		Append:         args.Append,
		UpdateUsername: *args.UpdateUsername,
	}

	if args.MapRoles {
		var roleResource = NewRolesAuthMap(args.RoleARN, args.Username, args.Groups)

		newMap, ok := upsertRole(authData.MapRoles, roleResource, opts)
		if ok {
			log.Printf("role %v has been updated\n", roleResource.RoleARN)
		} else {
			log.Printf("no updates needed to %v\n", roleResource.RoleARN)
		}
		authData.SetMapRoles(newMap)
	}

	if args.MapUsers {
		var userResource = NewUsersAuthMap(args.UserARN, args.Username, args.Groups)

		newMap, ok := upsertUser(authData.MapUsers, userResource, opts)
		if ok {
			log.Printf("role %v has been updated\n", userResource.UserARN)
		} else {
			log.Printf("no updates needed to %v\n", userResource.UserARN)
		}
		authData.SetMapUsers(newMap)
	}

	return UpdateAuthMap(b.KubernetesClient, authData, configMap)
}

func upsertRole(authMaps []*RolesAuthMap, resource *RolesAuthMap, opts *UpsertOptions) ([]*RolesAuthMap, bool) {
	var match bool
	var updated bool
	for _, existing := range authMaps {
		// Update
		if existing.RoleARN == resource.RoleARN {
			match = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				if opts.Append {
					existing.AppendGroups(resource.Groups)
				} else {
					existing.SetGroups(resource.Groups)
				}
				updated = true
			}
			if existing.Username != resource.Username {
				if opts.UpdateUsername {
					existing.SetUsername(resource.Username)
					updated = true
				}
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

func upsertUser(authMaps []*UsersAuthMap, resource *UsersAuthMap, opts *UpsertOptions) ([]*UsersAuthMap, bool) {
	var match bool
	var updated bool
	for _, existing := range authMaps {
		// Update
		if existing.UserARN == resource.UserARN {
			match = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				if opts.Append {
					existing.AppendGroups(resource.Groups)
				} else {
					existing.SetGroups(resource.Groups)
				}
				updated = true
			}
			if existing.Username != resource.Username {
				if opts.UpdateUsername {
					existing.SetUsername(resource.Username)
					updated = true
				}
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
