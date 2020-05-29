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
	"testing"

	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func create_MockConfigMap(client kubernetes.Interface) {
	role := NewRolesAuthMap("arn:aws:iam::00000000000:role/node-1",
		"system:node:{{EC2PrivateDNSName}}",
		[]string{"system:bootstrappers", "system:nodes"})

	user := NewUsersAuthMap("arn:aws:iam::00000000000:user/user-1",
		"admin",
		[]string{"system:masters"})

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AwsAuthName,
			Namespace: AwsAuthNamespace,
		},
		Data: map[string]string{
			"mapRoles": role.String(),
			"mapUsers": user.String(),
		},
	}
	_, err := client.CoreV1().ConfigMaps(AwsAuthNamespace).Create(configMap)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

}

func TestConfigMaps_Update(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	create_MockConfigMap(client)

	auth, cm, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	role := NewRolesAuthMap("arn:aws:iam::00000000000:role/node-2",
		"system:node:{{EC2PrivateDNSName}}",
		[]string{"system:bootstrappers", "system:nodes"})
	user := NewUsersAuthMap("arn:aws:iam::00000000000:user/user-2",
		"ops-user",
		[]string{"system:masters"})

	auth.MapRoles = append(auth.MapRoles, role)
	auth.MapUsers = append(auth.MapUsers, user)

	err = UpdateAuthMap(client, auth, cm)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, cm, err = ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	fmt.Println(auth.MapRoles[0])
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(2))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(2))
}

func TestConfigMaps_Read(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	create_MockConfigMap(client)

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal("arn:aws:iam::00000000000:role/node-1"))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("system:node:{{EC2PrivateDNSName}}"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.Equal([]string{"system:bootstrappers", "system:nodes"}))

	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal("arn:aws:iam::00000000000:user/user-1"))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("admin"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.Equal([]string{"system:masters"}))
}

func TestConfigMaps_Create(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}
