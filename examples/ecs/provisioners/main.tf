data "hcso_availability_zones" "default" {}

data "hcso_images_image" "default" {
  name        = var.image_name
  most_recent = true
}

data "hcso_compute_flavors" "default" {
  availability_zone = data.hcso_availability_zones.default.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcso_compute_keypair" "default" {
  name     = var.keypair_name
  key_file = var.private_key_path
}

resource "hcso_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcso_vpc_subnet" "default" {
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  vpc_id     = hcso_vpc.default.id
  gateway_ip = var.gateway_ip
}

resource "hcso_networking_secgroup" "default" {
  name = var.security_group_name
}

resource "hcso_vpc_eip" "default" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcso_compute_instance" "default" {
  name              = var.ecs_instance_name
  image_id          = data.hcso_images_image.default.id
  flavor_id         = data.hcso_compute_flavors.default.ids[0]
  availability_zone = data.hcso_availability_zones.default.names[0]
  key_pair          = hcso_compute_keypair.default.name
  user_data         = <<-EOF
#!/bin/bash
echo '${file("./test.txt")}' > /home/test.txt
EOF

  system_disk_type   = "SSD"
  system_disk_size   = 40

  security_groups = [
    hcso_networking_secgroup.default.name
  ]

  network {
    uuid = hcso_vpc_subnet.default.id
  }
}

resource "hcso_compute_eip_associate" "default" {
  public_ip   = hcso_vpc_eip.default.address
  instance_id = hcso_compute_instance.default.id
}

resource "null_resource" "provision" {
  depends_on = [hcso_compute_eip_associate.default]

  provisioner "remote-exec" {
    connection {
      user        = "root"
      private_key = file(var.private_key_path)
      host        = hcso_vpc_eip.default.address
    }

    inline = [
      "cat /home/test.txt"
    ]
  }
}
