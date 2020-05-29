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
	yaml "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	AwsAuthNamespace = "kube-system"
	AwsAuthName      = "aws-auth"
)

// ReadAuthMap reads the aws-auth config map and returns an AwsAuthData and the actualy ConfigMap objects
func ReadAuthMap(k kubernetes.Interface) (AwsAuthData, *v1.ConfigMap, error) {
	var authData AwsAuthData

	cm, err := k.CoreV1().ConfigMaps(AwsAuthNamespace).Get(AwsAuthName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			cm, err = CreateAuthMap(k)
			if err != nil {
				return authData, cm, err
			}
		} else {
			return authData, cm, err
		}
	}

	err = yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &authData.MapRoles)
	if err != nil {
		return authData, cm, err
	}

	err = yaml.Unmarshal([]byte(cm.Data["mapUsers"]), &authData.MapUsers)
	if err != nil {
		return authData, cm, err
	}

	return authData, cm, nil
}

func CreateAuthMap(k kubernetes.Interface) (*v1.ConfigMap, error) {
	configMapObject := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: "kube-system",
		},
	}
	configMap, err := k.CoreV1().ConfigMaps("kube-system").Create(configMapObject)
	if err != nil {
		return configMap, err
	}
	return configMap, nil
}

// UpdateAuthMap updates a given ConfigMap
func UpdateAuthMap(k kubernetes.Interface, authData AwsAuthData, cm *v1.ConfigMap) error {

	mapRoles, err := yaml.Marshal(authData.MapRoles)
	if err != nil {
		return err
	}

	mapUsers, err := yaml.Marshal(authData.MapUsers)
	if err != nil {
		return err
	}

	cm.Data = map[string]string{
		"mapRoles": string(mapRoles),
		"mapUsers": string(mapUsers),
	}

	cm, err = k.CoreV1().ConfigMaps(AwsAuthNamespace).Update(cm)
	if err != nil {
		return err
	}

	return nil
}
