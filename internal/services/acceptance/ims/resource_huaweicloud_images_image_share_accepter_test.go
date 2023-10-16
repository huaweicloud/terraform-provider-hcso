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

func getImsImageShareAccepterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

	imageId := state.Primary.Attributes["image_id"]
	getImageQueryParams := buildGetImageQueryParams(imageId)
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

func TestAccImsImageShareAccepter_basic(t *testing.T) {
	var obj interface{}

	rName := "hcso_images_image_share_accepter.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageShareAccepterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckSourceImage(t)
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageShareAccepter_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "image_id", acceptance.HCSO_IMAGE_SHARE_SOURCE_IMAGE_ID),
				),
			},
		},
	})
}

func testImsImageShareAccepter_basic() string {
	return fmt.Sprintf(`
resource "hcso_images_image_share_accepter" "test" {
 image_id = "%s"
}
`, acceptance.HCSO_IMAGE_SHARE_SOURCE_IMAGE_ID)
}
