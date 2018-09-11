package flexibleengine

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/servers"
	"fmt"
)

func TestAccComputeV2BmsInstance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2BmsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2BmsInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func TestAccComputeV2BmsInstance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2BmsInstance_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func testAccCheckComputeV2BmsInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2HWClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_compute_instance_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2BmsInstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2HWClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Bms Instance not found")
		}

		*instance = *found

		return nil
	}
}


var testAccComputeV2BmsInstance_basic = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2BmsInstance_timeout = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }

  timeouts {
    create = "10m"
  }
}
`, OS_NETWORK_ID)