## 1.2.1 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `openstack_compute_keypair_v2` can now generate a private key, however the private key will be stored in your Terraform state. Please use caution.
* The MAC addresses in `openstack_networking_port_v2.allowed_address_pairs` is no longer computed. This should not cause an issue for users since if an `allowed_address_pairs` MAC address was not specified, the AAP MAC will match `openstack_networking_port_v2.mac_address` [GH-244]

FEATURES

* __New Resource:__ `openstack_networking_subnetpool_v2` [GH-243]
* __New Resource:__ `openstack_identity_role_v3` [GH-250]
* __New Data Source:__ `openstack_networking_subnetpool_v2` [GH-243]
* __New Data Source:__ `openstack_identity_role_v3` [GH-250]

IMPROVEMENTS

* Added `additional_properties` to `openstack_compute_instance_v2` scheduler hints [GH-230]
* `openstack_compute_keypair_v2` can now generate a private key [GH-217]
* `openstack_networking_router_v2` can now optionally set a default gateway after it has been created [GH-209]
* Added `subnetpool_id` to `openstack_networking_subnet_v2` resource and data source [GH-249]
* Added `extra_specs` to `openstack_compute_flavor_v2` [GH-241]

BUG FIXES

* Fixed bug with `openstack_networking_network_v2` and `openstack_networking_subnet_v2` where the `OS_TENANT_ID` was incorrectly being used as a default value [GH-254]
* Correctly detect if an object storage container is deleted [GH-261]
* Fixed a few small bugs with `openstack_fw_rule_v1` updating [GH-224]
* Fixed an issue with `openstack_networking_port_v2` `allowed_address_pairs` and MAC addresses [GH-244]

## 1.2.0 (January 18, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* The way IP addresses for `allowed_address_pairs` in the `openstack_networking_port_v2` resource are stored in the Terraform state has changed. 
* The `external_gateway` argument in the `openstack_networking_router_v2` has been deprecated in favor of the more appropriately named `external_network_id`.

FEATURES

