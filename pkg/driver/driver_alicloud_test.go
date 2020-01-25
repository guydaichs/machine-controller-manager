package driver

import (
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func TestTagsOrdered(t *testing.T) {
	tags := map[string]string{
		"kubernetes.io/cluster/ali-test": "1",
		"kubernetes.io/role/worker":      "1",
		"taga":                           "tagvala",
		"tagb":                           "tagvalb",
		"tagc":                           "tagvalc",
	}
	c := &AlicloudDriver{}
	res, err := c.toInstanceTags(tags)
	if err != nil {
		t.Errorf("toInstanceTags in TestTagsOrdered should not generate error: %v", err)
	}

	expected := []ecs.RunInstancesTag{
		{
			Key:   "kubernetes.io/cluster/ali-test",
			Value: "1",
		},
		{
			Key:   "kubernetes.io/role/worker",
			Value: "1",
		},
		{
			Key:   "taga",
			Value: "tagvala",
		},
		{
			Key:   "tagb",
			Value: "tagvalb",
		},
		{
			Key:   "tagc",
			Value: "tagvalc",
		},
	}
	checkRunInstanceTags("Function TestTagsOrdered: ", t, res, expected)
}

func TestNoClusterTags(t *testing.T) {
	tags := map[string]string{
		"kubernetes.io/role/worker": "1",
		"taga":                      "tagvala",
		"tagb":                      "tagvalb",
		"tagc":                      "tagvalc",
	}
	c := &AlicloudDriver{}
	_, err := c.toInstanceTags(tags)
	if err == nil {
		t.Errorf("toInstanceTags in TestRandomOrderTags should return an error")
	}
}
func TestRandomOrderTags(t *testing.T) {
	tags := map[string]string{
		"taga":                           "tagvala",
		"tagb":                           "tagvalb",
		"kubernetes.io/cluster/ali-test": "1",
		"kubernetes.io/role/worker":      "1",
		"tagc":                           "tagvalc",
	}
	c := &AlicloudDriver{}
	res, err := c.toInstanceTags(tags)
	if err != nil {
		t.Errorf("toInstanceTags in TestRandomOrderTags should not generate error: %v", err)
	}

	expected := []ecs.RunInstancesTag{
		{
			Key:   "kubernetes.io/cluster/ali-test",
			Value: "1",
		},
		{
			Key:   "kubernetes.io/role/worker",
			Value: "1",
		},
		{
			Key:   "taga",
			Value: "tagvala",
		},
		{
			Key:   "tagb",
			Value: "tagvalb",
		},
		{
			Key:   "tagc",
			Value: "tagvalc",
		},
	}
	checkRunInstanceTags("Function TestRandomOrderTags: ", t, res, expected)
}
func TestIDToName(t *testing.T) {
	id := "i-uf69zddmom11ci7est12"

	c := &AlicloudDriver{}
	if "iZuf69zddmom11ci7est12Z" != c.idToName(id) {
		t.Error("idToName() is not working")
	}
}

func TestAliCloudDriverSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Driver AliCloud Suite")
}

var _ = Describe("Driver AliCloud", func() {
	Context("Generate Data Disk Requests", func() {

		It("should generate multiple data disk requests", func() {
			c := &AlicloudDriver{}
			dataDisks := []v1alpha1.AlicloudDataDisk{
				{
					Name: "dd1",
					Category: "cloud_efficiency",
					Description: "this is a disk",
					DeleteWithInstance: true,
					Encrypted: true,
					Size: 100,
				},
				{
					Name: "dd2",
					Category: "cloud_ssd",
					Description: "this is also a disk",
					DeleteWithInstance: false,
					Encrypted: false,
					Size: 50,
				},
			}

			generatedDataDisksRequests := c.generateDataDiskRequests(dataDisks)
			expectedDataDiskRequests := []ecs.RunInstancesDataDisk{
				{
					Size:               "100",
					Category:           "cloud_efficiency",
					Encrypted:          strconv.FormatBool(true),
					DiskName:           "dd1",
					Description:        "this is a disk",
					DeleteWithInstance: strconv.FormatBool(true),
				},
				{
					Size:               "50",
					Category:           "cloud_ssd",
					Encrypted:          strconv.FormatBool(false),
					DiskName:           "dd2",
					Description:        "this is also a disk",
					DeleteWithInstance: strconv.FormatBool(false),
				},
			}

			Expect(generatedDataDisksRequests).To(Equal(expectedDataDiskRequests))
		})
	})
})

// real[2..]'s order is NOT predicted as tags which generated them is a MAP!!!
func checkRunInstanceTags(leadErrMsg string, t *testing.T, real, expected []ecs.RunInstancesTag) {
	if len(real) != len(expected) {
		t.Errorf("%s: %s", leadErrMsg, "count of generated tags is not as expected")
		return
	}

	// index 0 and 1 is static
	if real[0] != expected[0] {
		t.Errorf("%s: tag %s should be at index %d", leadErrMsg, expected[0], 0)
	}
	if real[1] != expected[1] {
		t.Errorf("%s: tag %s should be at index %d", leadErrMsg, expected[1], 1)
	}

found:
	for i := 2; i < len(expected); i++ {
		for j := 2; j < len(real); j++ {
			if expected[i] == real[j] {
				continue found
			}
		}
		t.Errorf("%s: tag %s is not in real tags", leadErrMsg, expected[i])
		return
	}
}
