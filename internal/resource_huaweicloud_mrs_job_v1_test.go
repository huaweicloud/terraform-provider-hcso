package internal

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v1/job"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccMRSV1Job_basic(t *testing.T) {
	var jobGet job.Job

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMRSV1JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccMRSV1JobConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1JobExists("hcso_mrs_job_v1.job1", &jobGet),
					resource.TestCheckResourceAttr(
						"hcso_mrs_job_v1.job1", "job_state", "Completed"),
				),
			},
		},
	})
}

func testAccCheckMRSV1JobDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	mrsClient, err := config.MrsV1Client(HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating internal mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_mrs_job_v1" {
			continue
		}

		_, err := job.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return nil
			}
			// Temporarily pass 500 error due to some issue on server side.
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				return nil
			}
			return fmtp.Errorf("job still exists. err : %s", err)
		}
	}

	return nil
}

func testAccCheckMRSV1JobExists(n string, jobGet *job.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*config.Config)
		mrsClient, err := config.MrsV1Client(HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating internal mrs client: %s ", err)
		}

		found, err := job.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Job not found. ")
		}

		*jobGet = *found

		return nil
	}
}

var TestAccMRSV1JobConfig_basic = fmt.Sprintf(`
resource "hcso_mrs_cluster_v1" "cluster1" {
  cluster_name = "mrs-cluster-acc"
  region = "%s"
  billing_type = 12
  master_node_num = 2
  core_node_num = 3
  master_node_size = "c3.4xlarge.2.linux.bigdata"
  core_node_size = "c3.xlarge.4.linux.bigdata"
  available_zone_id = "ae04cf9d61544df3806a3feeb401b204"
  vpc_id = "%s"
  subnet_id = "%s"
  cluster_version = "MRS 1.6.3"
  volume_type = "SATA"
  volume_size = 100
  safe_mode = 0
  cluster_type = 0
  node_public_cert_name = "KeyPair-ci"
  cluster_admin_secret = ""
  component_list {
      component_name = "Hadoop"
  }
  component_list {
      component_name = "Spark"
  }
  component_list {
      component_name = "Hive"
  }
}

resource "hcso_mrs_job_v1" "job1" {
  job_type = 1
  job_name = "test_mapreduce_job1"
  cluster_id = "${hcso_mrs_cluster_v1.cluster1.id}"
  jar_path = "s3a://tf-mrs/program/hadoop-mapreduce-examples-2.7.5.jar"
  input = "s3a://tf-mrs/input/"
  output = "s3a://tf-mrs/output/"
  job_log = "s3a://tf-mrs/joblog/"
  arguments = "wordcount"
}`, HCSO_REGION_NAME, HCSO_VPC_ID, HCSO_NETWORK_ID)
