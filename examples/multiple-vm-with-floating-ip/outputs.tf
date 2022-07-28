output "floating_ip_addresses" {
  value = openstack_networking_floatingip_v2.fip.*.address
}
