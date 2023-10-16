package dds

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dds/v3/roles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDatabaseRoleFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS v3 client: %s ", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	name := state.Primary.Attributes["name"]
	opts := roles.ListOpts{
		Name:   state.Primary.Attributes["name"],
		DbName: state.Primary.Attributes["db_name"],
	}
	resp, err := roles.List(client, instanceId, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting role (%s) from DDS instance (%s): %v", name, instanceId, err)
	}
	if len(resp) < 1 {
		return nil, fmt.Errorf("unable to find role (%s) from DDS instance (%s)", name, instanceId)
	}
	role := resp[0]
	return &role, nil
}

func TestAccDatabaseRole_basic(t *testing.T) {
	var role roles.Role
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_dds_database_role.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getDatabaseRoleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRole_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.name",
						"hcso_dds_database_role.base", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "inherited_privileges",
						"hcso_dds_database_role.base", "inherited_privileges"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatabaseRoleImportStateIdFunc(),
			},
		},
	})
}

func testAccDatabaseRoleImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, dbName, name string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "hcso_dds_database_role" {
				instanceId = rs.Primary.Attributes["instance_id"]
				dbName = rs.Primary.Attributes["db_name"]
				name = rs.Primary.Attributes["name"]
			}
		}
		if instanceId == "" || dbName == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", instanceId, dbName, name)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, dbName, name), nil
	}
}

func testAccDatabaseRole_base(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_dds_instance" "test" {
  availability_zone = data.hcso_availability_zones.test.names[0]
  vpc_id            = hcso_vpc.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id

  name     = "%s"
  mode     = "Sharding"
  password = "Test@12345678"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDatabaseRole_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_dds_database_role" "base" {
  instance_id = hcso_dds_instance.test.id

  name    = "%[2]s-base"
  db_name = "admin"
}

resource "hcso_dds_database_role" "test" {
  instance_id = hcso_dds_instance.test.id

  name    = "%[2]s"
  db_name = "admin"

  roles {
    name    = hcso_dds_database_role.base.name
    db_name = "admin"
  }
}
`, testAccDatabaseRole_base(rName), rName)
}
