---
layout: "openstack"
page_title: "OpenStack: openstack_compute_instance_v2"
sidebar_current: "docs-openstack-resource-compute-instance-v2"
description: |-
  Manages a V2 VM instance resource within OpenStack.
---

# openstack\_compute\_instance\_v2

Manages a V2 VM instance resource within OpenStack.

~> **Note:** All arguments including the instance admin password will be stored
in the raw state as plain-text. [Read more about sensitive data in
state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### Basic Instance

```hcl
resource "openstack_compute_instance_v2" "basic" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  metadata = {
    this = "that"
  }

  network {
    name = "my_network"
  }
}
```

### Instance With Attached Volume

```hcl
resource "openstack_blockstorage_volume_v2" "myvol" {
  name = "myvol"
  size = 1
}

resource "openstack_compute_instance_v2" "myinstance" {
  name            = "myinstance"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }
}

resource "openstack_compute_volume_attach_v2" "attached" {
  instance_id = "${openstack_compute_instance_v2.myinstance.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.myvol.id}"
}
```

### Boot From Volume

```hcl
resource "openstack_compute_instance_v2" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    volume_size           = 5
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = "my_network"
  }
}
```

### Boot From an Existing Volume

```hcl
resource "openstack_blockstorage_volume_v1" "myvol" {
  name     = "myvol"
  size     = 5
  image_id = "<image-id>"
}

resource "openstack_compute_instance_v2" "boot-from-volume" {
  name            = "bootfromvolume"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "${openstack_blockstorage_volume_v1.myvol.id}"
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = "my_network"
  }
}
```

### Boot Instance, Create Volume, and Attach Volume as a Block Device

```hcl
resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_size           = 1
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Boot Instance and Attach Existing Volume as a Block Device

```hcl
resource "openstack_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    uuid                  = "${openstack_blockstorage_volume_v2.volume_1.id}"
    source_type           = "volume"
    destination_type      = "volume"
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Instance With Multiple Networks

```hcl
resource "openstack_networking_floatingip_v2" "myip" {
  pool = "my_pool"
}

resource "openstack_compute_instance_v2" "multi-net" {
  name            = "multi-net"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_first_network"
  }

  network {
    name = "my_second_network"
  }
}

resource "openstack_compute_floatingip_associate_v2" "myip" {
  floating_ip = "${openstack_networking_floatingip_v2.myip.address}"
  instance_id = "${openstack_compute_instance_v2.multi-net.id}"
  fixed_ip    = "${openstack_compute_instance_v2.multi-net.network.1.fixed_ip_v4}"
}
```

### Instance With Personality

```hcl
resource "openstack_compute_instance_v2" "personality" {
  name            = "personality"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  personality {
    file    = "/path/to/file/on/instance.txt"
    content = "contents of file"
  }

  network {
    name = "my_network"
  }
}
```

### Instance with Multiple Ephemeral Disks

```hcl
resource "openstack_compute_instance_v2" "multi-eph" {
  name            = "multi_eph"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
    guest_format          = "ext4"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instance with Boot Disk and Swap Disk

```hcl
resource "openstack_compute_flavor_v2" "flavor-with-swap" {
  name  = "flavor-with-swap"
  ram   = "8096"
  vcpus = "2"
  disk  = "20"
  swap  = "4096"
}

resource "openstack_compute_instance_v2" "vm-swap" {
  name            = "vm_swap"
  flavor_id       = "${openstack_compute_flavor_v2.flavor-with-swap.id}"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    guest_format          = "swap"
    volume_size           = 4
  }
}
```

### Instance with User Data (cloud-init)

```hcl
resource "openstack_compute_instance_v2" "instance_1" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  user_data       = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"

  network {
    name = "my_network"
  }
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the server instance. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new server.

* `name` - (Required) A unique name for the resource.

* `image_id` - (Optional; Required if `image_name` is empty and not booting
    from a volume. Do not specify if booting from a volume.) The image ID of
    the desired image for the server. Changing this rebuilds the existing
    server.

