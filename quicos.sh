#!/bin/bash

# Define the IP addresses of the VMs
r1="192.168.50.2"
r2="192.168.50.3"
c1="192.168.50.10"
s1="192.168.50.11"
vms=("$r1" "$r2", "$bond", "$sat", "$c1" "$s1")
session_name="quicos-testbed"
custom_known_hosts_file=./configs/ssh/known_hosts
# get the password from the environment variable BP
become_password=$BP

init() {
    if [ ! -d "./configs/ssh" ]; then
        mkdir -p ./configs/ssh
    fi

    vagrant up

    # For each VMs:
    #  - Add the VM to the known_hosts file
    #  - Copy the SSH key to the VM    
    for vm in "${vms[@]}"; do
        ssh-keyscan -H $vm >> $custom_known_hosts_file
        sshpass -p vagrant ssh-copy-id -f -i ~/.ssh/id_rsa.pub -o UserKnownHostsFile=$custom_known_hosts_file vagrant@$vm
    done
}

deinit() {
    vagrant destroy -f
    rm -rf ./configs/ssh
    rm -rf .vagrant 
}

tmux_create_session() {
    # If directory exists, remove it
    if [ -d "./configs/ssh" ]; then
        rm -rf ./configs/ssh
    fi

    # Create directory
    mkdir -p ./configs/ssh

    # check if the session already exists
    tmux has-session -t $session_name 2>/dev/null
    if [ $? == 0 ]; then
        echo "Session $session_name already exists. Please use 'attach' command to attach to the session."
        exit 1
    fi

    # Start a new tmux session
    tmux new-session -d -s $session_name

    # Split the window into 4 panes
    tmux split-window -h
    tmux split-window -v
    tmux select-pane -t 0
    tmux split-window -v

    # Arrange the panes evenly
    tmux select-layout tiled

    # Connect to each VM using SSH in each pane of the tmux window
    # Use the custom known_hosts file
    for i in "${!vms[@]}"; do
        tmux select-pane -t $i
        tmux send-keys "ssh -o UserKnownHostsFile=$custom_known_hosts_file vagrant@${vms[$i]}" C-m
        tmux send-keys "clear" C-m
    done
}

tmux_attach_session() {
    # Attach to the tmux session
    tmux attach-session -t $session_name
}

tmux_destroy_session() {
    # Kill the tmux session
    tmux kill-session -t $session_name
}

install_openback() {
    cd openbach/ansible
    # ansible-playbook -i ../../ansible-openbach/inventory/inventory install.yml --skip-tags check_resources -u vagrant -K --private-key ~/.ssh/id_rsa --ssh-extra-args="-o UserKnownHostsFile=../../configs/ssh/known_hosts"
    # expect << EOF
    # set timeout -1
#     spawn 
#     expect "BECOME password:"
#     send "$become_password\r"
#     expect eof
# EOF

    #    -e '{"openbach_jobs_folders": ["~/openbach_extra/external_jobs/stable_jobs/", "~/jobs_devel/"], "default_jobs": ["iperf", "fping", "socat", "my_brand_new_job", "mptcp", "squid"]}'
    ansible-playbook -i ../../ansible-openbach/inventory/inventory install.yml               \
                     -K  \
                     -e '{"openbach_jobs_folders": ["../../tutorial/jobs/"], "default_jobs": ["iperf"]}' \
                     -e project_name=quicos-test                                                  \
                     --skip-tags check_resources,configure_ntp_server -u vagrant -K          \
                     --private-key ~/.ssh/id_rsa                                             \
                     --ssh-extra-args "-o UserKnownHostsFile=../../configs/ssh/known_hosts"

    cd ../..
}

uninstall_openback() {
    cd openbach/ansible
    ansible-playbook -i ../../ansible-openbach/inventory/inventory uninstall.yml -u vagrant -K --private-key ~/.ssh/id_rsa --ssh-extra-args="-o UserKnownHostsFile=../../configs/ssh/known_hosts"
    cd ../..
}

if [ $# -eq 0 ]; then
    echo "Usage: $0 {init|deinit|tmux|tmux-attach|tmux-destroy|install}"
    exit 1
fi

case $1 in
    init)
        init
        ;;
    deinit)
        deinit
        ;;
    t | tmux)
        tmux_create_session
        ;;
    ta | tmux-attach)
        tmux_attach_session
        ;;
    td | tmux-destroy)
        tmux_destroy_session
        ;;
    i | install)
        install_openback
        ;;
    u | uninstall)
        uninstall_openback
        ;;
    *)
        echo "Usage: $0 {init|deinit|tmux|tmux-attach|tmux-destroy|install}"
        exit 1
        ;;
esac