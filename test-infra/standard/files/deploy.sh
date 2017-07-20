#!/bin/bash
set -x

sudo hostnamectl set-hostname localhost
sudo yum -y update

cd
sudo yum install -y -q nfs-utils
sudo mkdir /mnt/nfs
sudo chown nfsnobody:nfsnobody /mnt/nfs
sudo chmod 777 /mnt/nfs
echo "/mnt/nfs 127.0.0.1(rw,sync,no_root_squash,no_subtree_check)" | sudo tee /etc/exports
sudo exportfs -a

sudo yum install -y -q centos-release-openstack-ocata
sudo yum update -y -q
sudo yum install -y -q openstack-packstack crudini

# Run packstack
sudo packstack --answer-file /home/centos/files/packstack-answers.txt

# Configure LBaaSv2 and FWaaS
#sudo crudini --set /etc/neutron/neutron.conf DEFAULT debug True
#sudo crudini --set /etc/neutron/l3_agent.ini DEFAULT debug True
#sudo crudini --set /etc/neutron/neutron.conf DEFAULT service_plugins router,firewall,neutron_lbaas.services.loadbalancer.plugin.LoadBalancerPluginv2
#sudo crudini --set /etc/neutron/neutron_lbaas.conf service_providers service_provider LOADBALANCERV2:Haproxy:neutron_lbaas.drivers.haproxy.plugin_driver.HaproxyOnHostPluginDriver:default
#sudo crudini --set /etc/neutron/neutron.conf service_providers service_provider FIREWALL:Iptables:neutron.agent.linux.iptables_firewall.OVSHybridIptablesFirewallDriver:default
#sudo crudini --set /etc/neutron/lbaas_agent.ini DEFAULT interface_driver neutron.agent.linux.interface.OVSInterfaceDriver
#sudo crudini --set /etc/neutron/l3_agent.ini AGENT extensions fwaas
#sudo crudini --set /etc/neutron/neutron.conf fwaas enabled True
#sudo crudini --set /etc/neutron/neutron.conf fwaas driver iptables
#sudo crudini --set /etc/neutron/neutron.conf fwaas agent_version v1

#sudo neutron-db-manage --subproject neutron-lbaas upgrade head
#sudo neutron-db-manage --subproject neutron-fwaas upgrade head
#sudo systemctl disable neutron-lbaas-agent.service
#sudo systemctl restart neutron-server.service
#sudo systemctl restart neutron-l3-agent.service
#sudo systemctl enable neutron-lbaasv2-agent.service
#sudo systemctl start neutron-lbaasv2-agent.service

# Move findmnt to allow multiple mounts to 127.0.0.1:/mnt
sudo mv /bin/findmnt{,.orig}

# Prep the testing environment by creating the required testing resources and environment variables
sudo cp /root/keystonerc_demo /home/centos
sudo cp /root/keystonerc_admin /home/centos
sudo chown centos: /home/centos/keystonerc*
source /home/centos/keystonerc_admin
nova flavor-create m1.acctest 99 512 5 1 --ephemeral 10
nova flavor-create m1.resize 98 512 6 1 --ephemeral 10
_NETWORK_ID=$(openstack network show private -c id -f value)
_SUBNET_ID=$(openstack subnet show private_subnet -c id -f value)
_EXTGW_ID=$(openstack network show public -c id -f value)
_IMAGE_ID=$(openstack image show cirros -c id -f value)

echo "" >> /home/centos/keystonerc_admin
echo export OS_IMAGE_NAME="cirros" >> /home/centos/keystonerc_admin
echo export OS_IMAGE_ID="$_IMAGE_ID" >> /home/centos/keystonerc_admin
echo export OS_NETWORK_ID=$_NETWORK_ID >> /home/centos/keystonerc_admin
echo export OS_EXTGW_ID=$_EXTGW_ID >> /home/centos/keystonerc_admin
echo export OS_POOL_NAME="public" >> /home/centos/keystonerc_admin
echo export OS_FLAVOR_ID=99 >> /home/centos/keystonerc_admin
echo export OS_FLAVOR_ID_RESIZE=98 >> /home/centos/keystonerc_admin

echo "" >> /home/centos/keystonerc_demo
echo export OS_IMAGE_NAME="cirros" >> /home/centos/keystonerc_demo
echo export OS_IMAGE_ID="$_IMAGE_ID" >> /home/centos/keystonerc_demo
echo export OS_NETWORK_ID=$_NETWORK_ID >> /home/centos/keystonerc_demo
echo export OS_EXTGW_ID=$_EXTGW_ID >> /home/centos/keystonerc_demo
echo export OS_POOL_NAME="public" >> /home/centos/keystonerc_demo
echo export OS_FLAVOR_ID=99 >> /home/centos/keystonerc_demo
echo export OS_FLAVOR_ID_RESIZE=98 >> /home/centos/keystonerc_demo

# Configure Manila
manila type-create default_share_type True
sudo mkdir -p /var/lib/manila/.ssh
sudo ssh-keygen -f /var/lib/manila/.ssh/id_rsa -t rsa -N ''
sudo chown -R manila:manila /var/lib/manila/.ssh
sudo crudini --set /etc/manila/manila.conf DEFAULT default_share_type default_share_type
sudo crudini --set /etc/manila/manila.conf generic service_instance_user manila
sudo crudini --set /etc/manila/manila.conf generic service_instance_password manila
_MANILA_NET=$(openstack network show manila_service_network -c id -f value)
echo export OS_SHARE_NETWORK_ID=$_MANILA_NET >> /home/centos/keystonerc_admin
echo export OS_SHARE_NETWORK_ID=$_MANILA_NET >> /home/centos/keystonerc_demo

# Clean up the currently running services
sudo systemctl stop openstack-cinder-backup.service
sudo systemctl stop openstack-cinder-scheduler.service
sudo systemctl stop openstack-cinder-volume.service
sudo systemctl stop openstack-nova-cert.service
sudo systemctl stop openstack-nova-compute.service
sudo systemctl stop openstack-nova-conductor.service
sudo systemctl stop openstack-nova-consoleauth.service
sudo systemctl stop openstack-nova-novncproxy.service
sudo systemctl stop openstack-nova-scheduler.service
sudo systemctl stop neutron-dhcp-agent.service
sudo systemctl stop neutron-l3-agent.service
sudo systemctl stop neutron-lbaasv2-agent.service
sudo systemctl stop neutron-metadata-agent.service
sudo systemctl stop neutron-openvswitch-agent.service
sudo systemctl stop neutron-metering-agent.service
sudo systemctl stop openstack-manila-api.service
sudo systemctl stop openstack-manila-scheduler.service
sudo systemctl stop openstack-manila-share.service

sudo mysql -e "update services set deleted_at=now(), deleted=id" cinder
sudo mysql -e "update services set deleted_at=now(), deleted=id" nova
sudo mysql -e "update compute_nodes set deleted_at=now(), deleted=id" nova
for i in $(openstack network agent list -c ID -f value); do
  neutron agent-delete $i
done

sudo systemctl stop httpd

# Copy rc.local for post-boot configuration
sudo cp /home/centos/files/rc.local /etc
sudo chmod +x /etc/rc.local
