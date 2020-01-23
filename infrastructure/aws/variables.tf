variable "any_ip" {
  default     = "0.0.0.0/0"
  description = "Any IP Address"
  type        = "string"
}

variable "asg_max_size" {
  default     = "2"
  description = "The Maximum Number of Servers for the ASG"
  type        = "string"
}

variable "asg_min_size" {
  default     = "1"
  description = "The Minimum Number of Servers for the ASG"
  type        = "string"
}

variable "credentials_path" {
  description = "Path to Users AWS Credentials"
  type        = "string"
}

variable "http_port" {
  default     = "80"
  description = "The Communication Port for HTTP Requests"
  type        = "string"
}

variable "https_port" {
  default     = "443"
  description = "The Communication Port for Secured HTTP Requests"
  type        = "string"
}

variable "instance_type" {
  description = "AWS Istance Type"
  type        = "string"
}

variable "profile" {
  description = "The AWS Profile in Credentials Path to Use"
  type        = "string"
}

variable "project" {
  description = "The Project Name"
  type        = "string"
}

variable "region" {
  description = "AWS Region"
  type        = "string"
}

variable "ssh_port" {
  default     = "22"
  description = "The Secured Shell Communication Port"
  type        = "string"
}