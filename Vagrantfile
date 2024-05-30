# -*- mode: ruby -*-
# vi: set ft=ruby :

IMAGE_NAME = "generic/ubuntu2004"
USER = "vagrant"
PASSWORD = "vagrant"
SSH_KEY="~/.ssh/id_rsa.pub"

Vagrant.configure("2") do |config|
    config.vm.box = IMAGE_NAME    

    ["r1", "r2"].each do |name|
        config.vm.define name do |node|
            node.vm.hostname = name
            # node.ssh.username = USER
            # node.ssh.private_key_path = SSH_KEY
            node.vm.provider "libvirt" do |v|
                v.memory = 2048
                v.cpus = 1
                v.storage :file, size: "20G"
            end
            node.vm.network "private_network", ip: "192.168.50.#{name == "r1" ? 2 : 3}"
            node.vm.network "private_network", ip: "10.0.1.2", network: "10.0.1.0/24", dhcp_enabled: false, auto_config: false
            node.vm.network "private_network", ip: "10.0.1.2", network: "10.0.1.0/24", dhcp_enabled: false, auto_config: false
            node.vm.provision "ansible" do |ansible|
                ansible.playbook = "ansible/playbook.yml"
                ansible.compatibility_mode = "2.0"
                ansible.extra_vars = {
                    user: USER,
                }
            end
        end
    end
        
    ["c1", "s1"].each do |name|
        config.vm.define name do |node|
            node.vm.hostname = name
            # node.ssh.username = USER
            # node.ssh.private_key_path = SSH_KEY
            node.vm.provider "libvirt" do |v|
                # if name == "c1" 1024 else 4096
                if name == "s1"
                    v.memory = 2048
                    v.cpus = 1
                    v.storage :file, size: "15G"
                else
                    v.memory = 4096
                    v.cpus = 4
                    v.storage :file, size: "40G"
                end
            end
            node.vm.network "private_network", ip: "192.168.50.#{name == "c1" ? 10 : 11}"
            node.vm.network "private_network", ip: "10.0.1.2", network: "10.0.1.0/24", dhcp_enabled: false, auto_config: false
            node.vm.provision "ansible" do |ansible|
                ansible.playbook = "ansible/playbook.yml"
                ansible.compatibility_mode = "2.0"
                ansible.extra_vars = {
                    user: USER,
                }
            end
        end
    end
end