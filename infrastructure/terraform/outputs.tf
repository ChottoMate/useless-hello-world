output "public_ip" {
  value       = one(aws_instance.web[*].public_ip)
  description = "The public IP of the web server"
}

output "ssh_command" {
  value = "ssh -i useless-key ec2-user@${coalesce(one(aws_instance.web[*].public_ip), "NONE")}"
  description = "Connect command"
}

output "ecr_repository_urls" {
  description = "URL of ECR repositories"
  value = {
    gateway = aws_ecr_repository.repos["gateway-service"].repository_url
    hello   = aws_ecr_repository.repos["hello-service"].repository_url
    world   = aws_ecr_repository.repos["world-service"].repository_url
  }
}

output "github_actions_role_arn" {
  description = "IAM Role ARN for GitHub Actions"
  value       = aws_iam_role.github_actions.arn
}