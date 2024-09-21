resource "openstack_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = file("${var.ssh_key_file}.pub")
}

resource "openstack_compute_instance_v2" "my_instance" {
  name            = "my_instance"
  image_name      = var.image
  flavor_name     = var.flavor
  key_pair        = openstack_compute_keypair_v2.terraform.name
  security_groups = ["default"]
  network {
    name = var.network_name
  }
}

resource "openstack_networking_floatingip_v2" "fip" {
  pool = var.pool
}

data "openstack_networking_port_v2" "port" {
  device_id  = openstack_compute_instance_v2.my_instance.id
  network_id = openstack_compute_instance_v2.my_instance.network.0.uuid
}

resource "openstack_networking_floatingip_associate_v2" "fip_associate" {
  floating_ip = openstack_networking_floatingip_v2.fip.address
  port_id     = data.openstack_networking_port_v2.port.id
  connection {
    host        = openstack_networking_floatingip_v2.fip.address
    user        = var.ssh_user_name
    private_key = file(var.ssh_key_file)
  }

  provisioner "local-exec" {
    command = "echo ${openstack_networking_floatingip_v2.fip.address} > instance_ip.txt"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mkfs.ext4 ${openstack_compute_volume_attach_v2.attached.device}",
      "sudo mkdir /mnt/volume",
      "sudo mount ${openstack_compute_volume_attach_v2.attached.device} /mnt/volume",
      "sudo df -h /mnt/volume",
    ]
  }
}
