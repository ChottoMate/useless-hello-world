data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }
}

resource "aws_key_pair" "auth" {
  key_name   = "useless-key"
  public_key = file("~/useless-key.pub")
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.al2023.id
  instance_type = "t2.micro"

  subnet_id                   = aws_subnet.public_1a.id
  vpc_security_group_ids      = [aws_security_group.web.id]
  key_name                    = aws_key_pair.auth.key_name
  associate_public_ip_address = true

  user_data = <<-EOF
              #!/bin/bash
              dnf update -y
              dnf install -y docker git
              systemctl start docker
              systemctl enable docker
              usermod -a -G docker ec2-user

              # buildxのインストール
              mkdir -p /usr/local/lib/docker/cli-plugins
              curl -SL https://github.com/docker/buildx/releases/download/v0.19.0/buildx-v0.19.0.linux-amd64 -o /usr/local/lib/docker/cli-plugins/docker-buildx
              chmod +x /usr/local/lib/docker/cli-plugins/docker-buildx
              
              # Docker Compose (Plugin版) のインストール
              curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 -o /usr/local/lib/docker/cli-plugins/docker-compose
              chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
              EOF

  tags = {
    Name = "useless-ec2"
  }
}