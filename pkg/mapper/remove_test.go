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
	"testing"
	"time"

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
)

func TestMapper_Remove(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Remove(&MapperArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "admin",
		Groups:   []string{"system:masters"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}

func TestMapper_RemoveByARN(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Remove(&MapperArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}

func TestMapper_RemoveNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-2",
	})
	g.Expect(err).To(gomega.HaveOccurred())

	err = mapper.Remove(&MapperArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-2",
	})
	g.Expect(err).To(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveNotFoundForce(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		Force:    true,
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-2",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Remove(&MapperArguments{
		Force:    true,
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-2",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveByUsername(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.RemoveByUsername(&MapperArguments{
		Username: "system:node:{{EC2PrivateDNSName}}",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))

	err = mapper.RemoveByUsername(&MapperArguments{
		Username: "admin",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err = ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}

func TestMapper_RemoveByUsernameWithRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.RemoveByUsername(&MapperArguments{
		Username:      "system:node:{{EC2PrivateDNSName}}",
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))

	err = mapper.RemoveByUsername(&MapperArguments{
		Username:      "admin",
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err = ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}

func TestMapper_RemoveByUsernameNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.RemoveByUsername(&MapperArguments{
		Username: "",
	})
	g.Expect(err).To(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveByUsernameNotFoundForce(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.RemoveByUsername(&MapperArguments{
		Username: "",
		Force:    true,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveRole_PartialMatchARNOnly(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	// The ARN exists but with a different username — should fail to remove (no exact match)
	err := mapper.Remove(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "wrong-username",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).To(gomega.HaveOccurred())

	// Configmap should be unchanged
	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
}

func TestMapper_RemoveUser_PartialMatchARNOnly(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	// The ARN exists but with a different username — should fail to remove (no exact match)
	err := mapper.Remove(&MapperArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "wrong-username",
		Groups:   []string{"system:masters"},
	})
	g.Expect(err).To(gomega.HaveOccurred())

	// Configmap should be unchanged
	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveRole_GroupsMismatchNotOverriddenByUsername(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	// ARN and username match, but groups do NOT match — should fail to remove
	err := mapper.Remove(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"wrong-group"},
	})
	g.Expect(err).To(gomega.HaveOccurred())

	// Configmap should be unchanged
	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
}

func TestMapper_RemoveUser_GroupsMismatchNotOverriddenByUsername(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	// ARN and username match, but groups do NOT match — should fail to remove
	err := mapper.Remove(&MapperArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "admin",
		Groups:   []string{"wrong-group"},
	})
	g.Expect(err).To(gomega.HaveOccurred())

	// Configmap should be unchanged
	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveRole_NotFound_ErrorMessage(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/nonexistent",
	})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("rolemap"))
}

func TestMapper_RemoveUser_NotFound_ErrorMessage(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/nonexistent",
	})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("usermap"))
}

func TestMapper_RemoveWithRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Remove(&MapperArguments{
		MapRoles:      true,
		RoleARN:       "arn:aws:iam::00000000000:role/node-1",
		Username:      "system:node:{{EC2PrivateDNSName}}",
		Groups:        []string{"system:bootstrappers", "system:nodes"},
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Remove(&MapperArguments{
		MapUsers:      true,
		UserARN:       "arn:aws:iam::00000000000:user/user-1",
		Username:      "admin",
		Groups:        []string{"system:masters"},
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}
