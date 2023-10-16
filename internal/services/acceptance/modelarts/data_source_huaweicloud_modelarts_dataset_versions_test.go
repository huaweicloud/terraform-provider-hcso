package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDataSourceDatasetVersions_basic(t *testing.T) {
	dataSourceName := "data.hcso_modelarts_dataset_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasetVersions_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.id",
						"hcso_modelarts_dataset_version.test", "version_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.description",
						"hcso_modelarts_dataset_version.test", "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.name",
						"hcso_modelarts_dataset_version.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.split_ratio",
						"hcso_modelarts_dataset_version.test", "split_ratio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.status",
						"hcso_modelarts_dataset_version.test", "status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.name",
						"hcso_modelarts_dataset_version.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.files",
						"hcso_modelarts_dataset_version.test", "files"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.storage_path",
						"hcso_modelarts_dataset_version.test", "storage_path"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.is_current",
						"hcso_modelarts_dataset_version.test", "is_current"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.created_at",
						"hcso_modelarts_dataset_version.test", "created_at"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.updated_at",
						"hcso_modelarts_dataset_version.test", "updated_at"),
				),
			},
			{
				Config: testAccDataSourceDatasetVersions_name(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "versions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceDatasetVersions_basic(rName, obsName string) string {
	datasetVersion := testAccDatasetVersion_basic(rName, obsName)
	return fmt.Sprintf(`
%s

data "hcso_modelarts_dataset_versions" "test" {
  dataset_id  = hcso_modelarts_dataset.test.id
  split_ratio = "0,2.9"

  depends_on = [
    hcso_modelarts_dataset_version.test
  ]
}
`, datasetVersion)
}

func testAccDataSourceDatasetVersions_name(rName, obsName string) string {
	datasetVersion := testAccDatasetVersion_basic(rName, obsName)
	return fmt.Sprintf(`
%s

data "hcso_modelarts_dataset_versions" "test" {
  dataset_id  = hcso_modelarts_dataset.test.id
  split_ratio = "0,2.9"
  name        = "wrong_name"

  depends_on = [
    hcso_modelarts_dataset_version.test
  ]
}
`, datasetVersion)
}
