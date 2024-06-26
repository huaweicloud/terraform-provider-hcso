package rds

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

func getMysqlBinlogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME

	var (
		mysqlBinlogHttpUrl = "v3/{project_id}/instances/{instance_id}/binlog/clear-policy"
		mysqlBinlogProduct = "rds"
	)
	mysqlBinlogClient, err := cfg.NewServiceClient(mysqlBinlogProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}
	instanceID := state.Primary.Attributes["instance_id"]
	mysqlBinlogPath := mysqlBinlogClient.Endpoint + mysqlBinlogHttpUrl
	mysqlBinlogPath = strings.ReplaceAll(mysqlBinlogPath, "{project_id}", mysqlBinlogClient.ProjectID)
	mysqlBinlogPath = strings.ReplaceAll(mysqlBinlogPath, "{instance_id}", instanceID)
	mysqlBinlogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}
	mysqlBinlogResp, err := mysqlBinlogClient.Request("GET", mysqlBinlogPath, &mysqlBinlogOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS Mysql binlog")
	}
	mysqlBinlogRespBody, err := utils.FlattenResponse(mysqlBinlogResp)
	if err != nil {
		return nil, err
	}
	retentionHours := utils.PathSearch("binlog_retention_hours", mysqlBinlogRespBody, nil)
	if retentionHours == nil || int(retentionHours.(float64)) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return mysqlBinlogRespBody, nil
}

func TestAccMysqlBinlog_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_rds_mysql_binlog.test"
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getMysqlBinlogResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMysqlBinlog_basic(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcso_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "binlog_retention_hours", "6"),
				),
			},
			{
				Config: testMysqlBinlog_basic_update(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcso_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "binlog_retention_hours", "8"),
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

func testAccRdsInstance_mysql(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "hcso_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = true

  db {
    password = "%[3]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }
    
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name, pwd)
}

func testMysqlBinlog_basic(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

resource "hcso_rds_mysql_binlog" "test" {
  instance_id            = hcso_rds_instance.test.id
  binlog_retention_hours = 6
}
`, testAccRdsInstance_mysql(name, dbPwd))
}

func testMysqlBinlog_basic_update(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

resource "hcso_rds_mysql_binlog" "test" {
  instance_id            = hcso_rds_instance.test.id
  binlog_retention_hours = 8
}

`, testAccRdsInstance_mysql(name, dbPwd))
}
