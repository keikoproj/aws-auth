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

// Upsert update or inserts by rolearn
func (b *AuthMapper) Get(args *MapperArguments) (AwsAuthData, error) {
	args.IsGlobal = true
	args.Validate()

	if args.WithRetries {
		out, err := WithRetry(func() (interface{}, error) {
			return b.getAuth()
		}, args)
		if err != nil {
			return AwsAuthData{}, err
		}
		return out.(AwsAuthData), nil
	}

	return b.getAuth()
}

func (b *AuthMapper) getAuth() (AwsAuthData, error) {

	// Read the config map and return an AuthMap
	authData, _, err := ReadAuthMap(b.KubernetesClient)
	if err != nil {
		return AwsAuthData{}, err
	}

	return authData, nil
}
