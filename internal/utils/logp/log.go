package logp

import (
	"log"

	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
)

func Printf(format string, v ...interface{}) {
	newFormat := utils.BuildNewFormatByConfig(format)
	log.Printf(newFormat, v...)
}
