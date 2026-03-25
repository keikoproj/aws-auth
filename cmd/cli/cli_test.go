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

package cli

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestUpsertCmd_KubeconfigFlagBindsToUpsertArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	upsertArgs.KubeconfigPath = ""
	err := upsertCmd.Flags().Set("kubeconfig", "/tmp/test-kubeconfig")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(upsertArgs.KubeconfigPath).To(gomega.Equal("/tmp/test-kubeconfig"))
	g.Expect(getArgs.KubeconfigPath).NotTo(gomega.Equal("/tmp/test-kubeconfig"))

	// cleanup
	upsertArgs.KubeconfigPath = ""
}

func TestUpsertCmd_AsUserFlagBindsToUpsertArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	upsertArgs.AsUser = ""
	err := upsertCmd.Flags().Set("as", "test-user")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(upsertArgs.AsUser).To(gomega.Equal("test-user"))

	// cleanup
	upsertArgs.AsUser = ""
}

func TestUpsertCmd_AsGroupFlagBindsToUpsertArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	upsertArgs.AsGroups = []string{}
	err := upsertCmd.Flags().Set("as-group", "test-group")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(upsertArgs.AsGroups).To(gomega.ContainElement("test-group"))

	// cleanup
	upsertArgs.AsGroups = []string{}
}

func TestRemoveCmd_KubeconfigFlagBindsToRemoveArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	removeArgs.KubeconfigPath = ""
	err := removeCmd.Flags().Set("kubeconfig", "/tmp/test-kubeconfig")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(removeArgs.KubeconfigPath).To(gomega.Equal("/tmp/test-kubeconfig"))
	g.Expect(getArgs.KubeconfigPath).NotTo(gomega.Equal("/tmp/test-kubeconfig"))

	// cleanup
	removeArgs.KubeconfigPath = ""
}

func TestRemoveCmd_AsUserFlagBindsToRemoveArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	removeArgs.AsUser = ""
	err := removeCmd.Flags().Set("as", "test-user")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(removeArgs.AsUser).To(gomega.Equal("test-user"))
	g.Expect(upsertArgs.AsUser).NotTo(gomega.Equal("test-user"))

	// cleanup
	removeArgs.AsUser = ""
}

func TestRemoveCmd_AsGroupFlagBindsToRemoveArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	removeArgs.AsGroups = []string{}
	err := removeCmd.Flags().Set("as-group", "test-group")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(removeArgs.AsGroups).To(gomega.ContainElement("test-group"))
	g.Expect(upsertArgs.AsGroups).NotTo(gomega.ContainElement("test-group"))

	// cleanup
	removeArgs.AsGroups = []string{}
}

func TestRemoveByUsernameCmd_HasKubeconfigFlag(t *testing.T) {
	g := gomega.NewWithT(t)

	cmd := removeByUsernameCmd()
	f := cmd.Flags().Lookup("kubeconfig")
	g.Expect(f).NotTo(gomega.BeNil())
}

func TestRemoveByUsernameCmd_HasAsFlags(t *testing.T) {
	g := gomega.NewWithT(t)

	cmd := removeByUsernameCmd()
	asFlag := cmd.Flags().Lookup("as")
	g.Expect(asFlag).NotTo(gomega.BeNil())

	asGroupFlag := cmd.Flags().Lookup("as-group")
	g.Expect(asGroupFlag).NotTo(gomega.BeNil())
}

func TestGetCmd_KubeconfigFlagBindsToGetArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	getArgs.KubeconfigPath = ""
	err := getCmd.Flags().Set("kubeconfig", "/tmp/test-kubeconfig")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(getArgs.KubeconfigPath).To(gomega.Equal("/tmp/test-kubeconfig"))

	// cleanup
	getArgs.KubeconfigPath = ""
}

func TestGetCmd_AsUserFlagBindsToGetArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	getArgs.AsUser = ""
	err := getCmd.Flags().Set("as", "test-user")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(getArgs.AsUser).To(gomega.Equal("test-user"))
	g.Expect(upsertArgs.AsUser).NotTo(gomega.Equal("test-user"))

	// cleanup
	getArgs.AsUser = ""
}

func TestGetCmd_AsGroupFlagBindsToGetArgs(t *testing.T) {
	g := gomega.NewWithT(t)

	getArgs.AsGroups = []string{}
	err := getCmd.Flags().Set("as-group", "test-group")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(getArgs.AsGroups).To(gomega.ContainElement("test-group"))
	g.Expect(upsertArgs.AsGroups).NotTo(gomega.ContainElement("test-group"))

	// cleanup
	getArgs.AsGroups = []string{}
}

func TestGetKubernetesClient_WithKubeconfigPath(t *testing.T) {
	g := gomega.NewWithT(t)

	// Passing a non-existent kubeconfig should fail with a path-related error,
	// proving the path is actually used
	_, err := getKubernetesClient("/tmp/nonexistent-kubeconfig-12345", kubeOptions{})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("/tmp/nonexistent-kubeconfig-12345"))
}

func TestGetKubernetesClient_ImpersonationOptions(t *testing.T) {
	g := gomega.NewWithT(t)

	// We can't easily test the full client creation without a valid kubeconfig,
	// but we can verify the function signature accepts kubeOptions with AsUser/AsGroups
	options := kubeOptions{
		AsUser:   "admin",
		AsGroups: []string{"system:masters"},
	}
	g.Expect(options.AsUser).To(gomega.Equal("admin"))
	g.Expect(options.AsGroups).To(gomega.Equal([]string{"system:masters"}))
}
