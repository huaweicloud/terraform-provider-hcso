//nolint:revive
package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcso/internal"
)

var (
	HCSO_REGION_NAME                        = os.Getenv("HCSO_REGION_NAME")
	HCSO_CUSTOM_REGION_NAME                 = os.Getenv("HCSO_CUSTOM_REGION_NAME")
	HCSO_AVAILABILITY_ZONE                  = os.Getenv("HCSO_AVAILABILITY_ZONE")
	HCSO_ACCESS_KEY                         = os.Getenv("HCSO_ACCESS_KEY")
	HCSO_SECRET_KEY                         = os.Getenv("HCSO_SECRET_KEY")
	HCSO_USER_ID                            = os.Getenv("HCSO_USER_ID")
	HCSO_USER_NAME                          = os.Getenv("HCSO_USER_NAME")
	HCSO_PROJECT_ID                         = os.Getenv("HCSO_PROJECT_ID")
	HCSO_DOMAIN_ID                          = os.Getenv("HCSO_DOMAIN_ID")
	HCSO_DOMAIN_NAME                        = os.Getenv("HCSO_DOMAIN_NAME")
	HCSO_ENTERPRISE_PROJECT_ID_TEST         = os.Getenv("HCSO_ENTERPRISE_PROJECT_ID_TEST")
	HCSO_ENTERPRISE_MIGRATE_PROJECT_ID_TEST = os.Getenv("HCSO_ENTERPRISE_MIGRATE_PROJECT_ID_TEST")

	HCSO_FLAVOR_ID             = os.Getenv("HCSO_FLAVOR_ID")
	HCSO_FLAVOR_NAME           = os.Getenv("HCSO_FLAVOR_NAME")
	HCSO_IMAGE_ID              = os.Getenv("HCSO_IMAGE_ID")
	HCSO_IMAGE_NAME            = os.Getenv("HCSO_IMAGE_NAME")
	HCSO_VPC_ID                = os.Getenv("HCSO_VPC_ID")
	HCSO_NETWORK_ID            = os.Getenv("HCSO_NETWORK_ID")
	HCSO_SUBNET_ID             = os.Getenv("HCSO_SUBNET_ID")
	HCSO_ENTERPRISE_PROJECT_ID = os.Getenv("HCSO_ENTERPRISE_PROJECT_ID")
	HCSO_ADMIN                 = os.Getenv("HCSO_ADMIN")

	HCSO_MAPREDUCE_CUSTOM           = os.Getenv("HCSO_MAPREDUCE_CUSTOM")
	HCSO_MAPREDUCE_BOOTSTRAP_SCRIPT = os.Getenv("HCSO_MAPREDUCE_BOOTSTRAP_SCRIPT")

	HCSO_CNAD_ENABLE_FLAG       = os.Getenv("HCSO_CNAD_ENABLE_FLAG")
	HCSO_CNAD_PROJECT_OBJECT_ID = os.Getenv("HCSO_CNAD_PROJECT_OBJECT_ID")

	HCSO_OBS_BUCKET_NAME        = os.Getenv("HCSO_OBS_BUCKET_NAME")
	HCSO_OBS_DESTINATION_BUCKET = os.Getenv("HCSO_OBS_DESTINATION_BUCKET")

	HCSO_OMS_ENABLE_FLAG = os.Getenv("HCSO_OMS_ENABLE_FLAG")

	HCSO_DEPRECATED_ENVIRONMENT = os.Getenv("HCSO_DEPRECATED_ENVIRONMENT")
	HCSO_INTERNAL_USED          = os.Getenv("HCSO_INTERNAL_USED")

	HCSO_WAF_ENABLE_FLAG = os.Getenv("HCSO_WAF_ENABLE_FLAG")

	HCSO_DEST_REGION          = os.Getenv("HCSO_DEST_REGION")
	HCSO_DEST_PROJECT_ID      = os.Getenv("HCSO_DEST_PROJECT_ID")
	HCSO_DEST_PROJECT_ID_TEST = os.Getenv("HCSO_DEST_PROJECT_ID_TEST")
	HCSO_CHARGING_MODE        = os.Getenv("HCSO_CHARGING_MODE")
	HCSO_HIGH_COST_ALLOW      = os.Getenv("HCSO_HIGH_COST_ALLOW")
	HCSO_SWR_SHARING_ACCOUNT  = os.Getenv("HCSO_SWR_SHARING_ACCOUNT")

	HCSO_RAM_SHARE_ACCOUNT_ID          = os.Getenv("HCSO_RAM_SHARE_ACCOUNT_ID")
	HCSO_RAM_SHARE_RESOURCE_URN        = os.Getenv("HCSO_RAM_SHARE_RESOURCE_URN")
	HCSO_RAM_SHARE_UPDATE_ACCOUNT_ID   = os.Getenv("HCSO_RAM_SHARE_UPDATE_ACCOUNT_ID")
	HCSO_RAM_SHARE_UPDATE_RESOURCE_URN = os.Getenv("HCSO_RAM_SHARE_UPDATE_RESOURCE_URN")

	HCSO_CERTIFICATE_KEY_PATH         = os.Getenv("HCSO_CERTIFICATE_KEY_PATH")
	HCSO_CERTIFICATE_CHAIN_PATH       = os.Getenv("HCSO_CERTIFICATE_CHAIN_PATH")
	HCSO_CERTIFICATE_PRIVATE_KEY_PATH = os.Getenv("HCSO_CERTIFICATE_PRIVATE_KEY_PATH")
	HCSO_CERTIFICATE_SERVICE          = os.Getenv("HCSO_CERTIFICATE_SERVICE")
	HCSO_CERTIFICATE_PROJECT          = os.Getenv("HCSO_CERTIFICATE_PROJECT")
	HCSO_CERTIFICATE_PROJECT_UPDATED  = os.Getenv("HCSO_CERTIFICATE_PROJECT_UPDATED")
	HCSO_CERTIFICATE_NAME             = os.Getenv("HCSO_CERTIFICATE_NAME")
	HCSO_DMS_ENVIRONMENT              = os.Getenv("HCSO_DMS_ENVIRONMENT")
	HCSO_SMS_SOURCE_SERVER            = os.Getenv("HCSO_SMS_SOURCE_SERVER")

	HCSO_DLI_FLINK_JAR_OBS_PATH           = os.Getenv("HCSO_DLI_FLINK_JAR_OBS_PATH")
	HCSO_DLI_DS_AUTH_CSS_OBS_PATH         = os.Getenv("HCSO_DLI_DS_AUTH_CSS_OBS_PATH")
	HCSO_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH = os.Getenv("HCSO_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH")
	HCSO_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH   = os.Getenv("HCSO_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH")
	HCSO_DLI_DS_AUTH_KRB_CONF_OBS_PATH    = os.Getenv("HCSO_DLI_DS_AUTH_KRB_CONF_OBS_PATH")
	HCSO_DLI_DS_AUTH_KRB_TAB_OBS_PATH     = os.Getenv("HCSO_DLI_DS_AUTH_KRB_TAB_OBS_PATH")
	HCSO_DLI_AGENCY_FLAG                  = os.Getenv("HCSO_DLI_AGENCY_FLAG")

	HCSO_GITHUB_REPO_HOST        = os.Getenv("HCSO_GITHUB_REPO_HOST")        // Repository host (Github, Gitlab, Gitee)
	HCSO_GITHUB_PERSONAL_TOKEN   = os.Getenv("HCSO_GITHUB_PERSONAL_TOKEN")   // Personal access token (Github, Gitlab, Gitee)
	HCSO_GITHUB_REPO_PWD         = os.Getenv("HCSO_GITHUB_REPO_PWD")         // Repository password (DevCloud, BitBucket)
	HCSO_GITHUB_REPO_URL         = os.Getenv("HCSO_GITHUB_REPO_URL")         // Repository URL (Github, Gitlab, Gitee)
	HCSO_OBS_STORAGE_URL         = os.Getenv("HCSO_OBS_STORAGE_URL")         // OBS storage URL where ZIP file is located
	HCSO_BUILD_IMAGE_URL         = os.Getenv("HCSO_BUILD_IMAGE_URL")         // SWR Image URL for component deployment
	HCSO_BUILD_IMAGE_URL_UPDATED = os.Getenv("HCSO_BUILD_IMAGE_URL_UPDATED") // SWR Image URL for component deployment update

	HCSO_VOD_WATERMARK_FILE   = os.Getenv("HCSO_VOD_WATERMARK_FILE")
	HCSO_VOD_MEDIA_ASSET_FILE = os.Getenv("HCSO_VOD_MEDIA_ASSET_FILE")

	HCSO_CHAIR_EMAIL              = os.Getenv("HCSO_CHAIR_EMAIL")
	HCSO_GUEST_EMAIL              = os.Getenv("HCSO_GUEST_EMAIL")
	HCSO_MEETING_ACCOUNT_NAME     = os.Getenv("HCSO_MEETING_ACCOUNT_NAME")
	HCSO_MEETING_ACCOUNT_PASSWORD = os.Getenv("HCSO_MEETING_ACCOUNT_PASSWORD")
	HCSO_MEETING_APP_ID           = os.Getenv("HCSO_MEETING_APP_ID")
	HCSO_MEETING_APP_KEY          = os.Getenv("HCSO_MEETING_APP_KEY")
	HCSO_MEETING_USER_ID          = os.Getenv("HCSO_MEETING_USER_ID")
	HCSO_MEETING_ROOM_ID          = os.Getenv("HCSO_MEETING_ROOM_ID")

	HCSO_AAD_INSTANCE_ID = os.Getenv("HCSO_AAD_INSTANCE_ID")
	HCSO_AAD_IP_ADDRESS  = os.Getenv("HCSO_AAD_IP_ADDRESS")

	HCSO_WORKSPACE_AD_DOMAIN_NAME = os.Getenv("HCSO_WORKSPACE_AD_DOMAIN_NAME") // Domain name, e.g. "example.com".
	HCSO_WORKSPACE_AD_SERVER_PWD  = os.Getenv("HCSO_WORKSPACE_AD_SERVER_PWD")  // The password of AD server.
	HCSO_WORKSPACE_AD_DOMAIN_IP   = os.Getenv("HCSO_WORKSPACE_AD_DOMAIN_IP")   // Active domain IP, e.g. "192.168.196.3".
	HCSO_WORKSPACE_AD_VPC_ID      = os.Getenv("HCSO_WORKSPACE_AD_VPC_ID")      // The VPC ID to which the AD server and desktops belongs.
	HCSO_WORKSPACE_AD_NETWORK_ID  = os.Getenv("HCSO_WORKSPACE_AD_NETWORK_ID")  // The network ID to which the AD server belongs.

	HCSO_FGS_TRIGGER_LTS_AGENCY = os.Getenv("HCSO_FGS_TRIGGER_LTS_AGENCY")

	HCSO_KMS_ENVIRONMENT = os.Getenv("HCSO_KMS_ENVIRONMENT")

	HCSO_MULTI_ACCOUNT_ENVIRONMENT            = os.Getenv("HCSO_MULTI_ACCOUNT_ENVIRONMENT")
	HCSO_ORGANIZATIONS_ACCOUNT_NAME           = os.Getenv("HCSO_ORGANIZATIONS_ACCOUNT_NAME")
	HCSO_ORGANIZATIONS_INVITE_ACCOUNT_ID      = os.Getenv("HCSO_ORGANIZATIONS_INVITE_ACCOUNT_ID")
	HCSO_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID = os.Getenv("HCSO_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID")
	HCSO_ORGANIZATIONS_INVITATION_ID          = os.Getenv("HCSO_ORGANIZATIONS_INVITATION_ID")

	HCSO_IDENTITY_CENTER_ACCOUNT_ID = os.Getenv("HCSO_IDENTITY_CENTER_ACCOUNT_ID")

	HCSO_ER_TEST_ON = os.Getenv("HCSO_ER_TEST_ON") // Whether to run the ER related tests.

	// The OBS address where the HCL/JSON template archive (No variables) is located.
	HCSO_RF_TEMPLATE_ARCHIVE_NO_VARS_URI = os.Getenv("HCSO_RF_TEMPLATE_ARCHIVE_NO_VARS_URI")
	// The OBS address where the HCL/JSON template archive is located.
	HCSO_RF_TEMPLATE_ARCHIVE_URI = os.Getenv("HCSO_RF_TEMPLATE_ARCHIVE_URI")
	// The OBS address where the variable archive corresponding to the HCL/JSON template is located.
	HCSO_RF_VARIABLES_ARCHIVE_URI = os.Getenv("HCSO_RF_VARIABLES_ARCHIVE_URI")

	// The direct connection ID (provider does not support direct connection resource).
	HCSO_DC_DIRECT_CONNECT_ID = os.Getenv("HCSO_DC_DIRECT_CONNECT_ID")

	// The CFW instance ID
	HCSO_CFW_INSTANCE_ID = os.Getenv("HCSO_CFW_INSTANCE_ID")

	// The cluster ID of the CCE
	HCSO_CCE_CLUSTER_ID = os.Getenv("HCSO_CCE_CLUSTER_ID")
	// The partition az of the CCE
	HCSO_CCE_PARTITION_AZ = os.Getenv("HCSO_CCE_PARTITION_AZ")
	// The namespace of the workload is located
	HCSO_WORKLOAD_NAMESPACE = os.Getenv("HCSO_WORKLOAD_NAMESPACE")
	// The workload type deployed in CCE/CCI
	HCSO_WORKLOAD_TYPE = os.Getenv("HCSO_WORKLOAD_TYPE")
	// The workload name deployed in CCE/CCI
	HCSO_WORKLOAD_NAME = os.Getenv("HCSO_WORKLOAD_NAME")
	// The target region of SWR image auto sync
	HCSO_SWR_TARGET_REGION = os.Getenv("HCSO_SWR_TARGET_REGION")
	// The target organization of SWR image auto sync
	HCSO_SWR_TARGET_ORGANIZATION = os.Getenv("HCSO_SWR_TARGET_ORGANIZATION")

	// The ID of the CBR backup
	HCSO_IMS_BACKUP_ID = os.Getenv("HCSO_IMS_BACKUP_ID")

	// The SecMaster workspace ID
	HCSO_SECMASTER_WORKSPACE_ID = os.Getenv("HCSO_SECMASTER_WORKSPACE_ID")

	HCSO_MODELARTS_HAS_SUBSCRIBE_MODEL = os.Getenv("HCSO_MODELARTS_HAS_SUBSCRIBE_MODEL")

	// Deprecated
	HCSO_SRC_ACCESS_KEY = os.Getenv("HCSO_SRC_ACCESS_KEY")
	HCSO_SRC_SECRET_KEY = os.Getenv("HCSO_SRC_SECRET_KEY")
	HCSO_EXTGW_ID       = os.Getenv("HCSO_EXTGW_ID")
	HCSO_POOL_NAME      = os.Getenv("HCSO_POOL_NAME")

	HCSO_IMAGE_SHARE_SOURCE_IMAGE_ID = os.Getenv("HCSO_IMAGE_SHARE_SOURCE_IMAGE_ID")

	HCSO_CERTIFICATE_CONTENT         = os.Getenv("HCSO_CERTIFICATE_CONTENT")
	HCSO_CERTIFICATE_PRIVATE_KEY     = os.Getenv("HCSO_CERTIFICATE_PRIVATE_KEY")
	HCSO_CERTIFICATE_ROOT_CA         = os.Getenv("HCSO_CERTIFICATE_ROOT_CA")
	HCSO_NEW_CERTIFICATE_CONTENT     = os.Getenv("HCSO_NEW_CERTIFICATE_CONTENT")
	HCSO_NEW_CERTIFICATE_PRIVATE_KEY = os.Getenv("HCSO_NEW_CERTIFICATE_PRIVATE_KEY")
	HCSO_NEW_CERTIFICATE_ROOT_CA     = os.Getenv("HCSO_NEW_CERTIFICATE_ROOT_CA")

	HCSO_CODEARTS_RESOURCE_POOL_ID = os.Getenv("HCSO_CODEARTS_RESOURCE_POOL_ID")
	HCSO_CODEARTS_ENABLE_FLAG      = os.Getenv("HCSO_CODEARTS_ENABLE_FLAG")

	HCSO_EG_CHANNEL_ID = os.Getenv("HCSO_EG_CHANNEL_ID")
)

