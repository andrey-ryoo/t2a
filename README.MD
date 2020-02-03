# **t2a (Terrafirm To Ansible CLI)**

This app is created to convert terraform outputs from state to valid ansible inventory file

#### **How to use**

Write Terraform HCL describing your infrastructure.
Do not forget to write outputs.tf.

```hcl
output "master_addresses" {
  value = join(", ",openstack_compute_instance_v2.k8s-masters-seed.network[0].fixed_ip_v4)
}

output "minions_adresses" {
  value = join(", ",openstack_compute_instance_v2.k8s-minions[*].network[0].fixed_ip_v4)
}
```
Then run
```bash
#If you not compiled t2a do so with go get && go build
terraform apply && ./t2a
```
You will get something like

```yaml
all:
  children:
    master_addresses:
      hosts:
        10.219.150.27: 
        10.219.150.28: 
        10.219.150.30:
    minions_adresses:
      hosts:
        10.219.150.31: 
        10.219.150.32: 
        10.219.150.33:
```
Then you can do smth like:
```bash
ansible-playbook -u user -i inventory.yaml --private-key ~/.ssh/id_rsa your-playbook.yaml
```
Enjoy
