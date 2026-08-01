# Route53 module placeholder (Refinement 9)
# Future enhancement: Route53 hosted zone records and DNS routing for CloudStoreX API and Dashboard domains.

variable "environment" {
  type        = string
  description = "Deployment environment name"
  default     = "dev"
}
