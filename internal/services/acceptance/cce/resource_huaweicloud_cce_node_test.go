package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func TestAccNode_basic(t *testing.T) {
	var node nodes.Nodes

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "hcso_cce_node.test"
	// clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "hcso_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNode_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.extend_params.test_key", "test_val"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.extend_params.test_key", "test_val"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodeImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"taints", "extend_params",
				},
			},
			{
				Config: testAccNode_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.extend_params.test_key", "test_val"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.extend_params.test_key", "test_val"),
				),
			},
		},
	})
}

func TestAccNode_eip(t *testing.T) {
	var node nodes.Nodes

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_cce_node.test"
	// clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "hcso_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNode_auto_assign_eip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)),
				),
			},
			{
				Config: testAccNode_existing_eip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)),
				),
			},
		},
	})
}

func TestAccNode_volume_encryption(t *testing.T) {
	var node nodes.Nodes

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_cce_node.test"
	// clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "hcso_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNode_volume_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "data_volumes.0.kms_key_id"),
				),
			},
		},
	})
}

func TestAccNode_password(t *testing.T) {
	var node nodes.Nodes

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_cce_node.test"
	// clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "hcso_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNode_password(rName),
				Check: resource.ComposeTestCheckFunc(
					// use provisioner to check whether the node can login successfully
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccNode_password_update(rName),
				Check: resource.ComposeTestCheckFunc(
					// use provisioner to check whether the node can login successfully
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccNode_storage(t *testing.T) {
	var node nodes.Nodes

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_cce_node.test"
	// clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "hcso_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNode_storage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckNodeDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := cfg.CceV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "hcso_cce_cluster" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "hcso_cce_node" {
			continue
		}

		_, err := nodes.Get(cceClient, clusterId, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("node still exists")
		}
	}

	return nil
}

func testAccCheckNodeExists(n string, cluster string, node *nodes.Nodes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmt.Errorf("cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		if c.Primary.ID == "" {
			return fmt.Errorf("cluster id is not set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := config.CceV3Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CCE client: %s", err)
		}

		found, err := nodes.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("node not found")
		}

		*node = *found

		return nil
	}
}

func testAccCCENodeImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["hcso_cce_cluster.test"]
		if !ok {
			return "", fmt.Errorf("cluster not found: %s", cluster)
		}
		node, ok := s.RootModule().Resources["hcso_cce_node.test"]
		if !ok {
			return "", fmt.Errorf("node not found: %s", node)
		}
		if cluster.Primary.ID == "" || node.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", cluster.Primary.ID, node.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", cluster.Primary.ID, node.Primary.ID), nil
	}
}

func testAccNode_Base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

data "hcso_compute_flavors" "test" {
  availability_zone = data.hcso_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcso_kps_keypair" "test" {
  name = "%[2]s"
}

resource "hcso_cce_cluster" "test" {
  name                   = "%[2]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcso_vpc.test.id
  subnet_id              = hcso_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, common.TestVpc(rName), rName)
}

func testAccNode_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  key_pair          = hcso_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
    extend_params = {
      test_key = "test_val"
    }
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
    extend_params = {
      test_key = "test_val"
    }
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }
}
`, testAccNode_Base(rName), rName)
}

func testAccNode_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  key_pair          = hcso_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
    extend_params = {
      test_key = "test_val"
    }
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
    extend_params = {
      test_key = "test_val"
    }
  }

  tags = {
    foo = "bar"
    key = "value_update"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }

  lifecycle {
    ignore_changes = [
      taints, extend_params
    ]
  }
}
`, testAccNode_Base(rName), updateName)
}

func testAccNode_auto_assign_eip(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  key_pair          = hcso_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  // Assign EIP
  iptype="5_bgp"
  bandwidth_charge_mode="traffic"
  sharetype= "PER"
  bandwidth_size= 100
}
`, testAccNode_Base(rName), rName)
}

func testAccNode_existing_eip(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  key_pair          = hcso_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  // Assign existing EIP
  eip_id = hcso_vpc_eip.test.id
}
`, testAccNode_Base(rName), rName)
}

func testAccNode_volume_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  key_pair          = hcso_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = hcso_kms_key.test.id
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccNode_Base(rName), rName, rName)
}

func testAccNode_password(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  password          = "Test@123"
  eip_id            = hcso_vpc_eip.test.id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  provisioner "remote-exec" {
    connection {
      user     = "root"
      password = "Test@123"
      host     = hcso_vpc_eip.test.address
      port     = "22"
      timeout  = "60s"
    }

    inline = [
      "date"
    ]
  }
}
`, testAccNode_Base(rName), rName)
}

func testAccNode_password_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  password          = "Test@12345"
  eip_id            = hcso_vpc_eip.test.id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  provisioner "remote-exec" {
    connection {
      user     = "root"
      password = "Test@12345"
      host     = hcso_vpc_eip.test.address
      port     = "22"
      timeout  = "60s"
    }

    inline = [
      "date"
    ]
  }
}
`, testAccNode_Base(rName), rName)
}

func testAccNode_storage(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "hcso_cce_node" "test" {
  cluster_id        = hcso_cce_cluster.test.id
  name              = "%s"
  flavor_id         = data.hcso_compute_flavors.test.ids[0]
  availability_zone = data.hcso_availability_zones.test.names[0]
  key_pair          = hcso_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = hcso_kms_key.test.id
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = hcso_kms_key.test.id
  }

  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "100"
      match_label_count = 1
    }

    selectors {
      name                           = "user"
      type                           = "evs"
      match_label_size               = "100"
      match_label_metadata_encrypted = "1"
      match_label_metadata_cmkid     = hcso_kms_key.test.id
      match_label_count              = "1"
    }

    groups {
      name           = "vgpaas"
      selector_names = ["cceUse"]
      cce_managed    = true

      virtual_spaces {
        name        = "kubernetes"
        size        = "10%%"
        lvm_lv_type = "linear"
      }

      virtual_spaces {
        name            = "runtime"
        size            = "90%%"
        runtime_lv_type = "linear"
      }
    }

    groups {
      name           = "vguser"
      selector_names = ["user"]

      virtual_spaces {
        name        = "user"
        size        = "100%%"
        lvm_lv_type = "linear"
        lvm_path    = "/workspace"
      }
    }
  }
}
`, testAccNode_Base(rName), rName, rName)
}
