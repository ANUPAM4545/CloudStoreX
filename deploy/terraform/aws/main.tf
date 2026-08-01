terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source               = "./modules/vpc"
  environment          = var.environment
  cidr_block           = var.vpc_cidr
  availability_zones   = var.availability_zones
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
}

module "security_groups" {
  source      = "./modules/security-groups"
  environment = var.environment
  vpc_id      = module.vpc.vpc_id
  vpc_cidr    = var.vpc_cidr
}

module "rds" {
  source            = "./modules/rds"
  environment       = var.environment
  subnet_ids        = module.vpc.private_subnet_ids
  security_group_id = module.security_groups.db_security_group_id
  instance_class    = var.db_instance_class
  allocated_storage = var.db_allocated_storage
  db_password       = var.db_password
}

module "redis" {
  source            = "./modules/redis"
  environment       = var.environment
  subnet_ids        = module.vpc.private_subnet_ids
  security_group_id = module.security_groups.db_security_group_id
  node_type         = var.redis_node_type
}

module "s3" {
  source        = "./modules/s3"
  environment   = var.environment
  bucket_suffix = var.s3_bucket_suffix
  force_destroy = var.environment != "prod"
}

module "iam" {
  source        = "./modules/iam"
  environment   = var.environment
  s3_bucket_arn = module.s3.bucket_arn
}

module "ecr" {
  source      = "./modules/ecr"
  environment = var.environment
}

module "cloudwatch" {
  source            = "./modules/cloudwatch"
  environment       = var.environment
  retention_in_days = var.log_retention_days
}

# Placeholder modules (Refinement 9)
module "kms" {
  source      = "./modules/kms"
  environment = var.environment
}

module "route53" {
  source      = "./modules/route53"
  environment = var.environment
}

module "acm" {
  source      = "./modules/acm"
  environment = var.environment
}

module "alb" {
  source      = "./modules/alb"
  environment = var.environment
}
