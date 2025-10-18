#!/bin/sh
# Example commands demonstrating builtin modules
xconfig remote all -i ../../inventory -m shell -a 'hostname'
xconfig remote all -i ../../inventory -m command -a '/bin/echo hi'
xconfig remote all -i ../../inventory -m copy -a '../../motd.tmpl:/tmp/motd-copy'
