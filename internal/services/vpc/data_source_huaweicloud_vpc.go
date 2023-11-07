package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceVpcV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"routes": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "use hcso_vpc_route_table data source to get all routes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	vpcClient, err := conf.NetworkingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		CIDR:                d.Get("cidr").(string),
		EnterpriseProjectID: conf.DataGetEnterpriseProjectID(d),
	}

	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVpcs) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	vpc := refinedVpcs[0]

	log.Printf("[INFO] Retrieved vpc using given filter %s: %+v", vpc.ID, vpc)
	d.SetId(vpc.ID)

	d.Set("region", conf.GetRegion(d))
	d.Set("name", vpc.Name)
	d.Set("cidr", vpc.CIDR)
	d.Set("enterprise_project_id", vpc.EnterpriseProjectID)
	d.Set("status", vpc.Status)
	d.Set("description", vpc.Description)

	var s []map[string]interface{}
	for _, route := range vpc.Routes {
		mapping := map[string]interface{}{
			"destination": route.DestinationCIDR,
			"nexthop":     route.NextHop,
		}
		s = append(s, mapping)
	}
	d.Set("routes", s)

	// save VirtualPrivateCloudV2 tags
	if vpcV2Client, err := conf.NetworkingV2Client(conf.GetRegion(d)); err == nil {
		if resourceTags, err := tags.Get(vpcV2Client, "vpcs", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				return diag.Errorf("error saving tags to state for VPC (%s): %s", d.Id(), err)
			}
		} else {
			log.Printf("[WARN] Error fetching tags of VPC (%s): %s", d.Id(), err)
		}
	}

	return nil
}
