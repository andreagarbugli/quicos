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
- `openbach`: Contains the OpenBACH software to run the experiment
- `setup_ssh.sh`: A script to setup the SSH keys to access the virtual machines

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

### OpenBACH Installation

To install OpenBACH, we first need to add an SSH key to the virtual machines. We can do this with the following command:

```bash
./setup_ssh.sh
```

Then, we have to move to the `openbach/ansible`.

We can test if ansible can access the virtual machines with the following command:

```bash
ansible -i inventory/inventory -u vagrant -K -m ping all --private-key ~/.ssh/id_rsa --ssh-extra-args="-o UserKnownHostsFile=../../configs/ssh/known_hosts"
```

If the command returns `pong` for all the virtual machines, we can proceed with the installation of OpenBACH:

```bash
ansible-playbook -i inventory/inventory install.yml -u vagrant -K --private-key ~/.ssh/id_rsa --ssh-extra-args="-o UserKnownHostsFile=../../configs/ssh/known_hosts"
```