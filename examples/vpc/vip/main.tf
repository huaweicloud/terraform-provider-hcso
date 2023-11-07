data "hcso_availability_zones" "myaz" {}

data "hcso_compute_flavors" "myflavor" {
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcso_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "hcso_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcso_vpc_subnet" "subnet_1" {
  vpc_id      = hcso_vpc.vpc_1.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "hcso_compute_instance" "mycompute" {
  name              = "mycompute_${count.index}"
  image_id          = data.hcso_images_image.myimage.id
  flavor_id         = data.hcso_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.hcso_availability_zones.myaz.names[0]

  network {
    uuid = hcso_vpc_subnet.subnet_1.id
  }
  count = 2
}

resource "hcso_networking_vip" "vip_1" {
  network_id = hcso_vpc_subnet.subnet_1.id
}

# associate ports to the vip
resource "hcso_networking_vip_associate" "vip_associated" {
  vip_id   = hcso_networking_vip.vip_1.id
  port_ids = [
    hcso_compute_instance.mycompute[0].network[0].port,
    hcso_compute_instance.mycompute[1].network[0].port
  ]
}
