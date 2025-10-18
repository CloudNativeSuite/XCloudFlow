ansible all -i inventory -m shell -a "du -hs /var/log/deepflow-agent/" -u root
ansible all -i inventory -m shell -a "find /var/log/deepflow-agent/ -name '*.log' -mtime +2 -delete" -u root --become
ansible-playbook -i inventory logrotate-deepflow-agent.yml -D -C
ansible-playbook -i inventory logrotate-deepflow-agent.yml -D