* __New Resource:__ `openstack_db_database_v1` ([#179](https://github.com/terraform-providers/terraform-provider-openstack/issues/179))
* __New Resource:__ `openstack_db_user_v1` ([#180](https://github.com/terraform-providers/terraform-provider-openstack/issues/180))
* __New Resource:__ `openstack_db_configuration_v1` ([#185](https://github.com/terraform-providers/terraform-provider-openstack/issues/185))
* __New Data Source:__ `openstack_compute_flavor_v2` ([#190](https://github.com/terraform-providers/terraform-provider-openstack/issues/190))


IMPROVEMENTS

* Added `external_fixed_ips` to the `openstack_networking_router_v2` resource ([#178](https://github.com/terraform-providers/terraform-provider-openstack/issues/178))
* Added `ipv6_address_mode` and `ipv6_ra_mode` to the `openstack_networking_subnet_v2` resource and data source ([#193](https://github.com/terraform-providers/terraform-provider-openstack/issues/193))
* Several new `openstack_networking_subnet_v2` attributes are now accessible in the data source ([#199](https://github.com/terraform-providers/terraform-provider-openstack/issues/199))
* Added `availability_zone_hints` to the `openstack_networking_network_v2` resource and data source ([#196](https://github.com/terraform-providers/terraform-provider-openstack/issues/196))
* Added `availability_zone_hints` to the `openstack_networking_router_v2` resource ([#203](https://github.com/terraform-providers/terraform-provider-openstack/issues/203))
* User's password field in `openstack_db_instance_v2` resource has been marked sensitive ([#220](https://github.com/terraform-providers/terraform-provider-openstack/issues/220))
* `openstack_db_instance_v1` now supports setting a `configuration_id` ([#221](https://github.com/terraform-providers/terraform-provider-openstack/issues/221))

BUG FIXES

* Allow the same `ip_address` with a different `mac_address` to be specified multiple times in the `openstack_networking_port_v2` resource ([#168](https://github.com/terraform-providers/terraform-provider-openstack/issues/168))
* Fixed unhandled error checks which were causing crashes in `openstack_networking_secgroup_v2` and `openstack_networking_network_v2` data sources ([#201](https://github.com/terraform-providers/terraform-provider-openstack/issues/201))
* Fixed unhandled error check when creating `openstack_networking_floatingip_v2` ([#206](https://github.com/terraform-providers/terraform-provider-openstack/issues/206))
* Fixed region detection when using `clouds.yaml` ([#216](https://github.com/terraform-providers/terraform-provider-openstack/issues/216))
* Make `subnet_id` optional for `openstack_lb_member_v2` ([#189](https://github.com/terraform-providers/terraform-provider-openstack/issues/189))
* Fix ordering of DNS servers in `openstack_networking_subnet_v2` ([#226](https://github.com/terraform-providers/terraform-provider-openstack/issues/226))

## 1.1.0 (December 04, 2017)

FEATURES

* __New Resource:__ `openstack_objectstorage_object_v1` ([#146](https://github.com/terraform-providers/terraform-provider-openstack/issues/146))
* __New Resource:__ `openstack_db_instance_v1` ([#155](https://github.com/terraform-providers/terraform-provider-openstack/issues/155))

IMPROVEMENTS

* Better handling of mutually exclusive options `no_gateway` and `gateway_ip` in the `openstack_networking_subnet_v2` resource ([#136](https://github.com/terraform-providers/terraform-provider-openstack/issues/136))
* Can now authenticate with a `clouds.yaml` file ([#154](https://github.com/terraform-providers/terraform-provider-openstack/issues/154))

BUG FIXES

* Fixed issue with automatic detection of an Octavia client and Networking client ([#172](https://github.com/terraform-providers/terraform-provider-openstack/issues/172))
* Fixed issue with creating public flavors ([#177](https://github.com/terraform-providers/terraform-provider-openstack/issues/177))

## 1.0.0 (November 08, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* If your OpenStack cloud supports the Octavia Load Balancing service, you can now use it by setting the provider-level `use_octavia` argument to `true`. The `openstack_lb_*_v2` resources will then seamlessly use Octavia.

FEATURES

* __New Data Source:__ `openstack_networking_subnet_v2` ([#135](https://github.com/terraform-providers/terraform-provider-openstack/issues/135))
* __New Data Source:__ `openstack_dns_zone_v2` ([#145](https://github.com/terraform-providers/terraform-provider-openstack/issues/145))

IMPROVEMENTS

* `openstack_networking_router_v2`: Added `enable_snat` argument ([#140](https://github.com/terraform-providers/terraform-provider-openstack/issues/140))
* Added provider-level option of `use_octavia` to use the Octavia load balancing service ([#149](https://github.com/terraform-providers/terraform-provider-openstack/issues/149))

## 0.3.0 (October 23, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* The `openstack_networking_port_v2` resource had a significant update to how it handles security groups. If you have not explicitly defined security groups in the port resource, any security groups which were automatically applied by OpenStack (such as the `default` security group) will be removed upon the next apply. To prevent this from happening, add the ID of the security groups to the `security_group_ids` argument. If you are already explicitly specifying security groups, you should see no change in behavior.

IMPROVEMENTS

 * `openstack_networking_router_interface_v2` will now set `subnet_id` when importing ([#119](https://github.com/terraform-providers/terraform-provider-openstack/issues/119))
 * `openstack_networking_router_route_v2` can now be imported ([#120](https://github.com/terraform-providers/terraform-provider-openstack/issues/120))
 * `openstack_images_image_v2` resource and data source now supports reading and setting properties ([#113](https://github.com/terraform-providers/terraform-provider-openstack/issues/113))

BUG FIXES

  * `openstack_networking_port_v2`: Fixed issues with how security groups and allowed address pairs are applied and updated [[#114](https://github.com/terraform-providers/terraform-provider-openstack/issues/114)].

## 0.2.2 (September 15, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Unused `id` fields in the LBaaS v2 resources were removed. This should not cause any issues, but please report if you find otherwise.

FEATURES:

* __New Data Source:__ `openstack_networking_secgroup_v2` ([#86](https://github.com/terraform-providers/terraform-provider-openstack/issues/86))
* __New Resource:__: `openstack_compute_flavor_v2` ([#83](https://github.com/terraform-providers/terraform-provider-openstack/issues/83))

IMPROVEMENTS
 * Added `status` field to `openstack_networking_network_v2` data source ([#105](https://github.com/terraform-providers/terraform-provider-openstack/issues/105))
 * `openstack_networking_router_v2` can now be imported ([#111](https://github.com/terraform-providers/terraform-provider-openstack/issues/111))
 * `openstack_networking_router_interface_v2` can now be imported ([#112](https://github.com/terraform-providers/terraform-provider-openstack/issues/112))
 
BUG FIXES

* `openstack_lb_listener_v2`: Don't send `connection_limit` unless it has been set ([#90](https://github.com/terraform-providers/terraform-provider-openstack/issues/90))
* `openstack_lb_pool_v2`: Find Load Balancer via Listener ([#97](https://github.com/terraform-providers/terraform-provider-openstack/issues/97))
* LBaaS v2: Removed unused `id` fields ([#93](https://github.com/terraform-providers/terraform-provider-openstack/issues/93))
* `openstack_lb_monitor_v2`: Check if a monitor was successfully created before proceeding ([#102](https://github.com/terraform-providers/terraform-provider-openstack/issues/102))
* `openstack_networking_router_v2`: Fix region parameter ([#107](https://github.com/terraform-providers/terraform-provider-openstack/issues/107))
* `openstack_compute_instance_v2`: Fix regression bug with NIC detection ([#117](https://github.com/terraform-providers/terraform-provider-openstack/issues/117))

## 0.2.1 (August 23, 2017)

IMPROVEMENTS:

* `openstack_lb_loadbalancer_v2` timeouts have been lowered to 10 and 5 minutes ([#74](https://github.com/terraform-providers/terraform-provider-openstack/issues/74))

BUG FIXES:

* `openstack_images_image_v2` data source now sorts images by `CreatedAt` instead of `UpdatedAt` ([#78](https://github.com/terraform-providers/terraform-provider-openstack/issues/78))
* `openstack_networking_secgroup_v2` now re-reads security group before deleteing rules when `delete_default_rules => true` ([#82](https://github.com/terraform-providers/terraform-provider-openstack/issues/82))
* Fixed `openstack_compute_instance_v2` access IP address detection in dual-stack environments ([#85](https://github.com/terraform-providers/terraform-provider-openstack/issues/85))

## 0.2.0 (August 14, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Network detection in the `openstack_compute_instance_v2` resource was cleaned up and updated. There should be no incompatibilities, but you should do a `plan` before `apply` just to be safe.
* The `openstack_lb_loadbalancer_v2.provider` argument has been removed entirely. This was an erroneous argument from the beginning, so it should not be in use. However, if you do have it set in your configurations, please rename it to `loadbalancer_provider`.

FEATURES:

* __New Resource:__ `openstack_identity_project_v3` ([#50](https://github.com/terraform-providers/terraform-provider-openstack/issues/50))
* __New Resource:__ `openstack_identity_user_v3` ([#52](https://github.com/terraform-providers/terraform-provider-openstack/issues/52))

IMPROVEMENTS:

* `openstack_compute_instance_v2` now supports Neutron for network detection ([#39](https://github.com/terraform-providers/terraform-provider-openstack/issues/39))
* `openstack_compute_instance_v2` support for multiple NICs on the same network ([#39](https://github.com/terraform-providers/terraform-provider-openstack/issues/39))
* Added support for `TERMINATED_HTTPS` protocol in `openstack_lb_listener_v2` ([#49](https://github.com/terraform-providers/terraform-provider-openstack/issues/49))
* Improvements to LBaaS v2 resource coordination ([#59](https://github.com/terraform-providers/terraform-provider-openstack/issues/59))
* `openstack_lb_loadbalancer_v2.provider` has been removed. See notes above. ([#65](https://github.com/terraform-providers/terraform-provider-openstack/issues/65))

BUG FIXES:
* `openstack_lb_pool_v2` handling of `persistence` updated, `cookie_name` is now optional. ([#57](https://github.com/terraform-providers/terraform-provider-openstack/issues/57))
* `openstack_fw_firewall_v1.associated_routers` is now computed. ([#53](https://github.com/terraform-providers/terraform-provider-openstack/issues/53))
* All `openstack_fw_rule_v1` attributes are now passed during an update phase. ([#53](https://github.com/terraform-providers/terraform-provider-openstack/issues/53))
* `openstack_networking_secgroup_v2` now correctly updates description. ([#60](https://github.com/terraform-providers/terraform-provider-openstack/issues/60))
* `openstack_fw_firewall_v1` now correctly translates `value_specs` on create. ([#66](https://github.com/terraform-providers/terraform-provider-openstack/issues/66))

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
