package sdrs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccSDRSDomainDataSource_basic(t *testing.T) {
	rName := "data.hcso_sdrs_domain.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSDRSDomainDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "description"),
				),
			},
			{
				Config:      testAccCheckSDRSDomainDataSource_checkError,
				ExpectError: regexp.MustCompile(`your query returned no results. Please change your search criteria and try again`),
			},
		},
	})
}

const testAccCheckSDRSDomainDataSource_basic = `
data "hcso_sdrs_domain" "test" {}
`

const testAccCheckSDRSDomainDataSource_checkError = `
data "hcso_sdrs_domain" "test" {
  name = "error_check"
}
`
