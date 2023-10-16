package ces

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

func getCesAlarmTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getAlarmTemplate: Query CES alarm template
	var (
		getAlarmTemplateHttpUrl = "v2/{project_id}/alarm-templates/{template_id}"
		getAlarmTemplateProduct = "ces"
	)
	getAlarmTemplateClient, err := cfg.NewServiceClient(getAlarmTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CES Client: %s", err)
	}

	getAlarmTemplatePath := getAlarmTemplateClient.Endpoint + getAlarmTemplateHttpUrl
	getAlarmTemplatePath = strings.ReplaceAll(getAlarmTemplatePath, "{project_id}", getAlarmTemplateClient.ProjectID)
	getAlarmTemplatePath = strings.ReplaceAll(getAlarmTemplatePath, "{template_id}", state.Primary.ID)

	getAlarmTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAlarmTemplateResp, err := getAlarmTemplateClient.Request("GET", getAlarmTemplatePath, &getAlarmTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CES alarm template: %s", err)
	}
	return utils.FlattenResponse(getAlarmTemplateResp)
}

func TestAccCesAlarmTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "hcso_ces_alarm_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCesAlarmTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCesAlarmTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "It is a test template"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttr(rName, "policies.0.namespace", "SYS.APIG"),
					resource.TestCheckResourceAttr(rName, "policies.0.dimension_name", "api_id"),
					resource.TestCheckResourceAttr(rName, "policies.0.metric_name", "req_count_2xx"),
					resource.TestCheckResourceAttr(rName, "policies.0.period", "1"),
					resource.TestCheckResourceAttr(rName, "policies.0.filter", "average"),
					resource.TestCheckResourceAttr(rName, "policies.0.comparison_operator", "="),
					resource.TestCheckResourceAttr(rName, "policies.0.value", "10"),
					resource.TestCheckResourceAttr(rName, "policies.0.unit", "times/minute"),
					resource.TestCheckResourceAttr(rName, "policies.0.count", "3"),
					resource.TestCheckResourceAttr(rName, "policies.0.alarm_level", "2"),
					resource.TestCheckResourceAttr(rName, "policies.0.suppress_duration", "300"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "association_alarm_total"),
				),
			},
			{
				Config: testCesAlarmTemplate_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "It is an update template"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttr(rName, "policies.0.namespace", "SYS.DDS"),
					resource.TestCheckResourceAttr(rName, "policies.0.dimension_name", "mongodb_instance_id"),
					resource.TestCheckResourceAttr(rName, "policies.0.metric_name", "mongo003_insert_ps"),
					resource.TestCheckResourceAttr(rName, "policies.0.period", "300"),
					resource.TestCheckResourceAttr(rName, "policies.0.filter", "max"),
					resource.TestCheckResourceAttr(rName, "policies.0.comparison_operator", "<"),
					resource.TestCheckResourceAttr(rName, "policies.0.value", "300"),
					resource.TestCheckResourceAttr(rName, "policies.0.unit", "times/second"),
					resource.TestCheckResourceAttr(rName, "policies.0.count", "5"),
					resource.TestCheckResourceAttr(rName, "policies.0.alarm_level", "3"),
					resource.TestCheckResourceAttr(rName, "policies.0.suppress_duration", "3600"),
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

func testCesAlarmTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "hcso_ces_alarm_template" "test" {
  name        = "%s" 
  description = "It is a test template"

  policies {
	namespace           = "SYS.APIG"
	dimension_name      = "api_id"
	metric_name         = "req_count_2xx"
	period              = 1
	filter              = "average"
	comparison_operator = "="
	value               = "10"
	unit                = "times/minute"
	count               = 3
	alarm_level         = 2
	suppress_duration   = 300
  }
}
`, name)
}

func testCesAlarmTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcso_ces_alarm_template" "test" {
  name        = "%s"
  description = "It is an update template"

  policies {
	namespace           = "SYS.DDS"
	dimension_name      = "mongodb_instance_id"
	metric_name         = "mongo003_insert_ps"
	period              = 300
	filter              = "max"
	comparison_operator = "<"
	value               = "300"
	unit                = "times/second"
	count               = 5
	alarm_level         = 3
	suppress_duration   = 3600
  }
}
`, name)
}
