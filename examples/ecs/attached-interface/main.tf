data "hcso_availability_zones" "myaz" {}

data "hcso_compute_flavors" "myflavor" {
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcso_vpc_subnets" "mynet" {
  name = "subnet-default"
}

data "hcso_images_image" "myimage" {
  name_regex  = "^Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "hcso_compute_instance" "myinstance" {
  name              = "basic"
  image_id          = data.hcso_images_image.myimage.id
  flavor_id         = data.hcso_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.hcso_availability_zones.myaz.names[0]

  system_disk_type   = "SSD"
  system_disk_size   = 40

  network {
    uuid = data.hcso_vpc_subnets.mynet.subnets[0].id
  }
}

data "hcso_vpc" "myvpc" {
  name = "vpc-default"
}

resource "hcso_vpc_subnet" "attach" {
  name       = "subnet-attach"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = data.hcso_vpc.myvpc.id

  availability_zone = data.hcso_availability_zones.myaz.names[0]
}

resource "hcso_compute_interface_attach" "attached" {
  instance_id = hcso_compute_instance.myinstance.id
  network_id  = hcso_vpc_subnet.attach.id

  # This is optional
  fixed_ip = "192.168.1.100"
}
