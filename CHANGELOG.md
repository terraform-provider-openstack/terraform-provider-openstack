## 1.18.0 (Unreleased)
## 1.17.0 (April 01, 2019)

NOTES

* `extra_dhcp_option` in the `openstack_networking_port_v2` data source has been changed to a List. This is to resolve a bug where multiple DHCP options were not being rendered.


FEATURES

* __New Resource__: `openstack_identity_application_credential_v3` ([#660](https://github.com/terraform-providers/terraform-provider-openstack/issues/660))
* __New Data Source__: `openstack_blockstorage_availability_zones_v3` ([#652](https://github.com/terraform-providers/terraform-provider-openstack/issues/652))
* __New Data Source__: `openstack_sharedfilesystem_availability_zones_v2` ([#652](https://github.com/terraform-providers/terraform-provider-openstack/issues/652))
* __New Data Source__: `openstack_networking_trunk_v2` ([#626](https://github.com/terraform-providers/terraform-provider-openstack/issues/626))

IMPROVEMENTS

* Reduced API calls when updating `extra_dhcp_option` in `openstack_networking_port_v2` ([#689](https://github.com/terraform-providers/terraform-provider-openstack/issues/689))
* Added `port_security_enabled` to `openstack_networking_network_v2` ([#681](https://github.com/terraform-providers/terraform-provider-openstack/issues/681))
* Added `port_security_enabled` to `openstack_networking_port_v2` ([#682](https://github.com/terraform-providers/terraform-provider-openstack/issues/682))
* Added `prefix_length` to `openstack_networking_subnet_v2` ([#705](https://github.com/terraform-providers/terraform-provider-openstack/issues/705))
* Added `binding` to `openstack_networking_port_v2` ([#693](https://github.com/terraform-providers/terraform-provider-openstack/issues/693))
* Added `binding` to `openstack_networking_port_v2` data source ([#693](https://github.com/terraform-providers/terraform-provider-openstack/issues/693))
* Added `mtu` to `openstack_networking_network_v2` ([#708](https://github.com/terraform-providers/terraform-provider-openstack/issues/708))
* Added `mtu` to `openstack_networking_network_v2` data source ([#708](https://github.com/terraform-providers/terraform-provider-openstack/issues/708))
* Added `dns_name` and `dns_domain` to `openstack_networking_floatingip_v2` ([#706](https://github.com/terraform-providers/terraform-provider-openstack/issues/706))
* Added `dns_name` and `dns_domain` to `openstack_networking_floatingip_v2` data source ([#706](https://github.com/terraform-providers/terraform-provider-openstack/issues/706))
* Added `dns_domain` to `openstack_networking_network_v2` ([#706](https://github.com/terraform-providers/terraform-provider-openstack/issues/706))
* Added `dns_domain` to `openstack_networking_network_v2` data source ([#706](https://github.com/terraform-providers/terraform-provider-openstack/issues/706))
* Added `dns_name` and `dns_assignment` to `openstack_networking_port_v2` ([#706](https://github.com/terraform-providers/terraform-provider-openstack/issues/706))
* Added `dns_name` and `dns_assignment` to `openstack_networking_port_v2` data source ([#706](https://github.com/terraform-providers/terraform-provider-openstack/issues/706))
* Added `fixed_ip` to `openstack_networking_floatingip_associate_v2` ([#709](https://github.com/terraform-providers/terraform-provider-openstack/issues/709))
* Enable `fixed_ip` to be updated in `openstack_networking_floatingip_v2` ([#709](https://github.com/terraform-providers/terraform-provider-openstack/issues/709))
* Added ability to specify `cephx` as `access_type` and to retrieve the `access_key` in `openstack_sharedfilesystem_share_access_v2` ([#715](https://github.com/terraform-providers/terraform-provider-openstack/issues/715))

BUG FIXES

* Fixed bug in `openstack_identity_auth_scope_v3` data source where the `user_id` attribute was being set to the user's Name and not ID ([#660](https://github.com/terraform-providers/terraform-provider-openstack/issues/660))
* Fixed bug in Load Balancer resources for Contrail-based load balancers ([#691](https://github.com/terraform-providers/terraform-provider-openstack/issues/691))
* `extra_dhcp_option` in the `openstack_networking_port_v2` data source has been changed to a List. This is to resolve a bug where multiple DHCP options were not being rendered ([#695](https://github.com/terraform-providers/terraform-provider-openstack/issues/695))

## 1.16.0 (February 14, 2019)

NOTES

* The `openstack_networking_subnet_v2.host_routes` argument has been marked as deprecated. Please use the dedicated `openstack_networking_subnet_route_v2` resource instead.

FEATURES

* __New Data Source__: `openstack_compute_availability_zones_v2` ([#655](https://github.com/terraform-providers/terraform-provider-openstack/issues/655))

BUG FIXES

* The `openstack_networking_subnet_v2.host_routes` argument has been deprecated due to schema issues and conflicts with `openstack_networking_subnet_route_v2` ([#668](https://github.com/terraform-providers/terraform-provider-openstack/issues/668))
* A previously added validation to `openstack_networking_port_v2.fixed_ip.ip_address` was removed as it was causing problems for prior behavior of using an empty string ([#678](https://github.com/terraform-providers/terraform-provider-openstack/issues/678))

## 1.15.1 (February 08, 2019)

BUG FIXES

* Fixed issue where volume multiattachments would not be retried ([#540](https://github.com/terraform-providers/terraform-provider-openstack/issues/540))
* Reverted an incorrect schema validation for `openstack_networking_port_v2.allowed_address_pairs` ([#661](https://github.com/terraform-providers/terraform-provider-openstack/issues/661))

## 1.15.0 (February 06, 2019)

NOTES

* The `openstack_images_image_v2.update_at` attribute has been deprecated in favor of `updated_at` ([#617](https://github.com/terraform-providers/terraform-provider-openstack/issues/617))

FEATURES

* __New Resource__: `openstack_networking_addressscope_v2` ([#634](https://github.com/terraform-providers/terraform-provider-openstack/issues/634))
* __New Resource__: `openstack_networking_port_secgroup_associate_v2` ([#574](https://github.com/terraform-providers/terraform-provider-openstack/issues/574))

IMPROVEMENTS

* Added `flavor_id` to the `openstack_compute_flavor_v2` data source so flavors can be queried by ID ([#587](https://github.com/terraform-providers/terraform-provider-openstack/issues/587))
* `openstack_networking_port_ids_v2` data source can now return an empty set of results ([#631](https://github.com/terraform-providers/terraform-provider-openstack/issues/631))
* Added `description` to `openstack_networking_trunk_v2` resource ([#625](https://github.com/terraform-providers/terraform-provider-openstack/issues/625))
* Added `tags` to the networking data source to query by tags and `all_tags` to see a full list of tags ([#624](https://github.com/terraform-providers/terraform-provider-openstack/issues/624))
* `openstack_compute_instance_v2.admin_pass` is now a "sensitive" attribute ([#647](https://github.com/terraform-providers/terraform-provider-openstack/issues/647))
* Added support to authenticate with Application Credentials ([#642](https://github.com/terraform-providers/terraform-provider-openstack/issues/642))
* Added ability to specify region in `openstack_sharedfilesystem_share_access_v2` ([#654](https://github.com/terraform-providers/terraform-provider-openstack/issues/654))
* Added ability to specify region in `openstack_sharedfilesystem_share_v2` ([#654](https://github.com/terraform-providers/terraform-provider-openstack/issues/654))
* Added `all_tags` attribute to Networking resources to set tags provided by the OpenStack backend automatically ([#623](https://github.com/terraform-providers/terraform-provider-openstack/issues/623))

BUG FIXES

* Fixed `created_at`, `updated_at`, and `tag` fields in the `openstack_images_image_v2` data source ([#615](https://github.com/terraform-providers/terraform-provider-openstack/issues/615))
* Fixed `created_at` and `updated_at` fields in the `openstack_networking_subnetpool_v2` resource ([#619](https://github.com/terraform-providers/terraform-provider-openstack/issues/619))
* Fixed `created_at` and `updated_at` fields in the `openstack_networking_subnetpool_v2` data source ([#616](https://github.com/terraform-providers/terraform-provider-openstack/issues/616))
* Fixed issue where updating the description of a floating IP would cause the port to disassociate ([#606](https://github.com/terraform-providers/terraform-provider-openstack/issues/606))
* `admin_state_up` and `shared` fields of `openstack_networking_network_v2` are now correct boolean fields ([#593](https://github.com/terraform-providers/terraform-provider-openstack/issues/593))
* `external` field of `openstack_networking_network_v2` field will now show an actual value ([#593](https://github.com/terraform-providers/terraform-provider-openstack/issues/593))
* Fixed issue where `status` was being used as the query value for `network_id` in `openstack_networking_port_v2` data source ([#631](https://github.com/terraform-providers/terraform-provider-openstack/issues/631))
* Fixed issue where `status` was being used as the query value for `network_id` in `openstack_networking_port_ids_v2` data source ([#631](https://github.com/terraform-providers/terraform-provider-openstack/issues/631))
* The `openstack_images_image_v2` fields `update_at`, `updated_at`, and `created_at` all now set correctly ([#617](https://github.com/terraform-providers/terraform-provider-openstack/issues/617))
* Fixed issue with `openstack_dns_recordset_v2` where `records` would be returned out of order ([#636](https://github.com/terraform-providers/terraform-provider-openstack/issues/636))
* Fixed issue where `openstack_compute_volume_attach_v2` and `openstack_blockstorage_volume_v2` were trying to detach volumes at the same time ([#640](https://github.com/terraform-providers/terraform-provider-openstack/issues/640))
* Fixed a regression bug where destroying networks was failing on a 409 code ([#644](https://github.com/terraform-providers/terraform-provider-openstack/issues/644))
* Fixed an issue with `openstack_compute_instance_v2` where a 404 was triggering an error ([#647](https://github.com/terraform-providers/terraform-provider-openstack/issues/647))
* Fixed an issue where `all_fixed_ips` was not being set in `openstack_networking_port_v2` data source ([#649](https://github.com/terraform-providers/terraform-provider-openstack/issues/649))
* Fixed an issue where `openstack_networking_port_v2` would cause an API error ([#649](https://github.com/terraform-providers/terraform-provider-openstack/issues/649))
* Fixed an issue where Blockstorage volume resources couldn't be detached because they had been removed ([#641](https://github.com/terraform-providers/terraform-provider-openstack/issues/641))

## 1.14.0 (January 15, 2019)

NOTES

* The Load Balancer v2 resources have been updated to provide more efficient status checks. If you encounter any problems due to this, please report them and we will make it a priority to resolve.
* `openstack_networking_port_v2` will now set the `admin_state_up` to `true/UP` if it is left omitted from the resource configuration. This now correctly conforms to the OpenStack API. This should be a transparent change, but let us know if this causes you problems.

FEATURES

* __New Resource__: `openstack_lb_l7policy_v2` ([#527](https://github.com/terraform-providers/terraform-provider-openstack/issues/527))
* __New Resource__: `openstack_lb_l7rule_v2` ([#522](https://github.com/terraform-providers/terraform-provider-openstack/issues/522))
* __New Resource__: `openstack_sharedfilesystem_share_v2` ([#525](https://github.com/terraform-providers/terraform-provider-openstack/issues/525))
* __New Resource__: `openstack_sharedfilesystem_share_access_v2` ([#526](https://github.com/terraform-providers/terraform-provider-openstack/issues/526))
* __New Data Source__: `openstack_sharedfilesystem_share_v2` ([#564](https://github.com/terraform-providers/terraform-provider-openstack/issues/564))
* __New Data Source__: `openstack_networking_port_v2` ([#567](https://github.com/terraform-providers/terraform-provider-openstack/issues/567))
* __New Data Source__: `openstack_sharedfilesystem_sharenetwork_v2` ([#576](https://github.com/terraform-providers/terraform-provider-openstack/issues/576))
* __New Data Source__: `openstack_networking_port_ids_v2` ([#569](https://github.com/terraform-providers/terraform-provider-openstack/issues/569))
* __New Data Source__: `openstack_sharedfilesystem_snapshot_v2` ([#577](https://github.com/terraform-providers/terraform-provider-openstack/issues/577))

IMPROVEMENTS

* Provider options `swauth` and `use_octavia` will correctly use a default value of `false` when they are not specified. This is to help with compatibility for v0.12 ([#494](https://github.com/terraform-providers/terraform-provider-openstack/issues/494))
* Enhanced the pending status checks of the Load Balancer v2 resources ([#550](https://github.com/terraform-providers/terraform-provider-openstack/issues/550))
* Prioritized the status of Load Balancer v2 resources to first use the Load Balancer's master status ([#556](https://github.com/terraform-providers/terraform-provider-openstack/issues/556))
* Fix flavor detection in `openstack_compute_instance_v2` and `openstack_containerinfra_cluster_v1` for Terraform v0.12 ([#551](https://github.com/terraform-providers/terraform-provider-openstack/issues/551))
* Added the ability to import `openstack_lb_loadbalancer_v2` ([#524](https://github.com/terraform-providers/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_listener_v2` ([#524](https://github.com/terraform-providers/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_pool_v2` ([#524](https://github.com/terraform-providers/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_member_v2` ([#524](https://github.com/terraform-providers/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_monitor_v2` ([#524](https://github.com/terraform-providers/terraform-provider-openstack/issues/524))
* Added `device_type` and `disk_bus` to `openstack_compute_instance_v2` block device ([#558](https://github.com/terraform-providers/terraform-provider-openstack/issues/558))
* Added `transparent_vlan` to `openstack_networking_network_v2` ([#513](https://github.com/terraform-providers/terraform-provider-openstack/issues/513))
* Added `transparent_vlan` to `openstack_networking_network_v2` data source ([#538](https://github.com/terraform-providers/terraform-provider-openstack/issues/538))
* Added `max_retries` to the provider options ([#413](https://github.com/terraform-providers/terraform-provider-openstack/issues/413))
* Added the ability to override catalog endpoints ([#501](https://github.com/terraform-providers/terraform-provider-openstack/issues/501))
* Changed the `segments` attribute of the `openstack_networking_network_v2` to `TypeSet` [[#578](https://github.com/terraform-providers/terraform-provider-openstack/issues/578)] 

BUG FIXES

* `openstack_compute_interface_attach_v2` now correctly sets the `instance_id` [[#557](https://github.com/terraform-providers/terraform-provider-openstack/issues/557)] 
* `openstack_networking_port_v2` will now correctly set the `admin_state_up` to `true/UP` if left omitted ([#594](https://github.com/terraform-providers/terraform-provider-openstack/issues/594))
* Fixed out of range panic in `openstack_compute_instance_v2` when no IP addresses were detected ([#539](https://github.com/terraform-providers/terraform-provider-openstack/issues/539))

## 1.13.0 (December 18, 2018)

FEATURES

* __New Resource__: `openstack_sharedfilesystem_securityservice_v2` ([#515](https://github.com/terraform-providers/terraform-provider-openstack/issues/515))
* __New Resource__: `openstack_sharedfilesystem_sharenetwork_v2` ([#515](https://github.com/terraform-providers/terraform-provider-openstack/issues/515))
* __New Data Source__: `openstack_containerinfra_cluster_v1` ([#488](https://github.com/terraform-providers/terraform-provider-openstack/issues/488))
* __New Data Source__: `openstack_blockstorage_snapshot_v2` ([#448](https://github.com/terraform-providers/terraform-provider-openstack/issues/448))
* __New Data Source__: `openstack_blockstorage_snapshot_v3` ([#448](https://github.com/terraform-providers/terraform-provider-openstack/issues/448))

IMPROVEMENTS

* Added object versioning to `openstack_objectstorage_container_v1` ([#465](https://github.com/terraform-providers/terraform-provider-openstack/issues/465))
* Added support for soft affinities in `openstack_compute_servergroup_v2` ([#490](https://github.com/terraform-providers/terraform-provider-openstack/issues/490))
* Allow `default_pool_id` to be updated in `openstack_lb_listener_v2` ([#516](https://github.com/terraform-providers/terraform-provider-openstack/issues/516))
* Added `description` to `openstack_networking_router_v2` ([#529](https://github.com/terraform-providers/terraform-provider-openstack/issues/529))
* Added `description` to `openstack_networking_port_v2` ([#531](https://github.com/terraform-providers/terraform-provider-openstack/issues/531))
* Added `description` to `openstack_networking_subnet_v2` ([#533](https://github.com/terraform-providers/terraform-provider-openstack/issues/533))
* Added `description` to `openstack_networking_floatingip_v2` ([#534](https://github.com/terraform-providers/terraform-provider-openstack/issues/534))
* Added `description` to `openstack_networking_secgroup_v2` data source ([#535](https://github.com/terraform-providers/terraform-provider-openstack/issues/535))
* Added `description` to `openstack_networking_network_v2` ([#532](https://github.com/terraform-providers/terraform-provider-openstack/issues/532))
* Added `description` to `openstack_networking_subnet_v2` data source ([#528](https://github.com/terraform-providers/terraform-provider-openstack/issues/528))
* Added `description` to `openstack_networking_router_v2` data source ([#530](https://github.com/terraform-providers/terraform-provider-openstack/issues/530))
* Added `description` to `openstack_networking_network_v2` data source ([#536](https://github.com/terraform-providers/terraform-provider-openstack/issues/536))
* Added `description` to `openstack_networking_floatingip_v2` data source ([#523](https://github.com/terraform-providers/terraform-provider-openstack/issues/523))

BUG FIXES

* Allow instances to be in a state of `migrating` when performing a plan/refresh ([#496](https://github.com/terraform-providers/terraform-provider-openstack/issues/496))
* Fix issue when `openstack_networking_floatingip_v2`, `openstack_networking_router_v2`, `openstack_networking_subnet_v2`, and `openstack_networking_subnetpool_v2` tag updates send empty updates for the resource. ([#519](https://github.com/terraform-providers/terraform-provider-openstack/issues/519))

## 1.12.0 (November 13, 2018)

FEATURES

* __New Resource__: `openstack_compute_interface_attach_v2` ([#470](https://github.com/terraform-providers/terraform-provider-openstack/issues/470))

IMPROVEMENTS

* Added `tags` to `openstack_networking_network_v2` ([#454](https://github.com/terraform-providers/terraform-provider-openstack/issues/454))
* Added `tags` to `openstack_networking_subnet_v2` ([#459](https://github.com/terraform-providers/terraform-provider-openstack/issues/459))
* Added `tags` to `openstack_networking_subnetpool_v2` ([#460](https://github.com/terraform-providers/terraform-provider-openstack/issues/460))
* Added `tags` to `openstack_networking_port_v2` ([#461](https://github.com/terraform-providers/terraform-provider-openstack/issues/461))
* Added `tags` to `openstack_networking_secgroup_v2` ([#463](https://github.com/terraform-providers/terraform-provider-openstack/issues/463))
* Added `tags` to `openstack_networking_floatingip_v2` ([#466](https://github.com/terraform-providers/terraform-provider-openstack/issues/466))
* Added `tags` to `openstack_networking_router_v2` ([#467](https://github.com/terraform-providers/terraform-provider-openstack/issues/467))
* Added `extra_dhcp_options` to `openstack_networking_port_v2` ([#258](https://github.com/terraform-providers/terraform-provider-openstack/issues/258))
* Added `fingerprint` to `openstack_compute_keypair_v2` data source ([#481](https://github.com/terraform-providers/terraform-provider-openstack/issues/481))
* Added `extra_specs` to `openstack_compute_flavor_v2` data source ([#480](https://github.com/terraform-providers/terraform-provider-openstack/issues/480))

BUG FIXES

* Fixed issue with nova-network based environments having the `tenantnetworks` API disabled ([#485](https://github.com/terraform-providers/terraform-provider-openstack/issues/485))


## 1.11.0 (October 29, 2018)

FEATURES

* __New Resource__: `openstack_networking_trunk_v2` ([#446](https://github.com/terraform-providers/terraform-provider-openstack/issues/446))
* __New Resource__: `openstack_compute_flavor_access_v2` ([#447](https://github.com/terraform-providers/terraform-provider-openstack/issues/447))

IMPROVEMENTS

* Added `multiattach` argument and attribute for the `openstack_blockstorage_volume_v3` resource ([#431](https://github.com/terraform-providers/terraform-provider-openstack/issues/431))
* `openstack_dns_recordset_v2` can now accept IPv6 addresses with and without brackets ([#443](https://github.com/terraform-providers/terraform-provider-openstack/issues/443))
* Added `multiattach` argument for the `openstack_compute_volume_attach_v2` resource ([#442](https://github.com/terraform-providers/terraform-provider-openstack/issues/442))
* `openstack_lb_member_v2` resources can now use a weight of 0 ([#451](https://github.com/terraform-providers/terraform-provider-openstack/issues/451))

BUG FIXES

* Fixed an issue where environment variables were overwriting specified arguments ([#436](https://github.com/terraform-providers/terraform-provider-openstack/issues/436))
* Fixed an issue where security group rule descriptions were not working with older verisons of OpenStack ([#438](https://github.com/terraform-providers/terraform-provider-openstack/issues/438))

## 1.10.0 (October 01, 2018)

FEATURES

* __New Resource__: `openstack_containerinfra_cluster_v1` ([#421](https://github.com/terraform-providers/terraform-provider-openstack/issues/421))
* __New Data Source__: `openstack_containerinfra_clustertemplate_v1` ([#415](https://github.com/terraform-providers/terraform-provider-openstack/issues/415))

IMPROVEMENTS

* Added `description` argument for the `openstack_networking_secgroup_rule_v2` resource ([#416](https://github.com/terraform-providers/terraform-provider-openstack/issues/416))
* Added a vendor option of `ignore_resize_confirmation` to `openstack_compute_instance_v2` ([#422](https://github.com/terraform-providers/terraform-provider-openstack/issues/422))
* `openstack_compute_instance_v2` IP addresses are now visible in Rackspace. This provider still does not officially support Rackspace, though. ([#426](https://github.com/terraform-providers/terraform-provider-openstack/issues/426))
* Added `no_fixed_ip` argument to `openstack_networking_port_v2` which allows the port to not have an IP address ([#433](https://github.com/terraform-providers/terraform-provider-openstack/issues/433))

BUG FIXES

* Enabled instances to be in an `ERROR` state so they can be cleanly deleted ([#428](https://github.com/terraform-providers/terraform-provider-openstack/issues/428))

## 1.9.0 (September 05, 2018)

FEATURES

* __New Resource__: `openstack_objectstorage_tempurl_v1` ([#379](https://github.com/terraform-providers/terraform-provider-openstack/issues/379))
* __New Resource__: `openstack_containerinfra_clustertemplate_v1` ([#403](https://github.com/terraform-providers/terraform-provider-openstack/issues/403))
* __New Data Source__: `openstack_fw_policy_v1` ([#398](https://github.com/terraform-providers/terraform-provider-openstack/issues/398))
* __New Data Source__: `openstack_networking_router_v2` ([#401](https://github.com/terraform-providers/terraform-provider-openstack/issues/401))

IMPROVEMENTS

* The `openstack_images_image_v2` resource can now finally update properties. This update has been in progress over the last two release cycles. Please let us know if you encounter any problems ([#409](https://github.com/terraform-providers/terraform-provider-openstack/issues/409))

## 1.8.0 (August 08, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Support for `default_domain` has been added. This should not cause any issues, but please report any issues encountered.
* `openstack_images_image_v2.properties` has been set to `ForceNew`. If properties are modified, the image will be recreated. Previously, updates to the properties were only happening in the Terraform state and not actually reflected on the image itself.

FEATURES

* __New Data Source__: `openstack_identity_group_v3` ([#385](https://github.com/terraform-providers/terraform-provider-openstack/issues/385))
* __New Data Source__: `openstack_networking_floatingip_v2` ([#387](https://github.com/terraform-providers/terraform-provider-openstack/issues/387))

IMPROVEMENTS

* Added support for `default_domain` during authentication ([#329](https://github.com/terraform-providers/terraform-provider-openstack/issues/329))
* The upcoming OpenStack Rocky release will be automatically adding additional properties to the `openstack_images_image_v2` resource. This resource has been patched to account for this and to reconcile these server-provided properties with the user-provided properties. In addition, `openstack_images_image_v2.properties` has been set to `ForceNew` and will recreate the image when properties have been modified. Previously, any updates to the properties were only happening in the state and not actually reflected on the image itself. ([#390](https://github.com/terraform-providers/terraform-provider-openstack/issues/390))

BUG FIXES

* The addition of the `openstack_networking_network_v2.external` data source argument caused unintended behavior of results only containing external or non-external networks. This bug has been fixed and we apologize for the inconvenience ([#384](https://github.com/terraform-providers/terraform-provider-openstack/issues/384))
* The addition of the `openstack_compute_floatingip_associate_v2.wait_until_associated` argument caused the floating IP association to be recreated when updating to a later release of this provider. This was unintended and this has been resolved ([#395](https://github.com/terraform-providers/terraform-provider-openstack/issues/395))

## 1.7.0 (August 01, 2018)

FEATURES

* __New Data Source__: `openstack_identity_endpoint_v3` ([#377](https://github.com/terraform-providers/terraform-provider-openstack/issues/377))

IMPROVEMENTS

* Allow resize for stopped instances ([#348](https://github.com/terraform-providers/terraform-provider-openstack/issues/348))
* Added `power_state` to `openstack_compute_instance_v2` ([#350](https://github.com/terraform-providers/terraform-provider-openstack/issues/350))
* Added `external` to `openstack_networking_network_v2` resource ([#357](https://github.com/terraform-providers/terraform-provider-openstack/issues/357))
* Added `external` to `openstack_networking_network_v2` data source ([#358](https://github.com/terraform-providers/terraform-provider-openstack/issues/358))
* Return the default network uuid for `openstack_compute_instance_v2` ([#365](https://github.com/terraform-providers/terraform-provider-openstack/issues/365))
* Allow a specific floating IP to be specified in `openstack_networking_floatingip_v2` ([#371](https://github.com/terraform-providers/terraform-provider-openstack/issues/371))
* Allow `PROXY` protocol for `openstack_lb_pool_v2` ([#375](https://github.com/terraform-providers/terraform-provider-openstack/issues/375))

BUG FIXES

* Allow explicit values of `0` for `min_disk_gb` and `min_ram_mb` in the `openstack_images_image_v2` resource ([#351](https://github.com/terraform-providers/terraform-provider-openstack/issues/351))
* Make `peer_ep_group_id` optional in `openstack_vpnaas_site_connection` ([#353](https://github.com/terraform-providers/terraform-provider-openstack/issues/353))

## 1.6.0 (June 20, 2018)

FEATURES

* __New Resource__: `openstack_vpnaas_site_connection_v2` ([#330](https://github.com/terraform-providers/terraform-provider-openstack/issues/330))

IMPROVEMENTS

* Added `wait_until_associated` to `openstack_compute_floatingip_associate_v2` ([#310](https://github.com/terraform-providers/terraform-provider-openstack/issues/310))
* Added support for SSL settings in a `clouds.yaml` file ([#340](https://github.com/terraform-providers/terraform-provider-openstack/issues/340))

## 1.5.0 (May 15, 2018)

FEATURES

* __New Resource__: `openstack_blockstorage_volume_v3` ([#324](https://github.com/terraform-providers/terraform-provider-openstack/issues/324))
* __New Resource__: `openstack_blockstorage_volume_attach_v3` ([#324](https://github.com/terraform-providers/terraform-provider-openstack/issues/324))
* __New Resource__: `openstack_networking_subnet_route_v2` ([#314](https://github.com/terraform-providers/terraform-provider-openstack/issues/314))
* __New Resource__: `openstack_networking_floatingip_associate_v2` ([#313](https://github.com/terraform-providers/terraform-provider-openstack/issues/313))
* __New Resource__: `openstack_vpnaas_ipsec_policy_v2` ([#270](https://github.com/terraform-providers/terraform-provider-openstack/issues/270))
* __New Resource__: `openstack_vpnaas_service_v2` ([#300](https://github.com/terraform-providers/terraform-provider-openstack/issues/300))
* __New Resource__: `openstack_vpnaas_ike_policy_v2` ([#316](https://github.com/terraform-providers/terraform-provider-openstack/issues/316))
* __New Resource__: `openstack_vpnaas_endpoint_group_v2` ([#321](https://github.com/terraform-providers/terraform-provider-openstack/issues/321))
* __New Data Source__: `openstack_compute_keypair_v2` ([#307](https://github.com/terraform-providers/terraform-provider-openstack/issues/307))
* __New Data Source__: `openstack_identity_auth_scope_v3` ([#204](https://github.com/terraform-providers/terraform-provider-openstack/issues/204))

IMPROVEMENTS

* Added `verify_checksum` to `openstack_images_image_v2` resource so that checksum verification can be disabled ([#305](https://github.com/terraform-providers/terraform-provider-openstack/issues/305))
* The LBaaS v2 resources have lower "delay" times when waiting for state changes. This should speed up creation of a Load Balancing stack ([#297](https://github.com/terraform-providers/terraform-provider-openstack/issues/297))

BUG FIXES

* Fixed issue where `OS_IDENTITY_API_VERSION=2` was not recognized ([#315](https://github.com/terraform-providers/terraform-provider-openstack/issues/315))
* Fixed issue when using Identity v3 resources when an Identity v2 endpoint is published ([#320](https://github.com/terraform-providers/terraform-provider-openstack/issues/320))
* `openstack_networking_router_v2.distributed` will now pass `false` correctly ([#308](https://github.com/terraform-providers/terraform-provider-openstack/issues/308))
* `openstack_networking_router_v2.enable_snat` will now pass `false` correctly ([#309](https://github.com/terraform-providers/terraform-provider-openstack/issues/309))

## 1.4.0 (May 01, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* The OpenStack provider now has additional authentication options for `project_domain_name`, `project_domain_id`, `user_domain_name`, and `user_domain_id`. This will allow for more fine-grainted authentication scoping. This should not cause any problems with existing deployments, but please report any authentication issues after upgrading.

FEATURES

* __New Resource__: `openstack_identity_role_assignment_v3` ([#265](https://github.com/terraform-providers/terraform-provider-openstack/issues/265))
* __New Data Source__: `openstack_identity_project_v3` ([#251](https://github.com/terraform-providers/terraform-provider-openstack/issues/251))
* __New Data Source__: `openstack_identity_user_v3` ([#252](https://github.com/terraform-providers/terraform-provider-openstack/issues/252))

IMPROVEMENTS

* Added `member_status` to `openstack_images_image_v2` data source ([#269](https://github.com/terraform-providers/terraform-provider-openstack/issues/269))
* Add support for `OS_TOKEN` environment variable ([#272](https://github.com/terraform-providers/terraform-provider-openstack/issues/272))
* Added `force_destroy` to `openstack_objectstorage_container_v1` which will cause all objects in the container to be deleted when the container is deleted ([#276](https://github.com/terraform-providers/terraform-provider-openstack/issues/276))
* CIDR is now optional in `openstack_networking_subnet_v2` allowing a CIDR to be allocated from a subnet pool ([#294](https://github.com/terraform-providers/terraform-provider-openstack/issues/294))
* Added additional authentication options for domain scoping ([#290](https://github.com/terraform-providers/terraform-provider-openstack/issues/290))
* `openstack_images_image_v2` can now support OVA format ([#302](https://github.com/terraform-providers/terraform-provider-openstack/issues/302))

BUG FIXES

* `openstack_compute_instance_v2` resources can handle Availability Zones in the format of `az:host:node` ([#291](https://github.com/terraform-providers/terraform-provider-openstack/issues/291))

## 1.3.0 (March 14, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `openstack_compute_keypair_v2` can now generate a private key, however the private key will be stored in your Terraform state. Please use caution.
* The MAC addresses in `openstack_networking_port_v2.allowed_address_pairs` is no longer computed. This should not cause an issue for users since if an `allowed_address_pairs` MAC address was not specified, the AAP MAC will match `openstack_networking_port_v2.mac_address`.

FEATURES

* __New Resource:__ `openstack_networking_subnetpool_v2` ([#243](https://github.com/terraform-providers/terraform-provider-openstack/issues/243))
* __New Resource:__ `openstack_identity_role_v3` ([#250](https://github.com/terraform-providers/terraform-provider-openstack/issues/250))
* __New Data Source:__ `openstack_networking_subnetpool_v2` ([#243](https://github.com/terraform-providers/terraform-provider-openstack/issues/243))
* __New Data Source:__ `openstack_identity_role_v3` ([#250](https://github.com/terraform-providers/terraform-provider-openstack/issues/250))

IMPROVEMENTS

* Added `additional_properties` to `openstack_compute_instance_v2` scheduler hints ([#230](https://github.com/terraform-providers/terraform-provider-openstack/issues/230))
* `openstack_compute_keypair_v2` can now generate a private key ([#217](https://github.com/terraform-providers/terraform-provider-openstack/issues/217))
* `openstack_networking_router_v2` can now optionally set a default gateway after it has been created ([#209](https://github.com/terraform-providers/terraform-provider-openstack/issues/209))
* Added `subnetpool_id` to `openstack_networking_subnet_v2` resource and data source ([#249](https://github.com/terraform-providers/terraform-provider-openstack/issues/249))
* Added `extra_specs` to `openstack_compute_flavor_v2` ([#241](https://github.com/terraform-providers/terraform-provider-openstack/issues/241))
* Added `subnet_id` to `openstack_networking_floatingip_v2` ([#240](https://github.com/terraform-providers/terraform-provider-openstack/issues/240))

BUG FIXES

* Fixed bug with `openstack_networking_network_v2` and `openstack_networking_subnet_v2` where the `OS_TENANT_ID` was incorrectly being used as a default value ([#254](https://github.com/terraform-providers/terraform-provider-openstack/issues/254))
* Correctly detect if an object storage container is deleted ([#261](https://github.com/terraform-providers/terraform-provider-openstack/issues/261))
* Fixed a few small bugs with `openstack_fw_rule_v1` updating ([#224](https://github.com/terraform-providers/terraform-provider-openstack/issues/224))
* Fixed an issue with `openstack_networking_port_v2` `allowed_address_pairs` and MAC addresses ([#244](https://github.com/terraform-providers/terraform-provider-openstack/issues/244))

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
