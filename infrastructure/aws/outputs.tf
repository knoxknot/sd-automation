output "asg_name" {
  value = "${aws_autoscaling_group.sd_automation_asg.name}"
}

output "lb_dns_name" {
  value = "${aws_elb.knoxknot_lb.dns_name}"
}

