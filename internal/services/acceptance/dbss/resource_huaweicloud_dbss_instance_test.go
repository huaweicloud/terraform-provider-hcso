package dbss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	jmespath "github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dbss"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
)

func getInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getInstance: Query the DBSS instance detail
	var (
		getInstanceHttpUrl = "v1/{project_id}/dbss/audit/instances"
		getInstanceProduct = "dbss"
	)
	getInstanceClient, err := config.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)

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

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}

	instances, err := jmespath.Search("servers", getInstanceRespBody)
	if err != nil {
		return nil, fmt.Errorf("error parsing servers from response= %#v", getInstanceRespBody)
	}

	return dbss.FilterInstances(instances.([]interface{}), state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_dbss_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "flavor", "c3ne.xlarge.4"),
					resource.TestCheckResourceAttr(rName, "product_id", "00301-225396-0--0"),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", "dbss.bypassaudit.low"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode", "enterprise_project_id", "flavor", "period", "period_unit",
					"product_id",
				},
			},
		},
	})
}

func testInstance_base() string {
	return fmt.Sprintf(`
data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name   = "subnet-default"
  vpc_id = data.hcso_vpc.test.id
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_availability_zones" "test" {
  region = "%s"
} 
`, acceptance.HCSO_REGION_NAME)
}

func testInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_dbss_instance" "test" {
  name               = "%s"
  description        = "terraform test"
  flavor             = "c3ne.xlarge.4"
  product_id         = "00301-225396-0--0"
  resource_spec_code = "dbss.bypassaudit.low"
  availability_zone  = data.hcso_availability_zones.test.names[0]
  vpc_id             = data.hcso_vpc.test.id
  subnet_id          = data.hcso_vpc_subnet.test.id
  security_group_id  = data.hcso_networking_secgroup.test.id
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1
}
`, testInstance_base(), name)
}
