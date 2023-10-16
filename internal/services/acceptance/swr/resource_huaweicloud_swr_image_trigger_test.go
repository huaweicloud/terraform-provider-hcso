package swr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getSwrImageTriggerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getSwrImageTrigger: Query SWR image trigger
	var (
		getSwrImageTriggerHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}"
		getSwrImageTriggerProduct = "swr"
	)
	getSwrImageTriggerClient, err := cfg.NewServiceClient(getSwrImageTriggerProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <organization_name>/<repository_name>/<trigger_name>")
	}
	organization := parts[0]
	repository := parts[1]
	trigger := parts[2]

	getSwrImageTriggerPath := getSwrImageTriggerClient.Endpoint + getSwrImageTriggerHttpUrl
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{namespace}", organization)
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{repository}", repository)
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{trigger}", trigger)

	getSwrImageTriggerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageTriggerResp, err := getSwrImageTriggerClient.Request("GET",
		getSwrImageTriggerPath, &getSwrImageTriggerOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR image trigger: %s", err)
	}
	return utils.FlattenResponse(getSwrImageTriggerResp)
}

func TestAccSwrImageTrigger_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_swr_image_trigger.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImageTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkloadType(t)
			acceptance.TestAccPreCheckWorkloadName(t)
			acceptance.TestAccPreCheckCceClusterId(t)
			acceptance.TestAccPreCheckWorkloadNameSpace(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSwrImageTrigger_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"hcso_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"hcso_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "workload_type", acceptance.HCSO_WORKLOAD_TYPE),
					resource.TestCheckResourceAttr(rName, "workload_name", acceptance.HCSO_WORKLOAD_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HCSO_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "namespace", acceptance.HCSO_WORKLOAD_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "condition_value", ".*"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "cce"),
					resource.TestCheckResourceAttr(rName, "condition_type", "all"),
				),
			},
			{
				Config: testSwrImageTrigger_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"hcso_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"hcso_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "workload_type", acceptance.HCSO_WORKLOAD_TYPE),
					resource.TestCheckResourceAttr(rName, "workload_name", acceptance.HCSO_WORKLOAD_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HCSO_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "namespace", acceptance.HCSO_WORKLOAD_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "condition_value", ".*"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "cce"),
					resource.TestCheckResourceAttr(rName, "condition_type", "all"),
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

func testSwrImageTrigger_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_swr_image_trigger" "test" {
  organization    = hcso_swr_organization.test.name
  repository      = hcso_swr_repository.test.name
  workload_type   = "%[2]s"
  workload_name   = "%[3]s"
  cluster_id      = "%[4]s"
  namespace       = "%[5]s"
  condition_value = ".*"
  name            = "%[6]s"
  type            = "cce"
  condition_type  = "all"
}
`, testAccSWRRepository_basic(name), acceptance.HCSO_WORKLOAD_TYPE, acceptance.HCSO_WORKLOAD_NAME,
		acceptance.HCSO_CCE_CLUSTER_ID, acceptance.HCSO_WORKLOAD_NAMESPACE, name)
}

func testSwrImageTrigger_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_swr_image_trigger" "test" {
  organization    = hcso_swr_organization.test.name
  repository      = hcso_swr_repository.test.name
  workload_type   = "%[2]s"
  workload_name   = "%[3]s"
  cluster_id      = "%[4]s"
  namespace       = "%[5]s"
  condition_value = ".*"
  enabled         = "false"
  name            = "%[6]s"
  type            = "cce"
  condition_type  = "all"
}
`, testAccSWRRepository_basic(name), acceptance.HCSO_WORKLOAD_TYPE, acceptance.HCSO_WORKLOAD_NAME,
		acceptance.HCSO_CCE_CLUSTER_ID, acceptance.HCSO_WORKLOAD_NAMESPACE, name)
}
