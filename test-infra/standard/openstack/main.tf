variable "key_name" {}
variable "network_id" {}

variable "pool" {
  default = "public"
}

variable "flavor" {
  default = "m1.xlarge"
}

data "openstack_images_image_v2" "packstack_standard" {
  name = "packstack-standard-ocata"
  most_recent = true
}

resource "openstack_networking_floatingip_v2" "openstack_acc_tests" {
  pool = "${var.pool}"
}

resource "openstack_networking_secgroup_v2" "openstack_acc_tests" {
  name = "openstack_acc_tests"
  description = "Rules for openstack acceptance tests"
}

resource "openstack_networking_secgroup_rule_v2" "openstack_acc_tests_rule_1" {
  security_group_id = "${openstack_networking_secgroup_v2.openstack_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "tcp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "0.0.0.0/0"
}

resource "openstack_networking_secgroup_rule_v2" "openstack_acc_tests_rule_2" {
  security_group_id = "${openstack_networking_secgroup_v2.openstack_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "tcp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "::/0"
}

resource "openstack_networking_secgroup_rule_v2" "openstack_acc_tests_rule_3" {
  security_group_id = "${openstack_networking_secgroup_v2.openstack_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "udp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "0.0.0.0/0"
}

resource "openstack_networking_secgroup_rule_v2" "openstack_acc_tests_rule_4" {
  security_group_id = "${openstack_networking_secgroup_v2.openstack_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "udp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "::/0"
}

resource "openstack_networking_secgroup_rule_v2" "openstack_acc_tests_rule_5" {
  security_group_id = "${openstack_networking_secgroup_v2.openstack_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "icmp"
  remote_ip_prefix = "0.0.0.0/0"
}

resource "openstack_networking_secgroup_rule_v2" "openstack_acc_tests_rule_6" {
  security_group_id = "${openstack_networking_secgroup_v2.openstack_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "icmp"
  remote_ip_prefix = "::/0"
}

resource "openstack_compute_instance_v2" "openstack_acc_tests" {
  name = "openstack_acc_tests"
  image_id = "${data.openstack_images_image_v2.packstack_standard.id}"
  flavor_name = "${var.flavor}"
  key_pair = "${var.key_name}"

  security_groups = ["${openstack_networking_secgroup_v2.openstack_acc_tests.name}"]

  network {
    uuid = "${var.network_id}"
  }
}

resource "openstack_compute_floatingip_associate_v2" "openstack_acc_tests" {
  instance_id = "${openstack_compute_instance_v2.openstack_acc_tests.id}"
  floating_ip = "${openstack_networking_floatingip_v2.openstack_acc_tests.address}"
}

resource "null_resource" "rc_files" {
  provisioner "local-exec" {
    command = <<EOF
      while true ; do
        wget http://${openstack_compute_floatingip_associate_v2.openstack_acc_tests.floating_ip}/keystonerc_demo 2> /dev/null
        if [ $? = 0 ]; then
          break
        fi
        sleep 20
      done

      wget http://${openstack_compute_floatingip_associate_v2.openstack_acc_tests.floating_ip}/keystonerc_admin
    EOF
  }
}
