package aom

import (
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getPrometheusInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, _ := httpclient_go.NewHttpClientGo(conf, "aom", acceptance.HCSO_REGION_NAME)
	client.WithMethod(httpclient_go.MethodGet).WithUrlWithoutEndpoint(conf, "aom",
		conf.Region, "v1/"+conf.GetProjectID(conf.Region)+"/prometheus-instances?action=prom_for_cloud_service")

	resp, err := client.Do()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rlt := &entity.PrometheusInstanceParams{}

	err = json.Unmarshal(body, rlt)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("error getting HuaweiCloud Resource")
}

func TestAccAOMPrometheusInstance_basic(t *testing.T) {
	var ar []entity.AddAlarmRuleParams
	resourceName := "hcso_aom_prometheus_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getPrometheusInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckInternal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMPrometheusInstance_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAOMPrometheusInstance_basic() string {
	return `
resource "hcso_aom_prometheus_instance" "test" {
  prom_for_cloud_service  {
  ces_metric_namespaces =	["SYS.ELB","SYS.VPC","SYS.DMS","SYS.RDS"]
 }
}
`
}
