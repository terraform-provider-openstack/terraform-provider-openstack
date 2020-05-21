variable "image" {
  default = "Ubuntu 18.04"
}

variable "flavor" {
  default = "m1.small"
}

variable "ssh_key_file" {
  default = "~/.ssh/id_rsa"
}
variable "volume_size" {
  default = 1
}

variable "ssh_user_name" {
  default = "ubuntu"
}

variable "pool" {
  default = "public"
}

variable "network_name" {
  default = "my_network"
}

