output "resource_pool" {
  value = data.vsphere_compute_cluster.cluster.resource_pool_id
}

output "datastore" {
  value = data.vsphere_datastore.datastore.id
}

output "folder" {
  value = local.folder
}

output "network" {
  value = data.vsphere_network.network.id
}

output "datacenter" {
  value = data.vsphere_datacenter.datacenter.id
}

output "template" {
  value = data.vsphere_virtual_machine.template.id
}

output "guest_id" {
  value = data.vsphere_virtual_machine.template.guest_id
}

output "thin_disk" {
  value = data.vsphere_virtual_machine.template.disks.0.thin_provisioned
}

output "scrub_disk" {
  value = data.vsphere_virtual_machine.template.disks.0.eagerly_scrub
}

output "cluster_domain" {
  value = var.cluster_domain
}

output "cluster_id" {
  value = var.cluster_id
}

output "tags" {
  value = [vsphere_tag.tag.id]
}

output "master_memory" {
  value = var.vsphere_control_plane_memory_mib
}

output "master_num_cpus" {
  value = var.vsphere_control_plane_num_cpus
}

output "master_cores_per_socket" {
  value = var.vsphere_control_plane_cores_per_socket
}

output "master_disk_size" {
  value = var.vsphere_control_plane_disk_gib
}


