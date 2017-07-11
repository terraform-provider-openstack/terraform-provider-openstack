provider "aws" {
  region = "us-west-2"
}

data "aws_ami" "packstack_standard" {
  most_recent = true
  owners = ["self"]
  name_regex = "^packstack-standard-ocata"
}

resource "aws_spot_instance_request" "openstack_acc_tests" {
  ami = "${data.aws_ami.packstack_standard.id}"
  spot_price = "0.0441"
  instance_type = "m3.xlarge"
  wait_for_fulfillment = true
  spot_type = "one-time"

  root_block_device {
    volume_size = 40
    delete_on_termination = true
  }
}

resource "null_resource" "rc_files" {
  provisioner "local-exec" {
    command = <<EOF
      while true ; do
        wget http://${aws_spot_instance_request.openstack_acc_tests.public_ip}/keystonerc_demo 2> /dev/null
        if [ $? = 0 ]; then
          break
        fi
        sleep 20
      done

      wget http://${aws_spot_instance_request.openstack_acc_tests.public_ip}/keystonerc_admin
    EOF
  }
}
