output "vm_public_ip" {
  value       = aws_instance.compute.public_ip
  description = "Public IP of the VM"
}
