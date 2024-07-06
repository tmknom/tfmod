module "foo" {
  source = "../../module/foo"
}

module "bar" {
  source = "../../module/bar"
}

module "security-group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.1.2"
}
