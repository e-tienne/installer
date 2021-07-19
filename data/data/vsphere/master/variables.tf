variable "resource_pool" {
  type = string
}

variable "folder" {
  type = string
}

variable "datastore" {
  type = string
}

variable "network" {
  type = string
}

variable "datacenter" {
  type = string
}

variable "template" {
  type = string
}

variable "guest_id" {
  type = string
}

variable "master_memory" {
  type = number
}

variable "master_num_cpus" {
  type = number
}

variable "master_cores_per_socket" {
  type = number
}

variable "master_disk_size" {
  type = number
}

variable "tags" {
  type = list
}

variable "thin_disk" {
  type = bool
}

variable "scrub_disk" {
  type = bool
}
