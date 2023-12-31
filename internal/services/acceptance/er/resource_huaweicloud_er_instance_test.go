package er

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
)

func getInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getInstance: Query the Enterprise router instance detail
	var (
		getInstanceHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		getInstanceProduct = "er"
	)
	getInstanceClient, err := config.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.Replace(getInstancePath, "{project_id}", getInstanceClient.ProjectID, -1)
	getInstancePath = strings.Replace(getInstancePath, "{id}", state.Primary.ID, -1)

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Instance: %s", err)
	}
	return utils.FlattenResponse(getInstanceResp)
}

func TestAccInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_er_instance.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic_step1(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.hcso_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testInstance_basic_step2(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testInstance_basic_step1(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

resource "hcso_er_instance" "test" {
  availability_zones = slice(data.hcso_availability_zones.test.names, 0, 1)

  name = "%[2]s"
  asn  = %[3]d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, acceptance.HCSO_AVAILABILITY_ZONE, name, bgpAsNum)
}

func testInstance_basic_step2(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

resource "hcso_er_instance" "test" {
  availability_zones = slice(data.hcso_availability_zones.test.names, 0, 1)

  name = "%[2]s"
  asn  = %[3]d

  tags = {
    foo    = "baar"
    newkey = "value"
  }
}
`, acceptance.HCSO_AVAILABILITY_ZONE, name, bgpAsNum)
}
