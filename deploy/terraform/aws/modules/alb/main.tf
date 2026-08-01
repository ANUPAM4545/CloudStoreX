# ALB module placeholder (Refinement 9)
# Future enhancement: Application Load Balancer with HTTPS termination and target groups for CloudStoreX EKS/EC2 services.

variable "environment" {
  type        = string
  description = "Deployment environment name"
  default     = "dev"
}
