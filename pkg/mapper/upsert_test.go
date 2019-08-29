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

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
)

func TestMapper_UpsertInsert(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Upsert(&UpsertArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-2",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&UpsertArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-2",
		Username: "admin",
		Groups:   []string{"system:masters"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(2))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(2))
}

func TestMapper_UpsertUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Upsert(&UpsertArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&UpsertArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal("arn:aws:iam::00000000000:role/node-1"))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal("arn:aws:iam::00000000000:user/user-1"))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
}

func TestMapper_UpsertNotNeeded(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Upsert(&UpsertArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&UpsertArguments{
		MapUsers: true,
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "admin",
		Groups:   []string{"system:masters"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}
