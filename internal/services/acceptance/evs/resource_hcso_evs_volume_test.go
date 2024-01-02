package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func getVolumeResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.BlockStorageV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Huawei Cloud Stack Online block storage v2 client: %s", err)
	}
	return cloudvolumes.Get(c, state.Primary.ID).Extract()
}

func TestAccEvsVolume_basic(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_evs_volume.test"
	resourceName1 := "hcso_evs_volume.test.0"
	resourceName2 := "hcso_evs_volume.test.1"
	resourceName3 := "hcso_evs_volume.test.2"
	resourceName4 := "hcso_evs_volume.test.3"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(4),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.hcso_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"Created by acc test script."),
					resource.TestCheckResourceAttr(resourceName1, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName1, "size", "100"),
					resource.TestCheckResourceAttr(resourceName1, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName1, "tags.key", "value"),
					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_vbd_normal_volume"),
					resource.TestCheckResourceAttr(resourceName1, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName1, "multiattach", "false"),

					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_vbd_share_volume"),
					resource.TestCheckResourceAttr(resourceName2, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName2, "multiattach", "true"),

					resource.TestCheckResourceAttr(resourceName3, "name", rName+"_scsi_normal_volume"),
					resource.TestCheckResourceAttr(resourceName3, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName3, "multiattach", "false"),

					resource.TestCheckResourceAttr(resourceName4, "name", rName+"_scsi_share_volume"),
					resource.TestCheckResourceAttr(resourceName4, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName4, "multiattach", "true"),
				),
			},
			{
				Config: testAccEvsVolume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(4),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.hcso_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"Updated by acc test script."),
					resource.TestCheckResourceAttr(resourceName1, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName1, "size", "200"),
					resource.TestCheckResourceAttr(resourceName1, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName1, "tags.key", "value1"),
					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_vbd_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_vbd_share_volume_update"),
					resource.TestCheckResourceAttr(resourceName3, "name", rName+"_scsi_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName4, "name", rName+"_scsi_share_volume_update"),
				),
			},
		},
	})
}

func TestAccEvsVolume_withEpsId(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_evs_volume.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
				},
			},
		},
	})
}

func testAccEvsVolume_base() string {
	return `
	variable "volume_configuration" {
	  type = list(object({
		suffix      = string
		device_type = string
		volume_type = string
		multiattach = bool
		iops        = number
		throughput  = number
	  }))
	  default = [
		{
		  suffix = "vbd_normal_volume",
		  device_type = "VBD",
		  volume_type = "SSD",
		  multiattach = false,
		  iops = 0,
		  throughput = 0
		},
		{
		  suffix = "vbd_share_volume",
		  device_type = "VBD",
		  volume_type = "SSD",
		  multiattach = true,
		  iops = 0,
		  throughput = 0
		},
		{
		  suffix = "scsi_normal_volume",
		  device_type = "SCSI",
		  volume_type = "SSD",
		  multiattach = false,
		  iops = 0,
		  throughput = 0
		},
		{
		  suffix = "scsi_share_volume",
		  device_type = "SCSI",
		  volume_type = "SSD",
		  multiattach = true,
		  iops = 0,
		  throughput = 0
		},
	  ]
	}
	
	data "hcso_availability_zones" "test" {}
	`
}

func testAccEvsVolume_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.hcso_availability_zones.test.names[0]
  name              = "%s_${var.volume_configuration[count.index].suffix}"
  size              = 100
  description       = "Created by acc test script."
  volume_type       = var.volume_configuration[count.index].volume_type
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach
  iops              = var.volume_configuration[count.index].iops
  throughput        = var.volume_configuration[count.index].throughput

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsVolume_base(), rName)
}

func testAccEvsVolume_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.hcso_availability_zones.test.names[0]
  name              = "%s_${var.volume_configuration[count.index].suffix}_update"
  size              = 200
  description       = "Updated by acc test script."
  volume_type       = var.volume_configuration[count.index].volume_type
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach
  iops              = var.volume_configuration[count.index].iops
  throughput        = var.volume_configuration[count.index].throughput

  tags = {
    foo1 = "bar"
    key  = "value1"
  }
}
`, testAccEvsVolume_base(), rName)
}

func testAccEvsVolume_epsId(rName string) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

resource "hcso_evs_volume" "test" {
  name                  = "%s"
  description           = "test volume for epsID"
  availability_zone     = data.hcso_availability_zones.test.names[0]
  volume_type           = "SSD"
  size                  = 100
  enterprise_project_id = "%s"
}
`, rName, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST)
}
