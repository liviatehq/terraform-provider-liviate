data "liviate_template" "os" {
  template_filter = "executable"
  filter {
    name  = "name"
    value = "debian-13-std-30"
  }
}

data "liviate_service_offering" "small" {
  filter {
    name  = "name"
    value = "liviate-vcpu-2-4gb-20260317-01"
  }
}

data "liviate_zone" "primary" {
  filter {
    name  = "name"
    value = ".*"
  }
}

resource "liviate_ssh_keypair" "ssh" {
  name       = var.ssh_key_name
  public_key = file(var.public_key_path)
}

resource "liviate_instance" "web-01" {
  name             = "liviate-web-01"
  display_name     = "Web Server 01"
  service_offering = data.liviate_service_offering.small.name
  template         = data.liviate_template.os.id
  zone             = var.zone_id == "" ? data.liviate_zone.primary.id : var.zone_id
  root_disk_size   = 30
  keypair          = liviate_ssh_keypair.ssh.name
}

resource "liviate_instance" "web-02" {
  name             = "liviate-web-02"
  display_name     = "Web Server 02"
  service_offering = data.liviate_service_offering.small.name
  template         = data.liviate_template.os.id
  zone             = var.zone_id == "" ? data.liviate_zone.primary.id : var.zone_id
  root_disk_size   = 30
  keypair          = liviate_ssh_keypair.ssh.name
}
