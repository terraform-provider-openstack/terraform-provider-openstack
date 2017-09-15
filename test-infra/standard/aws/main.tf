provider "aws" {
  region = "us-west-2"
}

data "aws_ami" "packstack_standard" {
  most_recent = true
  owners = ["self"]
  name_regex = "^packstack-standard-ocata"
}

resource "random_id" "security_group_name" {
  prefix = "openstack_test_instance_allow_all_"
  byte_length = 8
}

resource "aws_spot_instance_request" "openstack_acc_tests" {
  ami = "${data.aws_ami.packstack_standard.id}"
  spot_price = "0.0500"
  instance_type = "m3.xlarge"
  wait_for_fulfillment = true
  spot_type = "one-time"

  security_groups = ["${aws_security_group.allow_all.name}"]

  root_block_device {
    volume_size = 40
    delete_on_termination = true
  }

  tags {
    Name = "OpenStack Acceptance Test Infra"
  }
}

resource "aws_security_group" "allow_all" {
  name        = "${random_id.security_group_name.hex}"
  description = "OpenStack Test Infra Allow all inbound/outbound traffic"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    ipv6_cidr_blocks = ["::/0"]
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
