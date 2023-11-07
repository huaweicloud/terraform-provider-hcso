package internal

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/helper/pathorcontents"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

//nolint:revive
var (
	HCSO_AVAILABILITY_ZONE          = os.Getenv("HCSO_AVAILABILITY_ZONE")
	HCSO_DEPRECATED_ENVIRONMENT     = os.Getenv("HCSO_DEPRECATED_ENVIRONMENT")
	HCSO_EXTGW_ID                   = os.Getenv("HCSO_EXTGW_ID")
	HCSO_FLAVOR_ID                  = os.Getenv("HCSO_FLAVOR_ID")
	HCSO_FLAVOR_NAME                = os.Getenv("HCSO_FLAVOR_NAME")
	HCSO_IMAGE_ID                   = os.Getenv("HCSO_IMAGE_ID")
	HCSO_IMAGE_NAME                 = os.Getenv("HCSO_IMAGE_NAME")
	HCSO_NETWORK_ID                 = os.Getenv("HCSO_NETWORK_ID")
	HCSO_SUBNET_ID                  = os.Getenv("HCSO_SUBNET_ID")
	HCSO_POOL_NAME                  = os.Getenv("HCSO_POOL_NAME")
	HCSO_REGION_NAME                = os.Getenv("HCSO_REGION_NAME")
	HCSO_ACCESS_KEY                 = os.Getenv("HCSO_ACCESS_KEY")
	HCSO_SECRET_KEY                 = os.Getenv("HCSO_SECRET_KEY")
	HCSO_VPC_ID                     = os.Getenv("HCSO_VPC_ID")
	HCSO_CCI_NAMESPACE              = os.Getenv("HCSO_CCI_NAMESPACE")
	HCSO_PROJECT_ID                 = os.Getenv("HCSO_PROJECT_ID")
	HCSO_DOMAIN_ID                  = os.Getenv("HCSO_DOMAIN_ID")
	HCSO_DOMAIN_NAME                = os.Getenv("HCSO_DOMAIN_NAME")
	HCSO_MRS_ENVIRONMENT            = os.Getenv("HCSO_MRS_ENVIRONMENT")
	HCSO_KMS_ENVIRONMENT            = os.Getenv("HCSO_KMS_ENVIRONMENT")
	HCSO_CCI_ENVIRONMENT            = os.Getenv("HCSO_CCI_ENVIRONMENT")
	HCSO_CDN_DOMAIN_NAME            = os.Getenv("HCSO_CDN_DOMAIN_NAME")
	HCSO_CDN_CERT_PATH              = os.Getenv("HCSO_CDN_CERT_PATH")
	HCSO_CDN_PRIVATE_KEY_PATH       = os.Getenv("HCSO_CDN_PRIVATE_KEY_PATH")
	HCSO_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("HCSO_ENTERPRISE_PROJECT_ID_TEST")
	HCSO_USER_ID                    = os.Getenv("HCSO_USER_ID")
	HCSO_CHARGING_MODE              = os.Getenv("HCSO_CHARGING_MODE")
)

var testAccProviders map[string]*schema.Provider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"internal": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HCSO_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}
}

func testAccPreCheckDeprecated(t *testing.T) {
	if HCSO_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}
}

func testAccPreCheckMrs(t *testing.T) {
	if HCSO_MRS_ENVIRONMENT == "" {
		t.Skip("This environment does not support MRS tests")
	}
}

func testAccPreCheckKms(t *testing.T) {
	if HCSO_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

func testAccPreCheckCDN(t *testing.T) {
	if HCSO_CDN_DOMAIN_NAME == "" {
		t.Skip("This environment does not support CDN tests")
	}
}

func testAccPreCheckCERT(t *testing.T) {
	if HCSO_CDN_CERT_PATH == "" || HCSO_CDN_PRIVATE_KEY_PATH == "" {
		t.Skip("This environment does not support CDN certificate tests")
	}
}

func testAccPreCheckCCINamespace(t *testing.T) {
	if HCSO_CCI_NAMESPACE == "" {
		t.Skip("This environment does not support CCI Namespace tests")
	}
}

func testAccPreCheckCCI(t *testing.T) {
	if HCSO_CCI_ENVIRONMENT == "" {
		t.Skip("This environment does not support CCI tests")
	}
}

func testAccPreCheckEpsID(t *testing.T) {
	if HCSO_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("This environment does not support Enterprise Project ID tests")
	}
}

func testAccPreCheckChargingMode(t *testing.T) {
	if HCSO_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Steps for configuring HuaweiCloud with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping HuaweiCloud CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloud CA by file: %s", diags[0].Summary)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping HuaweiCloud CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloud CA by string: %s", diags[0].Summary)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping HuaweiCloud client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(certFile)
	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(keyFile)

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloud Client keypair by file: %s", diags[0].Summary)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping HuaweiCloud client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloud Client keypair by contents: %s", diags[0].Summary)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmtp.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", varName)
	if err != nil {
		return "", fmtp.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmtp.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmtp.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}
