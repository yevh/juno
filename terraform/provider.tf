terraform {
  required_version = "~> 1.5"
  backend "s3" {
    bucket         = "nethermind-nubia-shared-dev-tf-state"
    key            = "juno/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "nethermind-nubia-shared-dev-tf-lock"
    encrypt        = true
  }


  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.57"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
  assume_role {
    role_arn = "arn:aws:iam::891377278198:role/github-administrator-access"
  }
  default_tags {
    tags = {
      Project     = "juno"
      Group       = "Nubia"
      ManagedBy   = "Terraform"
      ProjectCode = "SWI-JUN-01"
    }
  }
}
