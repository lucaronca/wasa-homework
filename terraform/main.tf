provider "aws" {
  region = "us-east-1"
}

module "ec2-setup" {
  source = "github.com/lucaronca/ec2-setup?ref=v0.0.2"
  ami = {
    name  = "amzn2-ami-kernel-5.10-hvm-2.0.20221210.1-x86_64-gp2"
    owner = "137112412989"
  }
  instance_tags = {
    name        = "SERVER01"
    environment = "DEV"
    os          = "AMAZON-LINUX"
  }
  allow_tls = false
}

output "instance_id" {
  description = "ID of the EC2 instance"
  value       = module.ec2-setup.instance_id
}

output "instance_public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = module.ec2-setup.instance_public_ip
}

output "instance_public_dns" {
  description = "Public IPv4 DNS address of the EC2 instance"
  value       = module.ec2-setup.instance_public_dns
}
