package deprecated

import (
	"fmt"
	"testing"
	"time"

	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// FirewallGroup is an HuaweiCloud firewall group.
type FirewallGroup struct {
	firewall_groups.FirewallGroup
	routerinsertion.FirewallGroupExt
}

func TestAccFWFirewallGroupV2_basic(t *testing.T) {
	var epolicyID *string
	var ipolicyID *string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallGroupV2_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2("hcso_fw_firewall_group_v2.fw_1", "", "", ipolicyID, epolicyID),
				),
			},
			{
				ResourceName:      "hcso_fw_firewall_group_v2.fw_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFWFirewallGroupV2_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2(
						"hcso_fw_firewall_group_v2.fw_1", "fw_1", "terraform acceptance test", ipolicyID, epolicyID),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port0(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("hcso_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_no_ports(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_no_ports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("hcso_fw_firewall_group_v2.fw_1", &firewall_group),
					resource.TestCheckResourceAttr("hcso_fw_firewall_group_v2.fw_1", "description", "firewall router test"),
					testAccCheckFWFirewallPortCount(&firewall_group, 0),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_update(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("hcso_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
			{
				Config: testAccFWFirewallV2_port_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("hcso_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 2),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_remove(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("hcso_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
			{
				Config: testAccFWFirewallV2_port_remove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("hcso_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 0),
				),
			},
		},
	})
}

func testAccCheckFWFirewallGroupV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	fwClient, err := config.FwV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_firewall_group" {
			continue
		}

		_, err = firewall_groups.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Firewall group (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWFirewallGroupV2Exists(n string, firewall_group *FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Exists) Error creating HuaweiCloud fw client: %s", err)
		}

		var found FirewallGroup
		err = firewall_groups.Get(fwClient, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Firewall group not found")
		}

		*firewall_group = found

		return nil
	}
}

func testAccCheckFWFirewallPortCount(firewall_group *FirewallGroup, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall_group.PortIDs) != expected {
			return fmtp.Errorf("Expected %d Ports, got %d", expected, len(firewall_group.PortIDs))
		}

		return nil
	}
}

func testAccCheckFWFirewallGroupV2(n, expectedName, expectedDescription string, ipolicyID *string, epolicyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Exists) Error creating HuaweiCloud fw client: %s", err)
		}

		var found *firewall_groups.FirewallGroup
		for i := 0; i < 5; i++ {
			// Firewall creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = firewall_groups.Get(fwClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					//lintignore:R018
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case found.Name != expectedName:
			err = fmtp.Errorf("Expected Name to be <%s> but found <%s>", expectedName, found.Name)
		case found.Description != expectedDescription:
			err = fmtp.Errorf("Expected Description to be <%s> but found <%s>",
				expectedDescription, found.Description)
		case found.IngressPolicyID == "":
			err = fmtp.Errorf("Ingress Policy should not be empty")
		case found.EgressPolicyID == "":
			err = fmtp.Errorf("Egress Policy should not be empty")
		case ipolicyID != nil && found.IngressPolicyID == *ipolicyID:
			err = fmtp.Errorf("Ingress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		case epolicyID != nil && found.EgressPolicyID == *epolicyID:
			err = fmtp.Errorf("Egress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		}

		if err != nil {
			return err
		}

		ipolicyID = &found.IngressPolicyID
		epolicyID = &found.EgressPolicyID

		return nil
	}
}

const testAccFWFirewallGroupV2_basic_1 = `
resource "hcso_fw_firewall_group_v2" "fw_1" {
  ingress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "hcso_fw_policy_v2" "policy_1" {
  name = "policy_1"
}
`

const testAccFWFirewallGroupV2_basic_2 = `
resource "hcso_fw_firewall_group_v2" "fw_1" {
  name = "fw_1"
  description = "terraform acceptance test"
  ingress_policy_id = "${hcso_fw_policy_v2.policy_2.id}"
  egress_policy_id = "${hcso_fw_policy_v2.policy_2.id}"
  admin_state_up = true

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "hcso_fw_policy_v2" "policy_2" {
  name = "policy_2"
}
`

var testAccFWFirewallV2_port = fmt.Sprintf(`
resource "hcso_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hcso_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  enable_dhcp = true
  network_id = "${hcso_networking_network_v2.network_1.id}"
}

resource "hcso_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"
}

resource "hcso_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hcso_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hcso_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.23"
  }
}

resource "hcso_networking_router_interface_v2" "router_interface_1" {
  router_id = "${hcso_networking_router_v2.router_1.id}"
  port_id = "${hcso_networking_port_v2.port_1.id}"
}

resource "hcso_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "hcso_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  #egress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  ports = [
	"${hcso_networking_port_v2.port_1.id}"
  ]
  depends_on = ["hcso_networking_router_interface_v2.router_interface_1"]
}
`, acceptance.HCSO_EXTGW_ID)

var testAccFWFirewallV2_port_add = fmt.Sprintf(`
resource "hcso_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hcso_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hcso_networking_network_v2.network_1.id}"
}

resource "hcso_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%[1]s"
}

resource "hcso_networking_router_v2" "router_2" {
  name = "router_2"
  admin_state_up = "true"
  external_network_id = "%[1]s"
}

resource "hcso_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hcso_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hcso_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.23"
  }
}

resource "hcso_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${hcso_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hcso_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.24"
  }
}

resource "hcso_networking_router_interface_v2" "router_interface_1" {
  router_id = "${hcso_networking_router_v2.router_1.id}"
  port_id = "${hcso_networking_port_v2.port_1.id}"
}

resource "hcso_networking_router_interface_v2" "router_interface_2" {
  router_id = "${hcso_networking_router_v2.router_2.id}"
  port_id = "${hcso_networking_port_v2.port_2.id}"
}

resource "hcso_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "hcso_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  ports = [
	"${hcso_networking_port_v2.port_1.id}",
	"${hcso_networking_port_v2.port_2.id}"
  ]
  depends_on = ["hcso_networking_router_interface_v2.router_interface_1", "hcso_networking_router_interface_v2.router_interface_2"]
}
`, acceptance.HCSO_EXTGW_ID)

const testAccFWFirewallV2_port_remove = `
resource "hcso_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "hcso_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
}
`

const testAccFWFirewallV2_no_ports = `
resource "hcso_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "hcso_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${hcso_fw_policy_v2.policy_1.id}"
}
`
