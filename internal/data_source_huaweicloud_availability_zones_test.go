package internal

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAvailabilityZones_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAvailabilityZonesConfig_all,
				Check: resource.ComposeTestCheckFunc(resource.TestMatchResourceAttr("data.hcso_availability_zones.all",
					"names.#", regexp.MustCompile(`[1-9]\d*`))),
			},
		},
	})
}

const testAccAvailabilityZonesConfig_all = `data "hcso_availability_zones" "all" {}`
