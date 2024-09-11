
resource "aws_iam_instance_profile" "ssm_instance_profile" {
  name = "ssm_instance_profile"
  role = aws_iam_role.ssm_role.name
}


resource "aws_instance" "compute" {
  ami           = "ami-02bb4eb2718dae347" # Ubuntu arm64 24.04
  instance_type = "c7g.2xlarge"           # 8 vCPUs

  subnet_id                   = aws_subnet.public.id
  vpc_security_group_ids      = [aws_security_group.allow_ssh_p2p.id]
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ssm_instance_profile.name

  root_block_device {
    volume_size = 200
    volume_type = "gp3"
  }

  tags = {
    Name = "juno-dev"
  }

  user_data = <<-EOF
      #!/bin/bash
      set -euxo pipefail
      mkdir -p /home/ubuntu/.ssh
      for name in "derrix060" "wojciechos" "rianhughes" "pnowosie" "weiihann" "kirugan" "AnkushinDaniil" "IronGauntlets" "thiagodeev"; do
        curl -s https://github.com/$name.keys >> /home/ubuntu/.ssh/authorized_keys
      done
      chown ubuntu:ubuntu /home/ubuntu/.ssh/authorized_keys
      chmod 600 /home/ubuntu/.ssh/authorized_keys

      apt-get update

      git clone https://github.com/NethermindEth/juno.git
      wget https://juno-snapshots.nethermind.dev/files/sepolia/latest -O latest.tar
      tar -xvf latest.tar
      rm -f latest.tar
    EOF
}

resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.compute.id
  allocation_id = aws_eip.compute_eip.id
}
