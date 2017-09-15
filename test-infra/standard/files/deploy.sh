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

# Configure Swift
sudo crudini --set /etc/swift/proxy-server.conf DEFAULT bind_ip 0.0.0.0

# Install PowerDNS
sudo mysql -e "CREATE DATABASE pdns default character set utf8 default collate utf8_general_ci"
sudo mysql -e "GRANT ALL PRIVILEGES ON pdns.* TO 'pdns'@'localhost' IDENTIFIED BY 'password'"
sudo yum install -y epel-release yum-plugin-priorities
sudo curl -o /etc/yum.repos.d/powerdns-auth-40.repo https://repo.powerdns.com/repo-files/centos-auth-40.repo
sudo yum install -y pdns pdns-backend-mysql

echo "daemon=no
allow-recursion=127.0.0.1
config-dir=/etc/powerdns
daemon=yes
disable-axfr=no
guardian=yes
local-address=0.0.0.0
local-ipv6=::
local-port=53
setgid=pdns
setuid=pdns
slave=yes
socket-dir=/var/run
version-string=powerdns
out-of-zone-additional-processing=no
webserver=yes
api=yes
api-key=someapikey
launch=gmysql
gmysql-host=127.0.0.1
gmysql-user=pdns
gmysql-dbname=pdns
gmysql-password=password" | sudo tee /etc/pdns/pdns.conf

sudo mysql pdns < /home/centos/files/pdns.sql
sudo systemctl restart pdns

# Install Designate
openstack user create --domain default --password password designate
openstack role add --project services --user designate admin
openstack service create --name designate --description "DNS" dns
openstack endpoint create --region RegionOne dns public http://127.0.0.1:9001/
sudo mysql -e "CREATE DATABASE designate CHARACTER SET utf8 COLLATE utf8_general_ci"
sudo mysql -e "CREATE DATABASE designate_pool_manager"
sudo mysql -e "GRANT ALL PRIVILEGES ON designate.* TO 'designate'@'localhost' IDENTIFIED BY 'password'"
sudo mysql -e "GRANT ALL PRIVILEGES ON designate_pool_manager.* TO 'designate'@'localhost' IDENTIFIED BY 'password'"
sudo mysql -e "GRANT ALL PRIVILEGES ON designate.* TO 'designate'@'localhost' IDENTIFIED BY 'password'"

sudo yum install -y openstack-designate\*
designate_conf="/etc/designate/designate.conf"
sudo cp /home/centos/files/pools.yaml /etc/designate/
sudo crudini --set $designate_conf DEFAULT debug True
sudo crudini --set $designate_conf DEFAULT debug True
sudo crudini --set $designate_conf DEFAULT notification_driver messaging
sudo crudini --set $designate_conf service:api enabled_extensions_v2 "quotas, reports"
sudo crudini --set $designate_conf keystone_authtoken auth_uri http://127.0.0.1:5000
sudo crudini --set $designate_conf keystone_authtoken auth_url http://127.0.0.1:35357
sudo crudini --set $designate_conf keystone_authtoken username designate
sudo crudini --set $designate_conf keystone_authtoken password password
sudo crudini --set $designate_conf keystone_authtoken project_name services
sudo crudini --set $designate_conf keystone_authtoken auth_type password
sudo crudini --set $designate_conf service:worker enabled true
sudo crudini --set $designate_conf service:worker notify true
sudo crudini --set $designate_conf storage:sqlalchemy connection mysql+pymysql://designate:password@127.0.0.1/designate
sudo -u designate designate-manage database sync

sudo systemctl enable designate-central designate-api
sudo systemctl enable designate-worker designate-producer designate-mdns
sudo systemctl restart designate-central designate-api
sudo systemctl restart designate-worker designate-producer designate-mdns

sudo -u designate designate-manage pool update


# Configure Manila
manila type-create default_share_type True
sudo mkdir -p /var/lib/manila/.ssh
sudo ssh-keygen -f /var/lib/manila/.ssh/id_rsa -t rsa -N ''
sudo chown -R manila:manila /var/lib/manila/.ssh
sudo crudini --set /etc/manila/manila.conf DEFAULT default_share_type default_share_type
sudo crudini --set /etc/manila/manila.conf generic service_instance_user manila
sudo crudini --set /etc/manila/manila.conf generic service_instance_password manila
source ~/keystone_demorc
_SUBNET_ID=$(openstack subnet show private_subnet -c id -f value)
manila share-network-create --neutron-net-id $_NETWORK_ID --neutron-subnet-id $_SUBNET_ID --name manila-network
_MANILA_NET=$(manila share-network-list --columns id | grep [0-9] | awk '{print $2}')
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
sudo systemctl stop designate-central designate-api
sudo systemctl stop designate-worker designate-producer designate-mdns

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