// TestAccProviders is a static map containing only the main provider instance.
//
// Deprecated: Terraform Plugin SDK version 2 uses TestCase.ProviderFactories
// but supports this value in TestCase.Providers for backwards compatibility.
// In the future Providers: TestAccProviders will be changed to
// ProviderFactories: TestAccProviderFactories
var TestAccProviders map[string]*schema.Provider

// TestAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// TestAccProvider is the "main" provider instance
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = internal.Provider()

	TestAccProviders = map[string]*schema.Provider{
		"hcso": TestAccProvider,
	}

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"hcso": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func preCheckRequiredEnvVars(t *testing.T) {
	if HCSO_REGION_NAME == "" {
		t.Fatal("HCSO_REGION_NAME must be set for acceptance tests")
	}
}

// use this function to precheck langding zone services, such as Organizations and Identity Center
// lintignore:AT003
func TestAccPreCheckMultiAccount(t *testing.T) {
	if HCSO_MULTI_ACCOUNT_ENVIRONMENT == "" {
		t.Skip("This environment does not support multi-account tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsAccountName(t *testing.T) {
	if HCSO_ORGANIZATIONS_ACCOUNT_NAME == "" {
		t.Skip("HCSO_ORGANIZATIONS_ACCOUNT_NAME must be set for the acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsInviteAccountId(t *testing.T) {
	if HCSO_ORGANIZATIONS_INVITE_ACCOUNT_ID == "" {
		t.Skip("HCSO_ORGANIZATIONS_INVITE_ACCOUNT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsInvitationId(t *testing.T) {
	if HCSO_ORGANIZATIONS_INVITATION_ID == "" {
		t.Skip("HCSO_ORGANIZATIONS_INVITATION_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsOrganizationalUnitId(t *testing.T) {
	if HCSO_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID == "" {
		t.Skip("HCSO_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckIdentityCenterAccountId(t *testing.T) {
	if HCSO_IDENTITY_CENTER_ACCOUNT_ID == "" {
		t.Skip("HCSO_IDENTITY_CENTER_ACCOUNT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HCSO_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

// lintignore:AT003
func TestAccPreCheckDestProjectIds(t *testing.T) {
	if HCSO_DEST_PROJECT_ID == "" || HCSO_DEST_PROJECT_ID_TEST == "" {
		t.Skip("HCSO_DEST_PROJECT_ID and HCSO_DEST_PROJECT_ID_TEST must be set for acceptance test.")
	}
}

// lintignore:AT003
func TestAccPrecheckDomainId(t *testing.T) {
	if HCSO_DOMAIN_ID == "" {
		t.Skip("HCSO_DOMAIN_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPrecheckDomainName(t *testing.T) {
	if HCSO_DOMAIN_NAME == "" {
		t.Skip("HCSO_DOMAIN_NAME must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPrecheckCustomRegion(t *testing.T) {
	if HCSO_CUSTOM_REGION_NAME == "" {
		t.Skip("HCSO_CUSTOM_REGION_NAME must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckDeprecated(t *testing.T) {
	if HCSO_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

// lintignore:AT003
func TestAccPreCheckInternal(t *testing.T) {
	if HCSO_INTERNAL_USED == "" {
		t.Skip("HCSO_INTERNAL_USED must be set for internal acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckEpsID(t *testing.T) {
	// The environment variables in tests take HCSO_ENTERPRISE_PROJECT_ID_TEST instead of HCSO_ENTERPRISE_PROJECT_ID to
	// ensure that other data-resources that support enterprise projects query the default project without being
	// affected by this variable.
	if HCSO_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Enterprise Project ID for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMigrateEpsID(t *testing.T) {
	if HCSO_ENTERPRISE_PROJECT_ID_TEST == "" || HCSO_ENTERPRISE_MIGRATE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Migrate Enterprise Project ID for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckUserId(t *testing.T) {
	if HCSO_USER_ID == "" {
		t.Skip("The environment variables does not support the user ID (HCSO_USER_ID) for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSms(t *testing.T) {
	if HCSO_SMS_SOURCE_SERVER == "" {
		t.Skip("HCSO_SMS_SOURCE_SERVER must be set for SMS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMrsCustom(t *testing.T) {
	if HCSO_MAPREDUCE_CUSTOM == "" {
		t.Skip("HCSO_MAPREDUCE_CUSTOM must be set for acceptance tests:custom type cluster of map reduce")
	}
}

// lintignore:AT003
func TestAccPreCheckMrsBootstrapScript(t *testing.T) {
	if HCSO_MAPREDUCE_BOOTSTRAP_SCRIPT == "" {
		t.Skip("HCSO_MAPREDUCE_BOOTSTRAP_SCRIPT must be set for acceptance tests: cluster of map reduce with bootstrap")
	}
}

// lintignore:AT003
func TestAccPreCheckFgsTrigger(t *testing.T) {
	if HCSO_FGS_TRIGGER_LTS_AGENCY == "" {
		t.Skip("HCSO_FGS_TRIGGER_LTS_AGENCY must be set for FGS trigger acceptance tests")
	}
}

// Deprecated
// lintignore:AT003
func TestAccPreCheckMaas(t *testing.T) {
	if HCSO_ACCESS_KEY == "" || HCSO_SECRET_KEY == "" || HCSO_SRC_ACCESS_KEY == "" || HCSO_SRC_SECRET_KEY == "" {
		t.Skip("HCSO_ACCESS_KEY, HCSO_SECRET_KEY, HCSO_SRC_ACCESS_KEY, and HCSO_SRC_SECRET_KEY  must be set for MAAS acceptance tests")
	}
}

func RandomAccResourceName() string {
	return fmt.Sprintf("tf_test_%s", acctest.RandString(5))
}

func RandomAccResourceNameWithDash() string {
	return fmt.Sprintf("tf-test-%s", acctest.RandString(5))
}

func RandomCidr() string {
	return fmt.Sprintf("172.16.%d.0/24", acctest.RandIntRange(0, 255))
}

func RandomCidrAndGatewayIp() (string, string) {
	seed := acctest.RandIntRange(0, 255)
	return fmt.Sprintf("172.16.%d.0/24", seed), fmt.Sprintf("172.16.%d.1", seed)
}

func RandomPassword() string {
	return fmt.Sprintf("%s%s%s%d",
		acctest.RandStringFromCharSet(2, "ABCDEFGHIJKLMNOPQRSTUVWXZY"),
		acctest.RandStringFromCharSet(3, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(2, "~!@#%^*-_=+?"),
		acctest.RandIntRange(1000, 9999))
}

// lintignore:AT003
func TestAccPrecheckWafInstance(t *testing.T) {
	if HCSO_WAF_ENABLE_FLAG == "" {
		t.Skip("Skip the WAF acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckCNADInstance(t *testing.T) {
	if HCSO_CNAD_ENABLE_FLAG == "" {
		t.Skip("Skip the CNAD acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckCNADProtectedObject(t *testing.T) {
	if HCSO_CNAD_PROJECT_OBJECT_ID == "" {
		t.Skip("Skipping test because HCSO_CNAD_PROJECT_OBJECT_ID is required for this acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckOmsInstance(t *testing.T) {
	if HCSO_OMS_ENABLE_FLAG == "" {
		t.Skip("Skip the OMS acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckAdminOnly(t *testing.T) {
	if HCSO_ADMIN == "" {
		t.Skip("Skipping test because it requires the admin privileges")
	}
}

// lintignore:AT003
func TestAccPreCheckReplication(t *testing.T) {
	if HCSO_DEST_REGION == "" || HCSO_DEST_PROJECT_ID == "" {
		t.Skip("Skip the replication policy acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckProjectId(t *testing.T) {
	if HCSO_DEST_PROJECT_ID_TEST == "" {
		t.Skip("Skipping test because it requires the test project id.")
	}
}

// lintignore:AT003
func TestAccPreCheckProject(t *testing.T) {
	if HCSO_ENTERPRISE_PROJECT_ID_TEST != "" {
		t.Skip("This environment does not support project tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBS(t *testing.T) {
	if HCSO_ACCESS_KEY == "" || HCSO_SECRET_KEY == "" {
		t.Skip("HCSO_ACCESS_KEY and HCSO_SECRET_KEY must be set for OBS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSBucket(t *testing.T) {
	if HCSO_OBS_BUCKET_NAME == "" {
		t.Skip("HCSO_OBS_BUCKET_NAME must be set for OBS object acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSDestinationBucket(t *testing.T) {
	if HCSO_OBS_DESTINATION_BUCKET == "" {
		t.Skip("HCSO_OBS_DESTINATION_BUCKET must be set for OBS destination tests")
	}
}

// lintignore:AT003
func TestAccPreCheckChargingMode(t *testing.T) {
	if HCSO_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

// lintignore:AT003
func TestAccPreCheckHighCostAllow(t *testing.T) {
	if HCSO_HIGH_COST_ALLOW == "" {
		t.Skip("Do not allow expensive testing")
	}
}

// lintignore:AT003
func TestAccPreCheckScm(t *testing.T) {
	if HCSO_CERTIFICATE_KEY_PATH == "" || HCSO_CERTIFICATE_CHAIN_PATH == "" ||
		HCSO_CERTIFICATE_PRIVATE_KEY_PATH == "" || HCSO_CERTIFICATE_SERVICE == "" ||
		HCSO_CERTIFICATE_PROJECT == "" || HCSO_CERTIFICATE_PROJECT_UPDATED == "" {
		t.Skip("HCSO_CERTIFICATE_KEY_PATH, HCSO_CERTIFICATE_CHAIN_PATH, HCSO_CERTIFICATE_PRIVATE_KEY_PATH, " +
			"HCSO_CERTIFICATE_SERVICE, HCSO_CERTIFICATE_PROJECT and HCSO_CERTIFICATE_TARGET_UPDATED " +
			"can not be empty for SCM certificate tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSWRDomian(t *testing.T) {
	if HCSO_SWR_SHARING_ACCOUNT == "" {
		t.Skip("HCSO_SWR_SHARING_ACCOUNT must be set for swr domian tests, " +
			"the value of HCSO_SWR_SHARING_ACCOUNT should be another IAM user name")
	}
}

// lintignore:AT003
func TestAccPreCheckRAM(t *testing.T) {
	if HCSO_RAM_SHARE_ACCOUNT_ID == "" || HCSO_RAM_SHARE_RESOURCE_URN == "" {
		t.Skip("HCSO_RAM_SHARE_ACCOUNT_ID and HCSO_RAM_SHARE_RESOURCE_URN " +
			"must be set for create ram resource tests.")
	}

	if HCSO_RAM_SHARE_UPDATE_ACCOUNT_ID == "" || HCSO_RAM_SHARE_UPDATE_RESOURCE_URN == "" {
		t.Skip("HCSO_RAM_SHARE_UPDATE_ACCOUNT_ID and HCSO_RAM_SHARE_UPDATE_RESOURCE_URN" +
			" must be set for update ram resource tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDms(t *testing.T) {
	if HCSO_DMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DMS tests")
	}
}

// lintignore:AT003
func TestAccPreCheckDliJarPath(t *testing.T) {
	if HCSO_DLI_FLINK_JAR_OBS_PATH == "" {
		t.Skip("HCSO_DLI_FLINK_JAR_OBS_PATH must be set for DLI Flink Jar job acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthCss(t *testing.T) {
	if HCSO_DLI_DS_AUTH_CSS_OBS_PATH == "" {
		t.Skip("HCSO_DLI_DS_AUTH_CSS_OBS_PATH must be set for DLI datasource CSS Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthKafka(t *testing.T) {
	if HCSO_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH == "" || HCSO_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH == "" {
		t.Skip("HCSO_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH,HCSO_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH must be set for DLI datasource Kafka Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthKrb(t *testing.T) {
	if HCSO_DLI_DS_AUTH_KRB_CONF_OBS_PATH == "" || HCSO_DLI_DS_AUTH_KRB_TAB_OBS_PATH == "" {
		t.Skip("HCSO_DLI_DS_AUTH_KRB_CONF_OBS_PATH,HCSO_DLI_DS_AUTH_KRB_TAB_OBS_PATH must be set for DLI datasource Kafka Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliAgency(t *testing.T) {
	if HCSO_DLI_AGENCY_FLAG == "" {
		t.Skip("HCSO_DLI_AGENCY_FLAG must be set for DLI datasource DLI agency acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRepoTokenAuth(t *testing.T) {
	if HCSO_GITHUB_REPO_HOST == "" || HCSO_GITHUB_PERSONAL_TOKEN == "" {
		t.Skip("Repository configurations are not completed for acceptance test of personal access token authorization.")
	}
}

// lintignore:AT003
func TestAccPreCheckRepoPwdAuth(t *testing.T) {
	if HCSO_DOMAIN_NAME == "" || HCSO_USER_NAME == "" || HCSO_GITHUB_REPO_PWD == "" {
		t.Skip("Repository configurations are not completed for acceptance test of password authorization.")
	}
}

// lintignore:AT003
func TestAccPreCheckComponent(t *testing.T) {
	if HCSO_DOMAIN_NAME == "" || HCSO_GITHUB_REPO_URL == "" || HCSO_OBS_STORAGE_URL == "" {
		t.Skip("Repository (package) configurations are not completed for acceptance test of component.")
	}
}

// lintignore:AT003
func TestAccPreCheckComponentDeployment(t *testing.T) {
	if HCSO_BUILD_IMAGE_URL == "" {
		t.Skip("SWR image URL configuration is not completed for acceptance test of component deployment.")
	}
}

// lintignore:AT003
func TestAccPreCheckImageUrlUpdated(t *testing.T) {
	if HCSO_BUILD_IMAGE_URL_UPDATED == "" {
		t.Skip("SWR image update URL configuration is not completed for acceptance test of component deployment.")
	}
}

// lintignore:AT003
func TestAccPreCheckVODWatermark(t *testing.T) {
	if HCSO_VOD_WATERMARK_FILE == "" {
		t.Skip("HCSO_VOD_WATERMARK_FILE must be set for VOD watermark template acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckVODMediaAsset(t *testing.T) {
	if HCSO_VOD_MEDIA_ASSET_FILE == "" {
		t.Skip("HCSO_VOD_MEDIA_ASSET_FILE must be set for VOD media asset acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckPwdAuth(t *testing.T) {
	if HCSO_MEETING_ACCOUNT_NAME == "" || HCSO_MEETING_ACCOUNT_PASSWORD == "" {
		t.Skip("The account name (HCSO_MEETING_ACCOUNT_NAME) or password (HCSO_MEETING_ACCOUNT_PASSWORD) is not " +
			"completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckAppAuth(t *testing.T) {
	if HCSO_MEETING_APP_ID == "" || HCSO_MEETING_APP_KEY == "" || HCSO_MEETING_USER_ID == "" {
		t.Skip("The app ID (HCSO_MEETING_APP_ID), app KEY (HCSO_MEETING_APP_KEY) or user ID (HCSO_MEETING_USER_ID) is not " +
			"completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckMeetingRoom(t *testing.T) {
	if HCSO_MEETING_ROOM_ID == "" {
		t.Skip("The vmr ID (HCSO_MEETING_ROOM_ID) is not completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckParticipants(t *testing.T) {
	if HCSO_CHAIR_EMAIL == "" || HCSO_GUEST_EMAIL == "" {
		t.Skip("The chair (HCSO_CHAIR_EMAIL) or guest (HCSO_GUEST_EMAIL) mailbox is not completed for acceptance test of " +
			"conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckAadForwardRule(t *testing.T) {
	if HCSO_AAD_INSTANCE_ID == "" || HCSO_AAD_IP_ADDRESS == "" {
		t.Skip("The instance information is not completed for AAD rule acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckScmCertificateName(t *testing.T) {
	if HCSO_CERTIFICATE_NAME == "" {
		t.Skip("HCSO_CERTIFICATE_NAME must be set for SCM acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckKms(t *testing.T) {
	if HCSO_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

// lintignore:AT003
func TestAccPreCheckProjectID(t *testing.T) {
	if HCSO_PROJECT_ID == "" {
		t.Skip("HCSO_PROJECT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkspaceAD(t *testing.T) {
	if HCSO_WORKSPACE_AD_DOMAIN_NAME == "" || HCSO_WORKSPACE_AD_SERVER_PWD == "" || HCSO_WORKSPACE_AD_DOMAIN_IP == "" ||
		HCSO_WORKSPACE_AD_VPC_ID == "" || HCSO_WORKSPACE_AD_NETWORK_ID == "" {
		t.Skip("The configuration of AD server is not completed for Workspace service acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckER(t *testing.T) {
	if HCSO_ER_TEST_ON == "" {
		t.Skip("Skip all ER acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRfArchives(t *testing.T) {
	if HCSO_RF_TEMPLATE_ARCHIVE_NO_VARS_URI == "" || HCSO_RF_TEMPLATE_ARCHIVE_URI == "" ||
		HCSO_RF_VARIABLES_ARCHIVE_URI == "" {
		t.Skip("Skip the archive URI parameters acceptance test for RF resource stack.")
	}
}

// lintignore:AT003
func TestAccPreCheckDcDirectConnection(t *testing.T) {
	if HCSO_DC_DIRECT_CONNECT_ID == "" {
		t.Skip("Skip the interface acceptance test because of the direct connection ID is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckCfw(t *testing.T) {
	if HCSO_CFW_INSTANCE_ID == "" {
		t.Skip("HCSO_CFW_INSTANCE_ID must be set for CFW acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadType(t *testing.T) {
	if HCSO_WORKLOAD_TYPE == "" {
		t.Skip("HCSO_WORKLOAD_TYPE must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadName(t *testing.T) {
	if HCSO_WORKLOAD_NAME == "" {
		t.Skip("HCSO_WORKLOAD_NAME must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCceClusterId(t *testing.T) {
	if HCSO_CCE_CLUSTER_ID == "" {
		t.Skip("HCSO_CCE_CLUSTER_ID must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadNameSpace(t *testing.T) {
	if HCSO_WORKLOAD_NAMESPACE == "" {
		t.Skip("HCSO_WORKLOAD_NAMESPACE must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrTargetRegion(t *testing.T) {
	if HCSO_SWR_TARGET_REGION == "" {
		t.Skip("HCSO_SWR_TARGET_REGION must be set for SWR image auto sync tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrTargetOrigination(t *testing.T) {
	if HCSO_SWR_TARGET_ORGANIZATION == "" {
		t.Skip("HCSO_SWR_TARGET_ORGANIZATION must be set for SWR image auto sync tests")
	}
}

// lintignore:AT003
func TestAccPreCheckImsBackupId(t *testing.T) {
	if HCSO_IMS_BACKUP_ID == "" {
		t.Skip("HCSO_IMS_BACKUP_ID must be set for IMS whole image with CBR backup id")
	}
}

// lintignore:AT003
func TestAccPreCheckSourceImage(t *testing.T) {
	if HCSO_IMAGE_SHARE_SOURCE_IMAGE_ID == "" {
		t.Skip("Skip the interface acceptance test because of the source image ID is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckSecMaster(t *testing.T) {
	if HCSO_SECMASTER_WORKSPACE_ID == "" {
		t.Skip("HCSO_SECMASTER_WORKSPACE_ID must be set for SecMaster acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCcePartitionAz(t *testing.T) {
	if HCSO_CCE_PARTITION_AZ == "" {
		t.Skip("Skip the interface acceptance test because of the cce partition az is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckCnEast3(t *testing.T) {
	if HCSO_REGION_NAME != "cn-east-3" {
		t.Skip("HCSO_REGION_NAME must be cn-east-3 for this test.")
	}
}

// lintignore:AT003
func TestAccPreCheckCertificateWithoutRootCA(t *testing.T) {
	if HCSO_CERTIFICATE_CONTENT == "" || HCSO_CERTIFICATE_PRIVATE_KEY == "" ||
		HCSO_NEW_CERTIFICATE_CONTENT == "" || HCSO_NEW_CERTIFICATE_PRIVATE_KEY == "" {
		t.Skip("HCSO_CERTIFICATE_CONTENT, HCSO_CERTIFICATE_PRIVATE_KEY, HCSO_NEW_CERTIFICATE_CONTENT and " +
			"HCSO_NEW_CERTIFICATE_PRIVATE_KEY must be set for simple acceptance tests of SSL certificate resource")
	}
}

// lintignore:AT003
func TestAccPreCheckCertificateFull(t *testing.T) {
	TestAccPreCheckCertificateWithoutRootCA(t)
	if HCSO_CERTIFICATE_ROOT_CA == "" || HCSO_NEW_CERTIFICATE_ROOT_CA == "" {
		t.Skip("HCSO_CERTIFICATE_ROOT_CA and HCSO_NEW_CERTIFICATE_ROOT_CA must be set for root CA validation")
	}
}

// lintignore:AT003
func TestAccPreCheckCodeArtsDeployResourcePoolID(t *testing.T) {
	if HCSO_CODEARTS_RESOURCE_POOL_ID == "" {
		t.Skip("HCSO_CODEARTS_RESOURCE_POOL_ID must be set for this acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckCodeArtsEnableFlag(t *testing.T) {
	if HCSO_CODEARTS_ENABLE_FLAG == "" {
		t.Skip("Skip the CodeArts acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckModelArtsHasSubscribeModel(t *testing.T) {
	if HCSO_MODELARTS_HAS_SUBSCRIBE_MODEL == "" {
		t.Skip("Subscribe two free models from market and set HCSO_MODELARTS_HAS_SUBSCRIBE_MODEL" +
			" for modelarts service acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckEgChannelId(t *testing.T) {
	if HCSO_EG_CHANNEL_ID == "" {
		t.Skip("The sub-resource acceptance test of the EG channel must set 'HCSO_EG_CHANNEL_ID'")
	}
}
