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

func TestMapper_UpsertInsert(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Upsert(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-2",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&MapperArguments{
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

	err := mapper.Upsert(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&MapperArguments{
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

	err := mapper.Upsert(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&MapperArguments{
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

func TestMapper_UpsertWithCreate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)

	err := mapper.Upsert(&MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&MapperArguments{
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

func TestMapper_UpsertMultipleInsert(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	role2 := &RolesAuthMap{
		RoleARN:  "arn:aws:iam::00000000000:role/node-2",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	}

	role3 := &RolesAuthMap{
		RoleARN:  "arn:aws:iam::00000000000:role/node-3",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	}

	err := mapper.UpsertMultiple([]*RolesAuthMap{role2, role3}, []*UsersAuthMap{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(3))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))

	mapUser2 := &UsersAuthMap{
		UserARN:  "arn:aws:iam::00000000000:user/user-2",
		Username: "admin",
		Groups:   []string{"system:masters"},
	}

	err = mapper.UpsertMultiple([]*RolesAuthMap{}, []*UsersAuthMap{mapUser2})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err = ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(3))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(2))
}

func TestMapper_UpsertMultipleUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	role1 := &RolesAuthMap{
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	}

	err := mapper.UpsertMultiple([]*RolesAuthMap{role1}, []*UsersAuthMap{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal("arn:aws:iam::00000000000:role/node-1"))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal("arn:aws:iam::00000000000:user/user-1"))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("admin"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.Equal([]string{"system:masters"}))

	mapUser1 := &UsersAuthMap{
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	}

	err = mapper.UpsertMultiple([]*RolesAuthMap{}, []*UsersAuthMap{mapUser1})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err = ReadAuthMap(client)
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

func TestMapper_UpsertMultipleNotNeeded(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	role1 := &RolesAuthMap{
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups:   []string{"system:bootstrappers", "system:nodes"},
	}

	mapUser1 := &UsersAuthMap{
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "admin",
		Groups:   []string{"system:masters"},
	}

	err := mapper.UpsertMultiple([]*RolesAuthMap{role1}, []*UsersAuthMap{mapUser1})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_UpsertMultipleWithCreate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)

	role1 := &RolesAuthMap{
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	}

	mapUser1 := &UsersAuthMap{
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "this:is:a:test",
		Groups:   []string{"system:some-role"},
	}

	err := mapper.UpsertMultiple([]*RolesAuthMap{role1}, []*UsersAuthMap{mapUser1})
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

func TestMapper_UpsertEmptyGroups(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)

	role1 := &RolesAuthMap{
		RoleARN:  "arn:aws:iam::00000000000:role/node-1",
		Username: "this:is:a:test",
		Groups:   []string{},
	}

	mapUser1 := &UsersAuthMap{
		UserARN:  "arn:aws:iam::00000000000:user/user-1",
		Username: "this:is:a:test",
		Groups:   []string{},
	}

	err := mapper.UpsertMultiple([]*RolesAuthMap{role1}, []*UsersAuthMap{mapUser1})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal("arn:aws:iam::00000000000:role/node-1"))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.BeNil())
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal("arn:aws:iam::00000000000:user/user-1"))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.BeNil())
}

func TestMapper_UpsertWithRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Upsert(&MapperArguments{
		MapRoles:      true,
		RoleARN:       "arn:aws:iam::00000000000:role/node-2",
		Username:      "system:node:{{EC2PrivateDNSName}}",
		Groups:        []string{"system:bootstrappers", "system:nodes"},
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&MapperArguments{
		MapUsers:      true,
		UserARN:       "arn:aws:iam::00000000000:user/user-2",
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
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(2))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(2))
}

func TestUpsertWithRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockConfigMap(client)

	err := mapper.Upsert(&MapperArguments{
		MapRoles:      true,
		RoleARN:       "arn:aws:iam::00000000000:role/node-1",
		Username:      "this:is:a:test",
		Groups:        []string{"system:some-role"},
		WithRetries:   true,
		MaxRetryCount: 12,
		MaxRetryTime:  1 * time.Millisecond,
		MinRetryTime:  1 * time.Millisecond,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&MapperArguments{
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
