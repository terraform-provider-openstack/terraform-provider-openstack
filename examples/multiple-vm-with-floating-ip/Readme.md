# Openstack deploy with Terraform

## Prerequists

Install Terraform

```
LATEST_VERSION=$(curl -sL https://releases.hashicorp.com/terraform/ | grep href | grep -v "\(alpha\|beta\|rc\)" | head -n2 | tail -n1 | sed -e 's#.*">.*_\(.*\)<.*#\1#g')
wget https://releases.hashicorp.com/terraform/${LATEST_VERSION}/terraform_${LATEST_VERSION}_linux_amd64.zip
sudo unzip terraform_${LATEST_VERSION}_linux_amd64.zip -d /usr/local/bin/
```

## Getting started

Source openstack configuration

```
source ./my-project-openrc.sh
```
Note: edit `openrc.sample` file

Init providers

```
terraform init
```

Plan for `n` instances

```
terraform plan
# or with specific variables
terraform plan -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'count=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```

Apply the plan

```
terraform apply -auto-approve
# or with specific variables
terraform apply -auto-approve \
                -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'count=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```


## Tear down

cleanup

```
terraform destroy -force
# or with specific variables
terraform destroy -force \
                -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'count=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```
