package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOTCVBSBackupV2_importBasic(t *testing.T) {
	resourceName := "flexibleengine_vbs_backup_v2.backup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
