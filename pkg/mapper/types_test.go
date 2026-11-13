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
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNew_LoggingDisabled(t *testing.T) {
	g := gomega.NewWithT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, false)
	g.Expect(mapper).NotTo(gomega.BeNil())
	g.Expect(mapper.KubernetesClient).To(gomega.Equal(client))
	g.Expect(mapper.Logger).NotTo(gomega.BeNil())
	// Logger should silently discard output
	mapper.Logger.Print("test")
}

func TestNew_LoggingEnabled(t *testing.T) {
	g := gomega.NewWithT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	g.Expect(mapper).NotTo(gomega.BeNil())
	g.Expect(mapper.Logger).NotTo(gomega.BeNil())
}

func TestNew_IndependentLoggers(t *testing.T) {
	g := gomega.NewWithT(t)
	client := fake.NewSimpleClientset()

	// Create one mapper with logging disabled, one with a custom buffer
	silent := New(client, false)
	var buf bytes.Buffer
	loud := &AuthMapper{
		KubernetesClient: client,
		Logger:           log.New(&buf, "", 0),
	}

	// Writing to the silent logger must not affect the loud logger
	silent.Logger.Print("should be discarded")
	loud.Logger.Print("should appear")

	g.Expect(buf.String()).To(gomega.ContainSubstring("should appear"))
	g.Expect(buf.String()).NotTo(gomega.ContainSubstring("should be discarded"))
}

func TestNewRolesAuthMap(t *testing.T) {
	g := gomega.NewWithT(t)
	r := NewRolesAuthMap("arn:aws:iam::123:role/foo", "myuser", []string{"group1", "group2"})
	g.Expect(r.RoleARN).To(gomega.Equal("arn:aws:iam::123:role/foo"))
	g.Expect(r.Username).To(gomega.Equal("myuser"))
	g.Expect(r.Groups).To(gomega.Equal([]string{"group1", "group2"}))
}

func TestNewUsersAuthMap(t *testing.T) {
	g := gomega.NewWithT(t)
	u := NewUsersAuthMap("arn:aws:iam::123:user/bar", "myuser", []string{"group1"})
	g.Expect(u.UserARN).To(gomega.Equal("arn:aws:iam::123:user/bar"))
	g.Expect(u.Username).To(gomega.Equal("myuser"))
	g.Expect(u.Groups).To(gomega.Equal([]string{"group1"}))
}

func TestRolesAuthMap_String(t *testing.T) {
	g := gomega.NewWithT(t)
	r := NewRolesAuthMap("arn:aws:iam::123:role/foo", "myuser", []string{"group1", "group2"})
	s := r.String()
	g.Expect(s).To(gomega.ContainSubstring("rolearn: arn:aws:iam::123:role/foo"))
	g.Expect(s).To(gomega.ContainSubstring("username: myuser"))
	g.Expect(s).To(gomega.ContainSubstring("groups:"))
	g.Expect(s).To(gomega.ContainSubstring("- group1"))
	g.Expect(s).To(gomega.ContainSubstring("- group2"))
}

func TestUsersAuthMap_String(t *testing.T) {
	g := gomega.NewWithT(t)
	u := NewUsersAuthMap("arn:aws:iam::123:user/bar", "myuser", []string{"system:masters"})
	s := u.String()
	g.Expect(s).To(gomega.ContainSubstring("userarn: arn:aws:iam::123:user/bar"))
	g.Expect(s).To(gomega.ContainSubstring("username: myuser"))
	g.Expect(s).To(gomega.ContainSubstring("groups:"))
	g.Expect(s).To(gomega.ContainSubstring("system:masters"))
}

func TestSetters_RolesAuthMap(t *testing.T) {
	g := gomega.NewWithT(t)
	r := NewRolesAuthMap("arn:aws:iam::123:role/foo", "original", []string{"g1"})

	r.SetUsername("updated")
	g.Expect(r.Username).To(gomega.Equal("updated"))

	r.SetGroups([]string{"g2", "g3"})
	g.Expect(r.Groups).To(gomega.Equal([]string{"g2", "g3"}))

	r.AppendGroups([]string{"g4"})
	g.Expect(r.Groups).To(gomega.Equal([]string{"g2", "g3", "g4"}))
}

func TestSetters_UsersAuthMap(t *testing.T) {
	g := gomega.NewWithT(t)
	u := NewUsersAuthMap("arn:aws:iam::123:user/bar", "original", []string{"g1"})

	u.SetUsername("updated")
	g.Expect(u.Username).To(gomega.Equal("updated"))

	u.SetGroups([]string{"g2", "g3"})
	g.Expect(u.Groups).To(gomega.Equal([]string{"g2", "g3"}))

	u.AppendGroups([]string{"g4"})
	g.Expect(u.Groups).To(gomega.Equal([]string{"g2", "g3", "g4"}))
}

func TestWithRetry_SucceedsFirstAttempt(t *testing.T) {
	g := gomega.NewWithT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	calls := 0
	fn := RetriableFunction(func() (interface{}, error) {
		calls++
		return "ok", nil
	})
	out, err := mapper.WithRetry(fn, &MapperArguments{
		MaxRetryCount: 3,
		MinRetryTime:  1 * time.Millisecond,
		MaxRetryTime:  2 * time.Millisecond,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(out).To(gomega.Equal("ok"))
	g.Expect(calls).To(gomega.Equal(1))
}

func TestWithRetry_ExhaustsRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	calls := 0
	fn := RetriableFunction(func() (interface{}, error) {
		calls++
		return nil, gomega.StopTrying("always fails")
	})
	_, err := mapper.WithRetry(fn, &MapperArguments{
		MaxRetryCount: 2,
		MinRetryTime:  1 * time.Millisecond,
		MaxRetryTime:  2 * time.Millisecond,
	})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("waiter timed out"))
	g.Expect(calls).To(gomega.Equal(2))
}

func TestValidate_InvalidRetryMaxCount(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{WithRetries: true, MaxRetryCount: 0, MapRoles: true, RoleARN: "arn"}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("--retry-max-count"))
}

func TestValidate_MissingRoleARN(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{MapRoles: true, RoleARN: ""}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("--rolearn not provided"))
}

func TestValidate_MissingUserARN(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{MapUsers: true, UserARN: ""}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("--userarn not provided"))
}

func TestValidate_MutuallyExclusiveFlags(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{MapUsers: true, MapRoles: true, UserARN: "arn", RoleARN: "arn"}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("mutually exclusive"))
}

func TestValidate_MissingUsername(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{OperationType: OperationUpsert, MapRoles: true, RoleARN: "arn"}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("--username not provided"))
}

func TestValidate_InvalidFormat(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{OperationType: OperationGet, Format: "json", IsGlobal: true}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("--format"))
}

func TestValidate_MissingMapSelection(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{MapUsers: false, MapRoles: false, IsGlobal: false}
	err := args.Validate()
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("must select"))
}

func TestValidate_Success(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{
		OperationType: OperationUpsert,
		MapRoles:      true,
		RoleARN:       "arn:aws:iam::123:role/foo",
		Username:      "myuser",
	}
	err := args.Validate()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(args.UpdateUsername).NotTo(gomega.BeNil())
	g.Expect(*args.UpdateUsername).To(gomega.BeTrue())
}

func TestValidate_GlobalSkipsMapSelection(t *testing.T) {
	g := gomega.NewWithT(t)
	args := &MapperArguments{IsGlobal: true}
	err := args.Validate()
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
