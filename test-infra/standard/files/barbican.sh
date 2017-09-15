openstack user create --domain default --password password barbican
openstack role add --project services --user barbican admin
openstack role create creator
openstack role add --project services --user barbican creator
openstack service create --name designate --description "Key Manager" key-manager
openstack endpoint create --region RegionOne key-manager public http://127.0.0.1:9311/
openstack endpoint create --region RegionOne key-manager admin http://127.0.0.1:9311/
openstack endpoint create --region RegionOne key-manager internal http://127.0.0.1:9311/
sudo mysql -e "CREATE DATABASE barbican CHARACTER SET utf8 COLLATE utf8_general_ci"
sudo mysql -e "GRANT ALL PRIVILEGES ON barbican.* TO 'barbican'@'localhost' IDENTIFIED BY 'password'"

sudo yum install -y openstack-barbican-api
sudo yum install -y openstack-barbican-worker

barbican_conf="/etc/barbican/barbican.conf"
sudo crudini --set $barbican_conf DEFAULT sql_connection mysql+pymysql://barbican:password@127.0.0.1/barbican
sudo crudini --set $barbican_conf DEFAULT transport_url rabbit://guest:guest@localhost/
sudo crudini --set $barbican_conf DEFAULT bind_host ::
sudo crudini --set $barbican_conf secretstore enabled_secretstore_plugins store_crypto
sudo crudini --set $barbican_conf crypto enabled_crypto_plugins simple_crypto
sudo crudini --set $barbican_conf simple_crypto_plugin kek YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY=
sudo crudini --set $barbican_conf keystone_authtoken auth_uri http://127.0.0.1:5000
sudo crudini --set $barbican_conf keystone_authtoken auth_url http://127.0.0.1:35357
sudo crudini --set $barbican_conf keystone_authtoken memcached_servers localhost:11211
sudo crudini --set $barbican_conf keystone_authtoken auth_type password
sudo crudini --set $barbican_conf keystone_authtoken project_domain_name default
sudo crudini --set $barbican_conf keystone_authtoken user_domain_name default
sudo crudini --set $barbican_conf keystone_authtoken project_name services
sudo crudini --set $barbican_conf keystone_authtoken username barbican
sudo crudini --set $barbican_conf keystone_authtoken password password
sudo -u barbican barbican-manage db upgrade

systemctl enable openstack-barbican-api.service
systemctl start openstack-barbican-api

sudo iptables -I INPUT -p tcp --dport 9311 -j ACCEPT
sudo ip6tables -I INPUT -p tcp --dport 9311 -j ACCEPT
