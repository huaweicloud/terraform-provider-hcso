package config

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	providerUserAgent string = "terraform-provider-iac"
)

type HCSOConfig struct {
	config.Config
}
