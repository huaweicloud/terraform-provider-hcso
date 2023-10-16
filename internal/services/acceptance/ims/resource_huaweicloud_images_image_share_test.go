package ims

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getImsImageShareResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getImage: Query IMS image
	var (
		getImageHttpUrl = "v2/cloudimages"
		getImageProduct = "ims"
	)
	getImageClient, err := cfg.NewServiceClient(getImageProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS Client: %s", err)
	}

	getImagePath := getImageClient.Endpoint + getImageHttpUrl
	getImagePath = strings.ReplaceAll(getImagePath, "{project_id}", getImageClient.ProjectID)

	getImageQueryParams := buildGetImageQueryParams(state.Primary.ID)
	getImagePath += getImageQueryParams

	getImageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getImageResp, err := getImageClient.Request("GET", getImagePath, &getImageOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS image: %s", err)
	}

	getImageRespBody, err := utils.FlattenResponse(getImageResp)
	if err != nil {
		return nil, err
	}

	images := utils.PathSearch("images", getImageRespBody, nil)
	if images == nil || len(images.([]interface{})) == 0 {
		return nil, fmt.Errorf("error get IMS share image")
	}

	return make(map[string]interface{}), nil
}

func buildGetImageQueryParams(id string) string {
	res := ""
	res = fmt.Sprintf("%s?id=%v", res, id)
	res = fmt.Sprintf("%s&__imagetype=%v", res, "shared")
	return res
}

func TestAccImsImageShare_basic(t *testing.T) {
	var obj interface{}

	imageName := acceptance.RandomAccResourceName()
	rName := "hcso_images_image_share.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageShareResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckReplication(t)
			acceptance.TestAccPreCheckProjectId(t)
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageShare_basic(imageName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "source_image_id",
						"hcso_images_image.test", "id"),
				),
			},
			{
				Config: testImsImageShare_update(imageName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "source_image_id",
						"hcso_images_image.test", "id"),
				),
			},
		},
	})
}

func testImsImageShare_basic(imageName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_images_image_share" "test" {
 source_image_id    = hcso_images_image.test.id
 target_project_ids = ["%[2]s"]
}
`, testAccImsImage_basic(imageName), acceptance.HCSO_DEST_PROJECT_ID)
}

func testImsImageShare_update(imageName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_images_image_share" "test" {
 source_image_id    = hcso_images_image.test.id
 target_project_ids = ["%[2]s"]
}
`, testAccImsImage_basic(imageName), acceptance.HCSO_DEST_PROJECT_ID_TEST)
}
