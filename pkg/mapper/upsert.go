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
	"os"
)

// Upsert update or inserts by rolearn
func (b *Bootstrapper) Upsert(args *UpsertArguments) error {
	args.validate()
	var resource = NewAuthMap(args.RoleARN, args.Username, args.Groups)

	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	if args.MapRoles {
		upsertMapRole(&authData, resource)
	}

	if args.MapUsers {
		upsertMapUser(&authData, resource)
	}

	// Update the config map and return an AuthMap
	err = UpdateAuthMap(b.KubernetesClient, authData, configMap)
	if err != nil {
		return err
	}

	return nil
}

func upsertMapRole(authMap *AwsAuthData, resource *AuthMap) {
	var match bool
	for _, existing := range authMap.MapRoles {
		if existing.RoleARN == resource.RoleARN {
			match = true
			existing.SetGroups(resource.Groups)
			existing.SetUsername(resource.Username)
		}
	}
	if !match {
		authMap.AddUniqueMapRole(resource)
	}
}

func upsertMapUser(authMap *AwsAuthData, resource *AuthMap) {
	var match bool
	for _, existing := range authMap.MapRoles {
		if existing.RoleARN == resource.RoleARN {
			match = true
			existing.SetGroups(resource.Groups)
			existing.SetUsername(resource.Username)
		}
	}
	if !match {
		authMap.AddUniqueMapUser(resource)
	}
}

func (args *UpsertArguments) validate() {
	if args.RoleARN == "" {
		fmt.Println("error: --rolearn not provided")
		os.Exit(1)
	}

	if len(args.Groups) == 0 {
		fmt.Println("error: --groups not provided")
		os.Exit(1)
	}

	if args.Username == "" {
		fmt.Println("error: --username not provided")
		os.Exit(1)
	}

	if args.MapUsers && args.MapRoles {
		fmt.Println("error: --mapusers and --maproles are mutually exclusive")
		os.Exit(1)
	}

	if !args.MapUsers && !args.MapRoles {
		fmt.Println("error: must select --mapusers or --maproles")
		os.Exit(1)
	}
}
