### Set a Provider ###
provider "aws" {
  region                  = "${var.region}"
  shared_credentials_file = "${var.credentials_path}"
  profile                 = "${var.profile}"
}

### Get Data from other Resources ###
# Get Availability Zones
data "aws_availability_zones" "available" {}

# Get Amazon Machine Image
data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/*ubuntu-xenial-16.04-amd64-server-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  owners = ["099720109477"] # Canonical
}

# Get Domain SSL ARN
data "aws_acm_certificate" "knoxknot" {
  domain   = "knoxknot.website"
  statuses = ["ISSUED"]
}

# Get Local Data
data "template_file" "user_data" {
  template = "${file("bootstrap-server.sh")}"
}

### Upload the Server Key ###
resource "aws_key_pair" "sd_automation_key" {
  key_name   = "sd-automation-key"
  public_key = "${file("~/.ssh/csproject_key.pub")}"
}

### Create Security Groups ###
# Create the Load Balancers Security Group
resource "aws_security_group" "knoxknot_lb_sg" {
  name        = "${var.project}-lbsg"
  description = "Knox Knot Load Balancers SG"

  ingress {
    from_port   = "${var.http_port}"
    to_port     = "${var.http_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.any_ip}"]
  }

  ingress {
    from_port   = "${var.https_port}"
    to_port     = "${var.https_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.any_ip}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.any_ip}"]
  }

  tags {
    Name        = "${var.project}-lbsg"
    Project     = "${var.project}"
  }
}

# Create the WebServer Security Group
resource "aws_security_group" "knoxknot_server_sg" {
  name        = "${var.project}-wssg"
  description = "Knox Knot WebServer SG"

  ingress {
    from_port       = "${var.http_port}"
    to_port         = "${var.http_port}"
    protocol        = "tcp"
    security_groups = ["${aws_security_group.knoxknot_lb_sg.id}"]
  }

  ingress {
    from_port       = "${var.https_port}"
    to_port         = "${var.https_port}"
    protocol        = "tcp"
    security_groups = ["${aws_security_group.knoxknot_lb_sg.id}"]
  }

  ingress {
    from_port   = "${var.ssh_port}"
    to_port     = "${var.ssh_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.any_ip}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.any_ip}"]
  }

  tags {
    Name        = "${var.project}-wssg"
    Project     = "${var.project}"
  }
}

# ### Create a Load Balancer and Associate its SG ###
# Create a Load Balancer
resource "aws_elb" "knoxknot_lb" {
  name               = "${var.project}-lb"
  security_groups    = ["${aws_security_group.knoxknot_lb_sg.id}"]
  availability_zones = ["${slice(data.aws_availability_zones.available.names,0,2)}"]

  listener {
    lb_port           = "${var.http_port}"
    lb_protocol       = "http"
    instance_port     = "${var.http_port}"
    instance_protocol = "http"
  }

  listener {
    lb_port            = "${var.https_port}"
    lb_protocol        = "https"
    instance_port      = "${var.http_port}"
    instance_protocol  = "http"
    ssl_certificate_id = "${data.aws_acm_certificate.knoxknot.arn}"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    interval            = 30
    target              = "HTTP:${var.http_port}/"
  }

  tags {
    Name        = "${var.project}-lb"
    Project     = "${var.project}"
  }
}

### Create the Server Clusters ###
# Create a Launch Configuration for the Web Servers
resource "aws_launch_configuration" "sd_automation_server_config" {
  image_id        = "${data.aws_ami.ubuntu.id}"
  instance_type   = "${var.instance_type}"
  key_name        = "${aws_key_pair.sd_automation_key.key_name}"
  security_groups = ["${aws_security_group.knoxknot_server_sg.name}"]
  user_data       = "${data.template_file.user_data.rendered}"

  lifecycle {
    create_before_destroy = true
  }
}

# Create the Web Server Clusters from the Launch Configuration
resource "aws_autoscaling_group" "sd_automation_asg" {
  name                 = "${var.project}-wsasg"
  launch_configuration = "${aws_launch_configuration.sd_automation_server_config.id}"
  load_balancers = ["${aws_elb.knoxknot_lb.name}"]
  availability_zones = ["${slice(data.aws_availability_zones.available.names,0,1)}"]
  health_check_type    = "ELB"
  min_size             = "${var.asg_min_size}"
  max_size             = "${var.asg_max_size}"

  tag {
    key                 = "Name"
    value               = "${var.project}-wsasg"
    propagate_at_launch = true
  }

  tag {
    key                 = "Project"
    value               = "${var.project}"
    propagate_at_launch = true
  }
}

# Define the Web Cluster Scaling Policies
resource "aws_autoscaling_policy" "sd_automation_asg_scale_policy" {
  name                      = "${var.project}-asp"
  adjustment_type           = "ChangeInCapacity"
  policy_type               = "TargetTrackingScaling"
  estimated_instance_warmup = 300
  autoscaling_group_name    = "${aws_autoscaling_group.sd_automation_asg.name}"

  target_tracking_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ASGAverageCPUUtilization"
    }

    target_value = 70.0
  }
}

resource "aws_cloudwatch_metric_alarm" "high_cpu_usage" {
  alarm_name          = "${var.project}-high-cpu-usage"
  namespace           = "AWS/EC2"
  metric_name         = "CPUUtilization"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  period              = 300
  statistic           = "Average"
  threshold           = 80
  unit                = "Percent"
  alarm_actions       = ["${aws_autoscaling_policy.sd_automation_asg_scale_policy.arn}"]

  dimensions = {
    AutoScalingGroupName = "${aws_autoscaling_group.sd_automation_asg.name}"
  }
}