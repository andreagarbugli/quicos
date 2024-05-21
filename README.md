# QUICOS

This is a project to create a virtual testbed for QUIC protocol. 

## Requirements

To run this project you need to have installed the following software:

- The Operating System must be Linux, the project was tested on Ubuntu 22.04
- `libvirt` and `qemu-kvm` (if you are using Linux)
- [Vagrant](https://www.vagrantup.com/) and the `vargrant-libvirt` plugin
- [Ansible](https://www.ansible.com/)

## Project Structure

The project is divided into the following directories and files:

- `Vagrantfile`: Contains the definition of the testbed, including the virtual machines and the network configuration
- `ansible`: Contains the Ansible playbooks to provision the virtual machines
- `ansible-openbach`: Contains the Ansible inventory for OpenBACH
- `openbach`: Contains the OpenBACH software to run the experiment
- `quicos.sh`: Script with some useful commands to run the testbed

## Setup the testbed

### Vagrant provisioning

To start the virtual testbed, you need to run the following command:

```bash
vagrant up
```

This command will create the virtual machines and provision them with the necessary software. If the provisioning fails, you can run the following command to re-run the provisioning:

```bash
vagrant provision
```

We can check the status of the virtual machines with the following command:

```bash
vagrant status
```

And we can access the virtual machines with the following command:

```bash
vagrant ssh <machine>
```

Where `<machine>` is the name of the virtual machine. The names of the virtual machines are defined in the `Vagrantfile`.

Alternatively, we can use the `quicos.sh` script to run the commands. For example, to start the virtual machines, we can run the following command:

```bash
./quicos.sh init
```

### OpenBACH Installation

We can test if ansible can access the virtual machines with the following command:

```bash
ansible -i ansible-openbach/inventory/inventory -u vagrant -K -m ping all \ 
    --private-key ~/.ssh/id_rsa \
    --ssh-extra-args="-o UserKnownHostsFile=./configs/ssh/known_hosts"
```

If the command returns `pong` for all the virtual machines, we can proceed with the installation of OpenBACH. 

To install OpenBACH, we can have two options:

1. Do `cd openbach/ansible` and run the following command:

```bash
ansible-playbook -i ../../ansible-openbach/inventory/inventory install.yml \
    -u vagrant -K --skip-tags check_resources \
    --private-key ~/.ssh/id_rsa \
    --ssh-extra-args="-o UserKnownHostsFile=../../configs/ssh/known_hosts"
```

2. Use the `quicos.sh` script:

```bash
./quicos.sh install
```