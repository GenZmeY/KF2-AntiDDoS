# KF2-AntiDDoS

Compiled versions for windows and linux are available on the [releases page](https://github.com/GenZmeY/KF2-AntiDDoS/releases).  
But you can build it yourself, for this there is a Makefile.

## How it works
The program parses the output of the KF2 server(s) and counts the number of connections. If the number of connections from one IP exceeds the threshold and it is still not known that this is a player, the program will execute a deny script passing it the IP as an argument.  
The program will periodically execute the allow script, passing it a set of IPs blocked in the last period.

## HowTo:
Program usage and parameters [see here](https://github.com/GenZmeY/KF2-AntiDDoS/blob/master/doc/README)

- Prepare an IP deny script for your firewall. The script must block the IP received by the first argument 
- Prepare an IP set allow script for your firewall. The script must unblock the set of IPs given by the arguments 
- Сreate a redirection of the output of all KF2 servers to the stdin of the program 
- In the parameters specify the scripts that you prepared and the shell that will execute them 

## Centos example 
(change paths and values as you need) 
### systemd service:
```
[Unit]
Description=kf2-antiddos
After=network-online.target
Wants=network-online.target

[Service]
User=root
Group=root
Type=simple
ExecStart=/bin/sh -c '/usr/bin/kf2-srv log tail | /usr/local/bin/kf2-antiddos-linux-amd64 /bin/bash /usr/local/share/kf2-antiddos/deny.sh /usr/local/share/kf2-antiddos/allow.sh'
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
### deny.sh
```
#!/bin/bash

firewall-cmd --add-rich-rule="rule family=ipv4 source address=$1 port port=7777-7815 protocol=udp reject"
firewall-cmd --add-rich-rule="rule family=ipv4 destination address=$1 reject"
```

### allow.sh
```
#!/bin/bash

for IP in $@
do
        firewall-cmd --remove-rich-rule="rule family=ipv4 source address=$IP port port=7777-7815 protocol=udp reject"
        firewall-cmd --remove-rich-rule="rule family=ipv4 destination address=$IP reject"
done
```

## Contributing
It would be great if someone set up and tried it on windows 

