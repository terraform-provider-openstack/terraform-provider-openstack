# Multiple Instances with Floating IPs

## Prerequisites

We assume you already installed terraform using the following instructions: https://www.terraform.io/downloads.html

## Getting started

First, download an openrc file from the OpenStack dashboard. If your cloud provider does not support downloading the openrc file from a dashboard, then you can use the openrc.sample file in this directory. Rename it to openrc and modify it using a text editor of your choice.

After you have modified it, source it into your shell environment:

```
source ./openrc
```

Init providers

```
terraform init
```


## Deploy the Configuration

plan the deployment

```
terraform plan
# or with specific variables
terraform plan -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'instance_count=3' \
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
                -var 'instance_count=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```


## Tear down

cleanup, watch the `-force` will not ask for any confirmation.

```
terraform destroy -force
# or with specific variables
terraform destroy -force \
                -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'instance_count=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```
