package ucs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getClusterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getCluster: Query the UCS Cluster detail
	var (
		getClusterHttpUrl = "v1/clusters/{id}"
		getClusterProduct = "ucs"
	)
	getClusterClient, err := cfg.NewServiceClient(getClusterProduct, "")
	if err != nil {
		return nil, fmt.Errorf("error creating UCS Client: %s", err)
	}

	getClusterPath := getClusterClient.Endpoint + getClusterHttpUrl
	getClusterPath = strings.ReplaceAll(getClusterPath, "{id}", state.Primary.ID)

	getClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getClusterResp, err := getClusterClient.Request("GET", getClusterPath, &getClusterOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Cluster: %s", err)
	}

	getClusterRespBody, err := utils.FlattenResponse(getClusterResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Cluster: %s", err)
	}

	return getClusterRespBody, nil
}

func TestAccCluster_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcso_ucs_cluster.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCluster_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "attachedcluster"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "privatek8s"),
					resource.TestCheckResourceAttr(rName, "manage_type", "discrete"),
					resource.TestCheckResourceAttr(rName, "cluster_name", name),
					resource.TestCheckResourceAttr(rName, "country", "CN"),
					resource.TestCheckResourceAttr(rName, "city", "110000"),
				),
			},
			{
				Config: testCluster_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "attachedcluster"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "privatek8s"),
					resource.TestCheckResourceAttr(rName, "manage_type", "grouped"),
					resource.TestCheckResourceAttr(rName, "cluster_name", name),
					resource.TestCheckResourceAttr(rName, "country", "US"),
					resource.TestCheckResourceAttr(rName, "city", "US"),
					resource.TestCheckResourceAttrPair(rName, "fleet_id", "hcso_ucs_fleet.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"annotations",
				},
			},
		},
	})
}

func TestAccCluster_cce(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcso_ucs_cluster.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCluster_cce(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "self"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "cce"),
					resource.TestCheckResourceAttr(rName, "manage_type", "grouped"),
					resource.TestCheckResourceAttrPair(rName, "fleet_id", "hcso_ucs_fleet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "hcso_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "cluster_region", acceptance.HCSO_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_project_id", acceptance.HCSO_PROJECT_ID),
				),
			},
			{
				Config: testCluster_cce_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "self"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "cce"),
					resource.TestCheckResourceAttr(rName, "manage_type", "discrete"),
					resource.TestCheckResourceAttr(rName, "fleet_id", ""),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "hcso_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "cluster_region", acceptance.HCSO_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_project_id", acceptance.HCSO_PROJECT_ID),
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

func testCluster_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcso_vpc.test.id
  subnet_id              = hcso_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "hcso_ucs_fleet" "test" {
  name        = "%[2]s"
  description = "created by terraform"
}
`, common.TestVpc(rName), rName)
}

func testCluster_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_ucs_cluster" "test" {
  category     = "attachedcluster"
  cluster_type = "privatek8s"
  cluster_name = "%s"
  country      = "CN"
  city         = "110000"

  annotations = {
    "kubeconfig" = hcso_cce_cluster.test.kube_config_raw
  }
}
`, testCluster_base(name), name)
}

func testCluster_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_ucs_cluster" "test" {
  category     = "attachedcluster"
  cluster_type = "privatek8s"
  fleet_id     = hcso_ucs_fleet.test.id
  cluster_name = "%s"
  country      = "US"
  city         = "US"

  annotations = {
    "kubeconfig" = hcso_cce_cluster.test.kube_config_raw
  }
}
`, testCluster_base(name), name)
}

func testCluster_cce(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_ucs_cluster" "test" {
  category           = "self"
  cluster_type       = "cce"
  fleet_id           = hcso_ucs_fleet.test.id
  cluster_id         = hcso_cce_cluster.test.id
  cluster_region     = "%s"
  cluster_project_id = "%s"
}
`, testCluster_base(name), acceptance.HCSO_REGION_NAME, acceptance.HCSO_PROJECT_ID)
}

func testCluster_cce_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_ucs_cluster" "test" {
  category           = "self"
  cluster_type       = "cce"
  cluster_id         = hcso_cce_cluster.test.id
  cluster_region     = "%s"
  cluster_project_id = "%s"
}
`, testCluster_base(name), acceptance.HCSO_REGION_NAME, acceptance.HCSO_PROJECT_ID)
}
