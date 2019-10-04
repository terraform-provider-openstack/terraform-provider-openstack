# Attach, format, and mount a Block Storage Volume to an Instance

This provides a template for attaching, formating, and mounting a Block Storage Volume to an Instance.

## Usage

Download and install [Terraform](https://www.terraform.io/downloads.html):

```sh
wget -P /tmp/ https://releases.hashicorp.com/terraform/0.12.9/terraform_0.12.9_linux_amd64.zip
unzip /tmp/terraform_0.12.9_linux_amd64.zip
sudo mv terraform /usr/local/bin/
terraform --version
```

Log in to the OpenStack dashboard, choose the project for which you want to download the OpenStack RC file, and run the following commands:

```sh
source ~/Downloads/PROJECT-openrc.sh
Please enter your OpenStack Password for project PROJECT as user username:
```

### Initialize providers

```sh
terraform init
```

### Generate an execution plan

```sh
terraform plan
# or with specific variables
terraform plan -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'volume_size=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```

### Apply the plan

```sh
terraform apply -auto-approve
# or with specific variables
terraform apply -auto-approve \
                -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'volume_size=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```

Upon completion,the instances floating IP can be viewed in `instance_ip.txt`.

### Destroy

 > watch the `-force` will not ask for any confirmation.

```sh
terraform destroy -force
# or with specific variables
terraform destroy -force \
                -var 'pool=gateway' \
                -var 'flavor=m02.c02.d20' \
                -var 'volume_size=3' \
                -var 'network_name=my-network' \
                -var 'ssh_key_file=./id_rsa_os'
```
