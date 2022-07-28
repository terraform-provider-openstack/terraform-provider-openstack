resource "openstack_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = var.volume_size
}

resource "openstack_compute_volume_attach_v2" "attached" {
  instance_id = openstack_compute_instance_v2.my_instance.id
  volume_id   = openstack_blockstorage_volume_v2.volume_1.id
  # Prevent re-creation
  #   lifecycle {
  #     ignore_changes = ["volume_id", "instance_id"]
  #   }
}
