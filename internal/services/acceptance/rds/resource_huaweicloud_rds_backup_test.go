package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getBackupResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getBackup: Query the RDS manual backup
	var (
		getBackupHttpUrl = "v3/{project_id}/backups"
		getBackupProduct = "rds"
	)
	getBackupClient, err := config.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Backup Client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.Replace(getBackupPath, "{project_id}", getBackupClient.ProjectID, -1)

	getBackupqueryParams := fmt.Sprintf("?instance_id=%s&backup_id=%s",
		state.Primary.Attributes["instance_id"], state.Primary.ID)
	getBackupPath = getBackupPath + getBackupqueryParams
	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Backup: %s", err)
	}

	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Backup: %s", err)
	}

	count := utils.PathSearch("total_count", getBackupRespBody, 0)
	if fmt.Sprintf("%v", count) == "0" {
		return nil, fmt.Errorf("error retrieving Backup: %s", err)
	}

	return getBackupRespBody, nil
}

func TestAccBackup_mysql_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_rds_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackup_mysql_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcso_rds_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "size"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBackupImportStateFunc(rName),
			},
		},
	})
}

func TestAccBackup_sqlserver_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_rds_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackup_sqlserver_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcso_rds_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "size"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBackupImportStateFunc(rName),
			},
		},
	})
}

func TestAccBackup_pg_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_rds_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackup_pg_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcso_rds_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "size"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBackupImportStateFunc(rName),
			},
		},
	})
}

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testBackup_mysql_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

data "hcso_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "hcso_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = hcso_networking_secgroup.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  vpc_id            = hcso_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "hcso_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = hcso_rds_instance.test.id
}
`, common.TestBaseNetwork(name), name)
}

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testBackup_sqlserver_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "hcso_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8631
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "hcso_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = hcso_rds_instance.test.id
}
`, common.TestBaseNetwork(name), name)
}

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testBackup_pg_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

data "hcso_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 8
}

resource "hcso_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = hcso_networking_secgroup.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  vpc_id            = hcso_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8632
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "hcso_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = hcso_rds_instance.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccBackupImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
