data "hcso_availability_zones" "myaz" {}

data "hcso_compute_flavors" "myflavor" {
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcso_vpc_subnets" "mynet" {
  name = "subnet-default"
}

data "hcso_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
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

resource "hcso_evs_volume" "myvolume" {
  name              = "myvolume"
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  volume_type       = "SSD"
  size              = 10
}

resource "hcso_compute_volume_attach" "attached" {
  instance_id = hcso_compute_instance.myinstance.id
  volume_id   = hcso_evs_volume.myvolume.id
}
