data "liviate_template" "os" {
  template_filter = "executable"
  filter {
    name  = "name"
    value = "^Ubuntu.*"
  }
}

data "liviate_service_offering" "small" {
  filter {
    name  = "name"
    value = "Small Instance"
  }
}

data "liviate_zone" "primary" {
  filter {
    name  = "name"
    value = ".*"
  }
}

resource "liviate_instance" "web" {
  name             = "liviate-web-01"
  display_name     = "Web Server 01"
  service_offering = data.liviate_service_offering.small.name
  template         = data.liviate_template.os.id
  zone             = var.zone_id == "" ? data.liviate_zone.primary.id : var.zone_id
  root_disk_size   = 20
}
