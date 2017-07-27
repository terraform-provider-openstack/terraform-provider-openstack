## 0.1.1 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Network detection in the `openstack_compute_instance_v2` resource was cleaned up and updated. There should be no incompatibilities, but you should do a `plan` before `apply` just to be safe.

IMPROVEMENTS:

* `openstack_compute_instance_v2` now supports Neutron for network detection [GH-39]
* `openstack_compute_instance_v2` support for multiple NICs on the same network [GH-39]
* Added support for `TERMINATED_HTTPS` protocol in `openstack_lb_listener_v2` [GH-49]

BUG FIXES:
* `openstack_lb_pool_v2` handling of `persistence` updated, `cookie_name` is now optional. [GH-57]

## 0.1.0 (June 21, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* You can now specify `region` in the provider block. All resources will inherit this region setting, or you can override it in the resource-level `region`. Make sure to do a `plan` before an `apply` to make sure the resource is not destroyed due to incorrectly determining the region! If you see this happening, either explicitly set the `region` in the resource or use `lifecycle.ignore_changes`. 
* `floating_ip` has been removed from `openstack_compute_instance_v2`. You must now use `openstack_compute_floatingip_associate_v2` to associate a Floating IP with an Instance.
* `volume` has been removed from `openstack_compute_instance_v2`. You must now use `openstack_compute_volume_attach_v2` to attach a Volume with an Instance.
* `member` has been removed from `openstack_lb_pool_v1`. You must now use `openstack_lb_member_v1` to add a LBaaS v1 Member to a Pool.


IMPROVEMENTS:

* Can specify `region` in the provider ([#25](https://github.com/terraform-providers/terraform-provider-openstack/issues/25))

BUG FIXES

* Wait for LoadBalancer to be active before creating Pools and Monitors ([#29](https://github.com/terraform-providers/terraform-provider-openstack/issues/29))
* Choose first network found with a matching name for compute instances ([#36](https://github.com/terraform-providers/terraform-provider-openstack/issues/36))
