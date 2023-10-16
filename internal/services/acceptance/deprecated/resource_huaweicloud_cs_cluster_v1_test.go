// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package deprecated

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCsClusterV1_basic(t *testing.T) {
	rName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckCsClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCsClusterV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCsClusterV1Exists(),
				),
			},
		},
	})
}

func testAccCsClusterV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_cs_cluster_v1" "cluster" {
  name = "%s"
}
	`, rName)
}

func testAccCheckCsClusterV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.CloudStreamV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_cs_cluster_v1" {
			continue
		}

		url, err := replaceVarsForTest(rs, "reserved_cluster/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmtp.Errorf("hcso_cs_cluster_v1 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCsClusterV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.CloudStreamV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["hcso_cs_cluster_v1.cluster"]
		if !ok {
			return fmtp.Errorf("Error checking hcso_cs_cluster_v1.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "reserved_cluster/{id}")
		if err != nil {
			return fmtp.Errorf("Error checking hcso_cs_cluster_v1.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("hcso_cs_cluster_v1.cluster is not exist")
			}
			return fmtp.Errorf("Error checking hcso_cs_cluster_v1.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}
