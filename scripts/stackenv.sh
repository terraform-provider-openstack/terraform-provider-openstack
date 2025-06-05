# Prep the testing environment by creating the required testing resources and
# environment variables. This env is for theopenlab CI jobs, you might need
# to modify this according to your setup

DEVSTACK_PATH=${DEVSTACK_PATH:-/opt/stack/new/devstack}
pushd $DEVSTACK_PATH
source openrc admin admin
openstack flavor create m1.acctest --id 99 --ram 512 --disk 10 --vcpu 1 --ephemeral 10
openstack flavor create m1.resize --id 98 --ram 512 --disk 11 --vcpu 1 --ephemeral 10
openstack flavor create m1.trove --id 97 --ram 512 --disk 10 --vcpu 1
openstack keypair create magnum

_NETWORK_ID=$(openstack network show private -c id -f value)
_SUBNET_ID=$(openstack subnet show private-subnet -c id -f value)
_EXTGW_ID=$(openstack network show public -c id -f value)
_IMAGE=$(openstack image list | grep -i cirros | head -n 1)
_IMAGE_ID=$(echo $_IMAGE | awk -F\| '{print $2}' | tr -d ' ')
_IMAGE_NAME=$(echo $_IMAGE | awk -F\| '{print $3}' | tr -d ' ')
_IMAGE_NAME=$(echo $_IMAGE | awk -F\| '{print $3}' | tr -d ' ')
_MAGNUM_IMAGE_ID=$(openstack image list --format value -c Name -c ID | grep coreos | cut -d ' ' -f 1)
if [ -z "$_MAGNUM_IMAGE_ID" ]; then
        _MAGNUM_IMAGE_ID=$(openstack image list --format value -c Name -c ID | grep -i atomic | cut -d ' ' -f 1)
fi
mysql_version=$(openstack datastore version list ${OS_DB_DATASTORE_TYPE} -f value -c Name --sort-column Name |head -n 1)

if [ -n "${OS_LB_ENVIRONMENT}" ]; then
        LB_FP_ID=`openstack loadbalancer flavorprofile create --provider amphora --flavor-data '{"loadbalancer_topology": "SINGLE"}' --name lb.acctest -f value -c id`
        openstack loadbalancer flavor create --name lb.acctest --flavorprofile $LB_FP_ID --description "Octavia flavor for acceptance tests" --enable
        echo export OS_LB_FLAVOR_NAME=lb.acctest >> openrc
fi

echo export OS_IMAGE_NAME="$_IMAGE_NAME" >> openrc
echo export OS_IMAGE_ID="$_IMAGE_ID" >> openrc
echo export OS_NETWORK_ID="$_NETWORK_ID" >> openrc
echo export OS_SUBNET_ID="$_SUBNET_ID" >> openrc
echo export OS_EXTGW_ID="$_EXTGW_ID" >> openrc
echo export OS_POOL_NAME="public" >> openrc
echo export OS_FLAVOR_ID=99 >> openrc
echo export OS_FLAVOR_ID_RESIZE=98 >> openrc
echo export OS_DOMAIN_ID=default >> openrc
echo export OS_MAGNUM_IMAGE_ID="$_MAGNUM_IMAGE_ID" >> openrc
echo export OS_MAGNUM_IMAGE="$_MAGNUM_IMAGE_ID" >> openrc
echo export OS_MAGNUM_FLAVOR=99 >> openrc
echo export OS_MAGNUM_KEYPAIR=magnum >> openrc
echo export OS_DB_DATASTORE_TYPE=mysql >> openrc
echo export OS_DB_DATASTORE_VERSION=${mysql_version} >> openrc

source openrc $1 $1
popd
