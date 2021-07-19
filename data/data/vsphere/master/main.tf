locals {
  // limitation of baremetal-runtimecfg.  The hostname must be master
  master_name = "master"
  description = "Created By OpenShift Installer"
}

provider "vsphere" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}

provider "vsphereprivate" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}

resource "vsphere_virtual_machine" "vm_master" {
  count = var.master_count

  name                 = "${var.cluster_id}-${local.master_name}-${count.index}"
  resource_pool_id     = var.resource_pool
  datastore_id         = var.datastore
  num_cpus             = var.master_num_cpus
  num_cores_per_socket = var.master_cores_per_socket
  memory               = var.master_memory
  guest_id             = var.guest_id
  folder               = var.folder
  enable_disk_uuid     = "true"
  annotation           = local.description

  wait_for_guest_net_timeout = 0
  wait_for_guest_ip_timeout  = 15

  network_interface {
    network_id = var.network
  }

  disk {
    label            = "disk0"
    size             = var.master_disk_size
    eagerly_scrub    = var.scrub_disk
    thin_provisioned = var.thin_disk
  }

  clone {
    template_uuid = var.template
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_master)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-${local.master_name}-${count.index}"
  }

  tags = var.tags
}

