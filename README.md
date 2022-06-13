# KF2-AntiDDoS

Compiled versions for windows and linux are available on the [releases page](https://github.com/GenZmeY/KF2-AntiDDoS/releases).  
But you can build it yourself, for this there is a Makefile.

## How it works
The program parses the output of the KF2 server(s) and counts the number of connections. If the number of connections from one IP exceeds the threshold and it is still not known that this is a player, the program will execute a deny script passing it the IP as an argument.  
The program will periodically execute the allow script, passing it a set of IPs blocked in the last period.

## HowTo:
```
Usage: <kf2_logs_output> | kf2-antiddos [option]... <shell> <deny_script> <allow_script>

kf2_logs_output            KF2 logs to redirect to stdin
shell                      shell to run deny_script and allow_script
deny_script                firewall deny script (takes IP as argument)
allow_script               firewall allow script (takes IPs as arguments)

Options:
  -j, --jobs N             allow N jobs at once
  -o, --output MODE        self|proxy|all|quiet
  -t, --deny-time TIME     minimum ip deny TIME (seconds)
  -c, --max-connections N  Skip N connections before run deny script
  -v, --version            Show version
  -h, --help               Show help
```

- Prepare an IP deny script for your firewall. The script must block the IP received by the first argument 
- Prepare an IP set allow script for your firewall. The script must unblock the set of IPs given by the arguments 
- Сreate a redirection of the output of all KF2 servers to the stdin of the program 
- In the parameters specify the scripts that you prepared and the shell that will execute them 

## Raw example
```
tail -f ./KFGame/Logs/Launch.log | ./kf2-antiddos-linux-amd64 /bin/bash ./deny.sh ./allow.sh
```

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

pay attention to this part:  
`/usr/bin/kf2-srv log tail`  
I use a self-written system to manage the server - the command specified here combines the output of all server logs into one stdout stream. If you want to protect several servers with antiddos, you also need to combine their logs into one stream. Replace this command with yours.

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
It would be great if someone set up, tried it on windows and share their experience 

