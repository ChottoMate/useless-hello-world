output "public_ip" {
  value       = aws_instance.web.public_ip
  description = "The public IP of the web server"
}

output "ssh_command" {
  value = "ssh -i useless-key ec2-user@${aws_instance.web.public_ip}"
  description = "Connect command"
}