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

// Remove removes by match of provided arguments
func (b *Bootstrapper) Remove(args *RemoveArguments) error {
	args.validate()
	var resource = NewAuthMap(args.RoleARN, args.Username, args.Groups)

	// Read the config map and return an AuthMap
	authData, configMap, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return err
	}

	if args.MapRoles {
		removeMapRole(&authData, resource)
	}

	if args.MapUsers {
		removeMapUser(&authData, resource)
	}

	// Update the config map and return an AuthMap
	err = UpdateAuthMap(b.KubernetesClient, authData, configMap)
	if err != nil {
		return err
	}

	return nil
}

func removeMapRole(authMap *AwsAuthData, resource *AuthMap) {
	authMap.RemoveMapRole(resource)
}

func removeMapUser(authMap *AwsAuthData, resource *AuthMap) {
	authMap.RemoveMapUser(resource)
}

func (args *RemoveArguments) validate() {
	if args.RoleARN == "" {
		fmt.Println("error: --rolearn not provided")
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