* `image_name` - (Optional; Required if `image_id` is empty and not booting
    from a volume. Do not specify if booting from a volume.) The name of the
    desired image for the server. Changing this rebuilds the existing server.

* `flavor_id` - (Optional; Required if `flavor_name` is empty) The flavor ID of
    the desired flavor for the server. Changing this resizes the existing server.

* `flavor_name` - (Optional; Required if `flavor_id` is empty) The name of the
    desired flavor for the server. Changing this resizes the existing server.

* `user_data` - (Optional) The user data to provide when launching the instance.
    Changing this creates a new server.

* `security_groups` - (Optional) An array of one or more security group names
    to associate with the server. Changing this results in adding/removing
    security groups from the existing server. *Note*: When attaching the
    instance to networks using Ports, place the security groups on the Port
    and not the instance. *Note*: Names should be used and not ids, as ids
    trigger unnecessary updates.

* `availability_zone_hints` - (Optional) The availability zone in which to
    create the server. This argument is preferred to `availability_zone`, when
    scheduling the server on a
    [particular](https://docs.openstack.org/nova/latest/admin/availability-zones.html)
    host or node. Conflicts with `availability_zone`. Changing this creates a
    new server.

* `availability_zone` - (Optional) The availability zone in which to create
    the server. Conflicts with `availability_zone_hints`. Changing this creates
    a new server.

* `network` - (Optional) An array of one or more networks to attach to the
    instance. The network object structure is documented below. Changing this
    creates a new server.

* `network_mode` - (Optional) Special string for `network` option to create
  the server. `network_mode` can be `"auto"` or `"none"`.
  Please see the following [reference](https://docs.openstack.org/api-ref/compute/?expanded=create-server-detail#id11) for more information. Conflicts with `network`.

* `metadata` - (Optional) Metadata key/value pairs to make available from
    within the instance. Changing this updates the existing server metadata.

* `config_drive` - (Optional) Whether to use the config_drive feature to
    configure the instance. Changing this creates a new server.

* `admin_pass` - (Optional) The administrative password to assign to the server.
    Changing this changes the root password on the existing server.

* `key_pair` - (Optional) The name of a key pair to put on the server. The key
    pair must already be created and associated with the tenant's account.
    Changing this creates a new server.

* `block_device` - (Optional) Configuration of block devices. The block_device
    structure is documented below. Changing this creates a new server.
    You can specify multiple block devices which will create an instance with
    multiple disks. This configuration is very flexible, so please see the
    following [reference](https://docs.openstack.org/nova/latest/user/block-device-mapping.html)
    for more information.

* `scheduler_hints` - (Optional) Provide the Nova scheduler with hints on how
    the instance should be launched. The available hints are described below.

* `personality` - (Optional) Customize the personality of an instance by
    defining one or more files and their contents. The personality structure
    is described below. Changing this rebuilds the existing server.

* `stop_before_destroy` - (Optional) Whether to try stop instance gracefully
    before destroying it, thus giving chance for guest OS daemons to stop correctly.
    If instance doesn't stop within timeout, it will be destroyed anyway.

* `force_delete` - (Optional) Whether to force the OpenStack instance to be
    forcefully deleted. This is useful for environments that have reclaim / soft
    deletion enabled.

* `power_state` - (Optional) Provide the VM state. Only 'active' and 'shutoff'
    are supported values. *Note*: If the initial power_state is the shutoff
    the VM will be stopped immediately after build and the provisioners like
    remote-exec or files are not supported.

* `tags` - (Optional) A set of string tags for the instance. Changing this
    updates the existing instance tags.

* `vendor_options` - (Optional) Map of additional vendor-specific options.
    Supported options are described below.

The `network` block supports:

* `uuid` - (Required unless `port`  or `name` is provided) The network UUID to
    attach to the server. Changing this creates a new server.

* `name` - (Required unless `uuid` or `port` is provided) The human-readable
    name of the network. Changing this creates a new server.

* `port` - (Required unless `uuid` or `name` is provided) The port UUID of a
    network to attach to the server. Changing this creates a new server.

* `fixed_ip_v4` - (Optional) Specifies a fixed IPv4 address to be used on this
    network. Changing this creates a new server.

* `access_network` - (Optional) Specifies if this network should be used for
    provisioning access. Accepts true or false. Defaults to false.

The `block_device` block supports:

* `uuid` - (Required unless `source_type` is set to `"blank"` ) The UUID of
    the image, volume, or snapshot. Changing this creates a new server.

* `source_type` - (Required) The source type of the device. Must be one of
    "blank", "image", "volume", or "snapshot". Changing this creates a new
    server.

* `volume_size` - The size of the volume to create (in gigabytes). Required
    in the following combinations: source=image and destination=volume,
    source=blank and destination=local, and source=blank and destination=volume.
    Changing this creates a new server.

* `guest_format` - (Optional) Specifies the guest server disk file system format,
    such as `ext2`, `ext3`, `ext4`, `xfs` or `swap`. Swap block device mappings
    have the following restrictions: source_type must be blank and destination_type
    must be local and only one swap disk per server and the size of the swap disk
    must be less than or equal to the swap size of the flavor. Changing this
    creates a new server.

* `boot_index` - (Optional) The boot index of the volume. It defaults to 0.
    Changing this creates a new server.

* `destination_type` - (Optional) The type that gets created. Possible values
    are "volume" and "local". Changing this creates a new server.

* `delete_on_termination` - (Optional) Delete the volume / block device upon
    termination of the instance. Defaults to false. Changing this creates a
    new server.

* `volume_type` - (Optional) The volume type that will be used, for example SSD
    or HDD storage. The available options depend on how your specific OpenStack
    cloud is configured and what classes of storage are provided. Changing this
    creates a new server.

* `device_type` - (Optional) The low-level device type that will be used. Most
    common thing is to leave this empty. Changing this creates a new server.

* `disk_bus` - (Optional) The low-level disk bus that will be used. Most common
    thing is to leave this empty. Changing this creates a new server.

The `scheduler_hints` block supports:

* `group` - (Optional) A UUID of a Server Group. The instance will be placed
    into that group.

* `different_host` - (Optional) A list of instance UUIDs. The instance will
    be scheduled on a different host than all other instances.

* `same_host` - (Optional) A list of instance UUIDs. The instance will be
    scheduled on the same host of those specified.

* `query` - (Optional) A conditional query that a compute node must pass in
    order to host an instance. The query must use the `JsonFilter` syntax
    which is described
    [here](https://docs.openstack.org/nova/latest/admin/configuration/schedulers.html#jsonfilter).
    At this time, only simple queries are supported. Compound queries using
    `and`, `or`, or `not` are not supported. An example of a simple query is:

    ```
    [">=", "$free_ram_mb", "1024"]
    ```

* `target_cell` - (Optional) The name of a cell to host the instance.

* `different_cell` - (Optional) The names of cells where not to build the instance.

* `build_near_host_ip` - (Optional) An IP Address in CIDR form. The instance
    will be placed on a compute node that is in the same subnet.

* `additional_properties` - (Optional) Arbitrary key/value pairs of additional
  properties to pass to the scheduler.

The `personality` block supports:

* `file` - (Required) The absolute path of the destination file.

* `content` - (Required) The contents of the file. Limited to 255 bytes.

The `vendor_options` block supports:

* `ignore_resize_confirmation` - (Optional) Boolean to control whether
    to ignore manual confirmation of the instance resizing. This can be helpful
    to work with some OpenStack clouds which automatically confirm resizing of
    instances after some timeout.

* `detach_ports_before_destroy` - (Optional) Whether to try to detach all attached
    ports to the vm before destroying it to make sure the port state is correct
    after the vm destruction. This is helpful when the port is not deleted.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `access_ip_v4` - The first detected Fixed IPv4 address.
* `access_ip_v6` - The first detected Fixed IPv6 address.
* `metadata` - See Argument Reference above.
* `security_groups` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `flavor_name` - See Argument Reference above.
* `network/uuid` - See Argument Reference above.
* `network/name` - See Argument Reference above.
* `network/port` - See Argument Reference above.
* `network/fixed_ip_v4` - The Fixed IPv4 address of the Instance on that
    network.
* `network/fixed_ip_v6` - The Fixed IPv6 address of the Instance on that
    network.
* `network/mac` - The MAC address of the NIC on that network.
* `all_metadata` - Contains all instance metadata, even metadata not set
    by Terraform.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the instance, which have
    been explicitly and implicitly added.
* `created` - The creation time of the instance.
* `updated` - The time when the instance was last updated.

## Notes

### Multiple Ephemeral Disks

It's possible to specify multiple `block_device` entries to create an instance
with multiple ephemeral (local) disks. In order to create multiple ephemeral
disks, the sum of the total amount of ephemeral space must be less than or
equal to what the chosen flavor supports.

The following example shows how to create an instance with multiple ephemeral
disks:

```hcl
resource "openstack_compute_instance_v2" "foo" {
  name            = "terraform-test"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image uuid>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instances and Security Groups

When referencing a security group resource in an instance resource, always
use the _name_ of the security group. If you specify the ID of the security
group, Terraform will remove and reapply the security group upon each call.
This is because the OpenStack Compute API returns the names of the associated
security groups and not their IDs.

Note the following example:

```hcl
resource "openstack_networking_secgroup_v2" "sg_1" {
  name = "sg_1"
}

resource "openstack_compute_instance_v2" "foo" {
  name            = "terraform-test"
  security_groups = ["${openstack_networking_secgroup_v2.sg_1.name}"]
}
```

### Instances and Ports

Neutron Ports are a great feature and provide a lot of functionality. However,
there are some notes to be aware of when mixing Instances and Ports:

* In OpenStack environments prior to the Kilo release, deleting or recreating
an Instance will cause the Instance's Port(s) to be deleted. One way of working
around this is to taint any Port(s) used in Instances which are to be recreated.
See [here](https://review.openstack.org/#/c/126309/) for further information.

* When attaching an Instance to one or more networks using Ports, place the
security groups on the Port and not the Instance. If you place the security
groups on the Instance, the security groups will not be applied upon creation,
but they will be applied upon a refresh. This is a known OpenStack bug.

* Network IP information is not available within an instance for networks that
are attached with Ports. This is mostly due to the flexibility Neutron Ports
provide when it comes to IP addresses. For example, a Neutron Port can have
multiple Fixed IP addresses associated with it. It's not possible to know which
single IP address the user would want returned to the Instance's state
information. Therefore, in order for a Provisioner to connect to an Instance
via it's network Port, customize the `connection` information:

```hcl
resource "openstack_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"

  network_id = "0a1d0a27-cffa-4de3-92c5-9d3fd3f2e74d"

  security_group_ids = [
    "2f02d20a-8dca-49b7-b26f-b6ce9fddaf4f",
    "ca1e5ed7-dae8-4605-987b-fadaeeb30461",
  ]
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"

  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }

  connection {
    user        = "root"
    host        = "${openstack_networking_port_v2.port_1.fixed_ip.0.ip_address}"
    private_key = "~/path/to/key"
  }

  provisioner "remote-exec" {
    inline = [
      "echo terraform executed > /tmp/foo",
    ]
  }
}
```

### Instances and Networks

Instances almost always require a network. Here are some notes to be aware of
with how Instances and Networks relate:

* In scenarios where you only have one network available, you can create an
instance without specifying a `network` block. OpenStack will automatically
launch the instance on this network.

* If you have access to more than one network, you will need to specify a network
with a `network` block. Not specifying a network will result in the following
error:

```
* openstack_compute_instance_v2.instance: Error creating OpenStack server:
Expected HTTP response code [201 202] when accessing [POST https://example.com:8774/v2.1/servers], but got 409 instead
{"conflictingRequest": {"message": "Multiple possible networks found, use a Network ID to be more specific.", "code": 409}}
```

* If you intend to use the `openstack_compute_interface_attach_v2` resource,
you still need to make sure one of the above points is satisfied. An instance
cannot be created without a valid network configuration even if you intend to
use `openstack_compute_interface_attach_v2` after the instance has been created.

## Importing instances

Importing instances can be tricky, since the nova api does not offer all
information provided at creation time for later retrieval.
Network interface attachment order, and number and sizes of ephemeral
disks are examples of this.

### Importing basic instance
Assume you want to import an instance with one ephemeral root disk,
and one network interface.

Your configuration would look like the following:

```hcl
resource "openstack_compute_instance_v2" "basic_instance" {
  name            = "basic"
  flavor_id       = "<flavor_id>"
  key_pair        = "<keyname>"
  security_groups = ["default"]
  image_id =  "<image_id>"

  network {
    name = "<network_name>"
  }
}

```
Then you execute
```
terraform import openstack_compute_instance_v2.basic_instance <instance_id>
```

### Importing an instance with multiple emphemeral disks

The importer cannot read the emphemeral disk configuration
of an instance, so just specify image_id as in the configuration
of the basic instance example.

### Importing instance with multiple network interfaces.

Nova returns the network interfaces grouped by network, thus not in creation
order.
That means that if you have multiple network interfaces you must take
care of the order of networks in your configuration.


As example we want to import an instance with one ephemeral root disk,
and 3 network interfaces.

Examples

```hcl
resource "openstack_compute_instance_v2" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "<flavor_id"
  key_pair        = "<keyname>"
  image_id        = "<image_id>"
  security_groups = ["default"]

  network {
    name = "<network1>"
  }
  network {
    name = "<network2>"
  }
  network {
    name = "<network1>"
    fixed_ip_v4 = "<fixed_ip_v4>"
  }

}
```

In the above configuration the networks are out of order compared to what nova
and thus the import code returns, which means the plan will not
be empty after import.

So either with care check the plan and modify configuration, or read the
network order in the state file after import and modify your
configuration accordingly.

 * A note on ports. If you have created a neutron port independent of an
 instance, then the import code has no way to detect that the port is created
 idenpendently, and therefore on deletion of imported instances you might have
 port resources in your project, which you expected to be created by the
 instance and thus to also be deleted with the instance.


### Importing instances with multiple block storage volumes.

We have an instance with two block storage volumes, one bootable and one
non-bootable.
Note that we only configure the bootable device as block_device.
The other volumes can be specified as `openstack_blockstorage_volume_v2`

```hcl
resource "openstack_compute_instance_v2" "instance_2" {
  name            = "instance_2"
  image_id        = "<image_id>"
  flavor_id       = "<flavor_id>"
  key_pair        = "<keyname>"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image_id>"
    source_type           = "image"
    destination_type      = "volume"
    boot_index            = 0
    delete_on_termination = true
  }

   network {
    name = "<network_name>"
  }
}
resource "openstack_blockstorage_volume_v2" "volume_1" {
  size = 1
  name = "<vol_name>"
}
resource "openstack_compute_volume_attach_v2" "va_1" {
  volume_id   = "${openstack_blockstorage_volume_v2.volume_1.id}"
  instance_id = "${openstack_compute_instance_v2.instance_2.id}"
}
```
To import the instance outlined in the above configuration
do the following:

```
terraform import openstack_compute_instance_v2.instance_2 <instance_id>
import openstack_blockstorage_volume_v2.volume_1 <volume_id>
terraform import openstack_compute_volume_attach_v2.va_1
<instance_id>/<volume_id>
```

* A note on block storage volumes, the importer does not read
  delete_on_termination flag, and always assumes true. If you
  import an instance created with delete_on_termination false,
  you end up with "orphaned" volumes after destruction of
  instances.
