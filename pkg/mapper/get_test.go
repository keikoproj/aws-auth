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

func TestMapper_Get(t *testing.T) {
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

	data, err := mapper.Get(&MapperArguments{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(data).NotTo(gomega.Equal(AwsAuthData{}))

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(auth).To(gomega.Equal(data))
}

func TestMapper_GetRetry(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)

	data, err := mapper.Get(&MapperArguments{
		WithRetries:   true,
		MinRetryTime:  1 * time.Second,
		MaxRetryTime:  2 * time.Second,
		MaxRetryCount: 5,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(data).To(gomega.Equal(AwsAuthData{}))

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(auth).To(gomega.Equal(data))
}

func TestMapper_GetRetryFail(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := New(client, true)
	create_MockMalformedConfigMap(client)

	data, err := mapper.Get(&MapperArguments{
		WithRetries:   true,
		MinRetryTime:  100 * time.Millisecond,
		MaxRetryTime:  200 * time.Millisecond,
		MaxRetryCount: 5,
	})

	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("waiter timed out"))
	g.Expect(data).To(gomega.Equal(AwsAuthData{}))

	_, _, err = ReadAuthMap(client)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("cannot unmarshal"))

}
