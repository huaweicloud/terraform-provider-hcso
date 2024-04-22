package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/antiddos"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/tms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"

	"github.com/huaweicloud/terraform-provider-hcso/internal/hcso_config"
)

const (
	defaultCloud string = "myhuaweicloud.com"
)

// Provider returns a schema.Provider for HuaweiCloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["region"],
				InputDefault: "cn-north-1",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_REGION_NAME",
				}, nil),
			},

			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["access_key"],
				RequiredWith: []string{"secret_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_ACCESS_KEY",
				}, nil),
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_SECRET_KEY",
				}, nil),
			},

			"security_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["security_token"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("HCSO_SECURITY_TOKEN", nil),
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_DOMAIN_ID",
				}, ""),
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_DOMAIN_NAME",
				}, ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_USER_NAME",
				}, ""),
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_USER_ID",
				}, ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: descriptions["password"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_USER_PASSWORD",
				}, ""),
			},

			"assume_role": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_agency_name"],
							DefaultFunc: schema.EnvDefaultFunc("HCSO_ASSUME_ROLE_AGENCY_NAME", nil),
						},
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_domain_name"],
							DefaultFunc: schema.EnvDefaultFunc("HCSO_ASSUME_ROLE_DOMAIN_NAME", nil),
						},
					},
				},
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_PROJECT_ID",
				}, nil),
			},

			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_PROJECT_NAME",
				}, nil),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_id"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_ID", ""),
			},

			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_name"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_NAME", ""),
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["token"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_AUTH_TOKEN",
				}, ""),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["insecure"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_INSECURE",
				}, false),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"agency_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_NAME", nil),
				Description:  descriptions["agency_name"],
				RequiredWith: []string{"agency_domain_name"},
			},

			"agency_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_DOMAIN_NAME", nil),
				Description:  descriptions["agency_domain_name"],
				RequiredWith: []string{"agency_name"},
			},

			"delegated_project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELEGATED_PROJECT", ""),
				Description: descriptions["delegated_project"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["auth_url"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_AUTH_URL",
					"OS_AUTH_URL",
				}, nil),
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cloud"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_CLOUD", ""),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoints"],
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"regional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["regional"],
			},

			"shared_config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_config_file"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_SHARED_CONFIG_FILE", ""),
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_PROFILE", ""),
			},

			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["enterprise_project_id"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_ENTERPRISE_PROJECT_ID", ""),
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_MAX_RETRIES", 5),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"hcso_availability_zones": huaweicloud.DataSourceAvailabilityZones(),

			"hcso_cbr_backup":   cbr.DataSourceBackup(),
			"hcso_cbr_vaults":   cbr.DataSourceVaults(),
			"hcso_cbr_policies": cbr.DataSourcePolicies(),

			"hcso_cce_addon_template": cce.DataSourceAddonTemplate(),
			"hcso_cce_cluster":        cce.DataSourceCCEClusterV3(),
			"hcso_cce_clusters":       cce.DataSourceCCEClusters(),
			"hcso_cce_node":           cce.DataSourceNode(),
			"hcso_cce_node_pool":      cce.DataSourceCCENodePoolV3(),
			"hcso_cce_nodes":          cce.DataSourceNodes(),

			"hcso_dcs_maintainwindow": dcs.DataSourceDcsMaintainWindow(),

			"hcso_dns_zones": dns.DataSourceZones(),

			"hcso_compute_flavors":      ecs.DataSourceEcsFlavors(),
			"hcso_compute_instance":     ecs.DataSourceComputeInstance(),
			"hcso_compute_instances":    ecs.DataSourceComputeInstances(),
			"hcso_compute_servergroups": ecs.DataSourceComputeServerGroups(),

			"hcso_evs_volumes": evs.DataSourceEvsVolumesV2(),

			"hcso_elb_certificate":       elb.DataSourceELBCertificateV3(),
			"hcso_elb_flavors":           elb.DataSourceElbFlavorsV3(),
			"hcso_elb_ipgroups":          elb.DataSourceElbIpGroups(),
			"hcso_elb_l7policies":        elb.DataSourceElbL7policies(),
			"hcso_elb_l7rules":           elb.DataSourceElbL7rules(),
			"hcso_elb_listeners":         elb.DataSourceElbListeners(),
			"hcso_elb_loadbalancers":     elb.DataSourceElbLoadbalances(),
			"hcso_elb_logtanks":          elb.DataSourceElbLogtanks(),
			"hcso_elb_members":           elb.DataSourceElbMembers(),
			"hcso_elb_pools":             elb.DataSourcePools(),
			"hcso_elb_security_policies": elb.DataSourceElbSecurityPolicies(),

			"hcso_images_image":  ims.DataSourceImagesImageV2(),
			"hcso_images_images": ims.DataSourceImagesImages(),

			"hcso_lb_listeners":    lb.DataSourceListeners(),
			"hcso_lb_loadbalancer": lb.DataSourceELBV2Loadbalancer(),

			"hcso_mapreduce_clusters": mrs.DataSourceMrsClusters(),

			"hcso_rds_backups":              rds.DataSourceBackup(),
			"hcso_rds_engine_versions":      rds.DataSourceRdsEngineVersionsV3(),
			"hcso_rds_flavors":              rds.DataSourceRdsFlavor(),
			"hcso_rds_instances":            rds.DataSourceRdsInstances(),
			"hcso_rds_mysql_accounts":       rds.DataSourceRdsMysqlAccounts(),
			"hcso_rds_mysql_databases":      rds.DataSourceRdsMysqlDatabases(),
			"hcso_rds_parametergroups":      rds.DataSourceParametergroups(),
			"hcso_rds_pg_plugins":           rds.DataSourcePgPlugins(),
			"hcso_rds_sqlserver_collations": rds.DataSourceSQLServerCollations(),
			"hcso_rds_storage_types":        rds.DataSourceStoragetype(),

			"hcso_sfs_turbos": sfs.DataSourceTurbos(),

			"hcso_vpc":                    vpc.DataSourceVpcV1(),
			"hcso_vpcs":                   vpc.DataSourceVpcs(),
			"hcso_vpc_subnet_ids":         vpc.DataSourceVpcSubnetIdsV1(),
			"hcso_vpc_subnets":            vpc.DataSourceVpcSubnets(),
			"hcso_networking_port":        vpc.DataSourceNetworkingPortV2(),
			"hcso_vpc_peering_connection": vpc.DataSourceVpcPeeringConnectionV2(),
			"hcso_networking_secgroups":   vpc.DataSourceNetworkingSecGroups(),

			"hcso_vpc_bandwidth": eip.DataSourceBandWidth(),
			"hcso_vpc_eip":       eip.DataSourceVpcEip(),
			"hcso_vpc_eips":      eip.DataSourceVpcEips(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"hcso_antiddos_basic": antiddos.ResourceCloudNativeAntiDdos(),

			"hcso_cbr_backup_share": cbr.ResourceBackupShare(),
			"hcso_cbr_checkpoint":   cbr.ResourceCheckpoint(),
			"hcso_cbr_policy":       cbr.ResourcePolicy(),
			"hcso_cbr_vault":        cbr.ResourceVault(),

			"hcso_cce_cluster":     cce.ResourceCluster(),
			"hcso_cce_node":        cce.ResourceNode(),
			"hcso_cce_node_attach": cce.ResourceNodeAttach(),
			"hcso_cce_addon":       cce.ResourceAddon(),
			"hcso_cce_node_pool":   cce.ResourceNodePool(),

			"hcso_compute_instance":         ecs.ResourceComputeInstance(),
			"hcso_compute_servergroup":      ecs.ResourceComputeServerGroup(),
			"hcso_compute_interface_attach": ecs.ResourceComputeInterfaceAttach(),
			"hcso_compute_keypair":          huaweicloud.ResourceComputeKeypairV2(),
			"hcso_compute_eip_associate":    ecs.ResourceComputeEIPAssociate(),
			"hcso_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),

			"hcso_elb_certificate":     elb.ResourceCertificateV3(),
			"hcso_elb_ipgroup":         elb.ResourceIpGroupV3(),
			"hcso_elb_l7policy":        elb.ResourceL7PolicyV3(),
			"hcso_elb_l7rule":          elb.ResourceL7RuleV3(),
			"hcso_elb_listener":        elb.ResourceListenerV3(),
			"hcso_elb_loadbalancer":    elb.ResourceLoadBalancerV3(),
			"hcso_elb_logtank":         elb.ResourceLogTank(),
			"hcso_elb_member":          elb.ResourceMemberV3(),
			"hcso_elb_monitor":         elb.ResourceMonitorV3(),
			"hcso_elb_pool":            elb.ResourcePoolV3(),
			"hcso_elb_security_policy": elb.ResourceSecurityPolicy(),
			"hcso_lts_group":           lts.ResourceLTSGroup(),
			"hcso_lts_stream":          lts.ResourceLTSStream(),

			"hcso_evs_volume":   evs.ResourceEvsVolume(),
			"hcso_evs_snapshot": evs.ResourceEvsSnapshotV2(),

			"hcso_images_image":                ims.ResourceImsImage(),
			"hcso_images_image_copy":           ims.ResourceImsImageCopy(),
			"hcso_images_image_share":          ims.ResourceImsImageShare(),
			"hcso_images_image_share_accepter": ims.ResourceImsImageShareAccepter(),

			"hcso_lb_l7policy":     lb.ResourceL7PolicyV2(),
			"hcso_lb_l7rule":       lb.ResourceL7RuleV2(),
			"hcso_lb_loadbalancer": lb.ResourceLoadBalancer(),
			"hcso_lb_listener":     lb.ResourceListener(),
			"hcso_lb_member":       lb.ResourceMemberV2(),
			"hcso_lb_monitor":      lb.ResourceMonitorV2(),
			"hcso_lb_pool":         lb.ResourcePoolV2(),
			"hcso_lb_whitelist":    lb.ResourceWhitelistV2(),

			"hcso_mapreduce_cluster": mrs.ResourceMRSClusterV2(),
			"hcso_mapreduce_job":     mrs.ResourceMRSJobV2(),

			"hcso_nat_private_gateway":    nat.ResourcePrivateGateway(),
			"hcso_nat_private_snat_rule":  nat.ResourcePrivateSnatRule(),
			"hcso_nat_private_transit_ip": nat.ResourcePrivateTransitIp(),

			"hcso_rds_backup":                       rds.ResourceBackup(),
			"hcso_rds_cross_region_backup_strategy": rds.ResourceBackupStrategy(),
			"hcso_rds_instance":                     rds.ResourceRdsInstance(),
			"hcso_rds_mysql_account":                rds.ResourceMysqlAccount(),
			"hcso_rds_mysql_binlog":                 rds.ResourceMysqlBinlog(),
			"hcso_rds_mysql_database_privilege":     rds.ResourceMysqlDatabasePrivilege(),
			"hcso_rds_mysql_database":               rds.ResourceMysqlDatabase(),
			"hcso_rds_parametergroup":               rds.ResourceRdsConfiguration(),
			"hcso_rds_pg_account":                   rds.ResourcePgAccount(),
			"hcso_rds_pg_database":                  rds.ResourcePgDatabase(),
			"hcso_rds_pg_plugin":                    rds.ResourceRdsPgPlugin(),
			"hcso_rds_read_replica_instance":        rds.ResourceRdsReadReplicaInstance(),
			"hcso_rds_sql_audit":                    rds.ResourceSQLAudit(),
			"hcso_rds_sqlserver_account":            rds.ResourceSQLServerAccount(),
			"hcso_rds_sqlserver_database_privilege": rds.ResourceSQLServerDatabasePrivilege(),
			"hcso_rds_sqlserver_database":           rds.ResourceSQLServerDatabase(),

			"hcso_sfs_turbo": sfs.ResourceSFSTurbo(),

			"hcso_tms_tags": tms.ResourceTmsTag(),

			"hcso_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"hcso_vpc_address_group":               vpc.ResourceVpcAddressGroup(),
			"hcso_vpc_subnet":                      vpc.ResourceVpcSubnetV1(),
			"hcso_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"hcso_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"hcso_networking_secgroup":             vpc.ResourceNetworkingSecGroup(),
			"hcso_networking_secgroup_rule":        vpc.ResourceNetworkingSecGroupRule(),
			"hcso_networking_vip":                  vpc.ResourceNetworkingVip(),

			"hcso_vpc_bandwidth":            eip.ResourceVpcBandWidthV1(),
			"hcso_vpc_bandwidth_associate":  eip.ResourceBandWidthAssociate(),
			"hcso_vpc_eip":                  eip.ResourceVpcEIPV1(),
			"hcso_vpc_eip_associate":        eip.ResourceEIPAssociate(),
			"hcso_networking_eip_associate": eip.ResourceEIPAssociate(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11 cc
			terraformVersion = "0.11+compatible"
		}

		return configureProvider(ctx, d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The HuaweiCloud region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"project_id": "The ID of the project to login with.",

		"project_name": "The name of the project to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"domain_id": "The ID of the Domain to scope to.",

		"domain_name": "The name of the Domain to scope to.",

		"access_key":     "The access key of the HuaweiCloud to use.",
		"secret_key":     "The secret key of the HuaweiCloud to use.",
		"security_token": "The security token to authenticate with a temporary security credential.",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"agency_name": "The name of agency",

		"agency_domain_name": "The name of domain who created the agency (Identity v3).",

		"delegated_project": "The name of delegated project (Identity v3).",

		"assume_role_agency_name": "The name of agency for assume role.",

		"assume_role_domain_name": "The name of domain for assume role.",

		"cloud": "The endpoint of cloud provider, defaults to myhuaweicloud.com",

		"endpoints": "The custom endpoints used to override the default endpoint URL.",

		"regional": "Whether the service endpoints are regional",

		"shared_config_file": "The path to the shared config file. If not set, the default is ~/.hcloud/config.json.",

		"profile": "The profile name as set in the shared config file.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"enterprise_project_id": "enterprise project id",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData, terraformVersion string) (interface{},
	diag.Diagnostics) {
	var tenantName, tenantID, delegatedProject, identityEndpoint string
	region := d.Get("region").(string)
	isRegional := d.Get("regional").(bool)
	cloud := getCloudDomain(d.Get("cloud").(string))

	// project_name is prior to tenant_name
	// if neither of them was set, use region as the default project
	if v, ok := d.GetOk("project_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else if v, ok := d.GetOk("tenant_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else {
		tenantName = region
	}

	// project_id is prior to tenant_id
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		tenantID = v.(string)
	} else {
		tenantID = d.Get("tenant_id").(string)
	}

	// Use region as delegated_project if it's not set
	if v, ok := d.GetOk("delegated_project"); ok && v.(string) != "" {
		delegatedProject = v.(string)
	} else {
		delegatedProject = region
	}

	// use auth_url as identityEndpoint if specified
	if v, ok := d.GetOk("auth_url"); ok {
		identityEndpoint = v.(string)
	} else {
		// use cloud as basis for identityEndpoint
		identityEndpoint = fmt.Sprintf("https://iam.%s.%s/v3", region, cloud)
	}

	hcsoConfig := hcso_config.HCSOConfig{
		Config: config.Config{
			AccessKey:           d.Get("access_key").(string),
			SecretKey:           d.Get("secret_key").(string),
			CACertFile:          d.Get("cacert_file").(string),
			ClientCertFile:      d.Get("cert").(string),
			ClientKeyFile:       d.Get("key").(string),
			DomainID:            d.Get("domain_id").(string),
			DomainName:          d.Get("domain_name").(string),
			IdentityEndpoint:    identityEndpoint,
			Insecure:            d.Get("insecure").(bool),
			Password:            d.Get("password").(string),
			Token:               d.Get("token").(string),
			SecurityToken:       d.Get("security_token").(string),
			Region:              region,
			TenantID:            tenantID,
			TenantName:          tenantName,
			Username:            d.Get("user_name").(string),
			UserID:              d.Get("user_id").(string),
			AgencyName:          d.Get("agency_name").(string),
			AgencyDomainName:    d.Get("agency_domain_name").(string),
			DelegatedProject:    delegatedProject,
			Cloud:               cloud,
			RegionClient:        isRegional,
			MaxRetries:          d.Get("max_retries").(int),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
			SharedConfigFile:    d.Get("shared_config_file").(string),
			Profile:             d.Get("profile").(string),
			TerraformVersion:    terraformVersion,
			RegionProjectIDMap:  make(map[string]string),
			RPLock:              new(sync.Mutex),
			SecurityKeyLock:     new(sync.Mutex),
		},
	}

	hcsoConfig.Metadata = &hcsoConfig.Config

	// get assume role
	assumeRoleList := d.Get("assume_role").([]interface{})
	if len(assumeRoleList) == 0 {
		// without assume_role block in provider
		delegatedAgencyName := os.Getenv("HCSO_ASSUME_ROLE_AGENCY_NAME")
		delegatedDomianName := os.Getenv("HCSO_ASSUME_ROLE_DOMAIN_NAME")
		if delegatedAgencyName != "" && delegatedDomianName != "" {
			hcsoConfig.AssumeRoleAgency = delegatedAgencyName
			hcsoConfig.AssumeRoleDomain = delegatedDomianName
		}
	} else {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		hcsoConfig.AssumeRoleAgency = assumeRole["agency_name"].(string)
		hcsoConfig.AssumeRoleDomain = assumeRole["domain_name"].(string)
	}

	// get custom endpoints
	endpoints, err := flattenProviderEndpoints(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	hcsoConfig.Endpoints = endpoints

	if err := hcsoConfig.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &hcsoConfig.Config, nil
}

func flattenProviderEndpoints(d *schema.ResourceData) (map[string]string, error) {
	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		endpoint := strings.TrimSpace(val.(string))
		// check empty string
		if endpoint == "" {
			return nil, fmt.Errorf("the value of customer endpoint %s must be specified", key)
		}

		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}
		epMap[key] = endpoint
	}

	// unify the endpoint which has multiple versions
	for key := range endpoints {
		ep, ok := epMap[key]
		if !ok {
			continue
		}

		multiKeys := config.GetServiceDerivedCatalogKeys(key)
		for _, k := range multiKeys {
			epMap[k] = ep
		}
	}

	log.Printf("[DEBUG] customer endpoints: %+v", epMap)
	return epMap, nil
}

func getCloudDomain(cloud string) string {
	// first, use the specified value
	if cloud != "" {
		return cloud
	}

	return defaultCloud
}
