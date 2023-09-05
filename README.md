# KF2-AntiDDoS
**DDoS protection of the kf2 server from one of the attacks faced by community**  

Compiled versions for windows and linux are available on the [releases page](https://github.com/GenZmeY/KF2-AntiDDoS/releases).  
But you can build it yourself, for this there is a Makefile.

## ⚠️ Note ⚠️
### UPDATE 10.04.2023:
This tool has served well, but since its inception, the community has moved forward in protecting KF2 servers from DDoS.

I highly recommend paying attention to the solution from [baztheallmighty](https://forums.tripwireinteractive.com/index.php?members/baztheallmighty.110378/):  
https://forums.tripwireinteractive.com/index.php?threads/kf2-or-any-unreal-engine-3-server-on-redhat-centos-rocky-alma-linux-ddos-defense-with-the-help-of-firewalld.2337631/post-2355522  
This method limits the number of connections from each IP so junk traffic is dropped before it even reaches the kf2 server. **It is much more efficient than this tool.**  
***
If you want to continue using this tool for any reason, it will be useful to reduce the `ConnectionTimeout` so that fake connections are closed faster and do not overload the server:  
**PCServer-KFEngine.ini / LinuxServer-KFEngine.ini**  
```ini
[IpDrv.TcpNetDriver]
...
ConnectionTimeout=20.0
```
Thanks to [o2xVc3UuXp0NyBihrUnu](https://forums.tripwireinteractive.com/index.php?members/o2xvc3uuxp0nybihrunu.95080/) for [finding and sharing this setting](https://forums.tripwireinteractive.com/index.php?threads/kf2-or-any-unreal-engine-3-server-on-redhat-centos-rocky-alma-linux-ddos-defense-with-the-help-of-firewalld.2337631/page-5#post-2355506).
***
The main discussion of the DDoS issue is here:  
[forums.tripwireinteractive.com/KF2 Sever DDos Defence](https://forums.tripwireinteractive.com/index.php?threads/kf2-or-any-unreal-engine-3-server-on-redhat-centos-rocky-alma-linux-ddos-defense-with-the-help-of-firewalld.2337631/)  
You might find it helpful to follow this thread.  
### UPDATE 06.09.2023:
[o2xVc3UuXp0NyBihrUnu](https://forums.tripwireinteractive.com/index.php?members/o2xvc3uuxp0nybihrunu.95080/) adapted the [baztheallmighty](https://forums.tripwireinteractive.com/index.php?members/baztheallmighty.110378/) idea for firewall-cmd, which is quite handy:
```
firewall-cmd --permanent --direct --add-rule ipv4 filter INPUT 0 -p udp --dport 7777:7797 -m connlimit --connlimit-above 5 --connlimit-mask 20 -j DROP
```
**Source:** https://forums.tripwireinteractive.com/index.php?threads/kf2-or-any-unreal-engine-3-server-on-redhat-centos-rocky-alma-linux-ddos-defense-with-the-help-of-firewalld.2337631/post-2358698

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
I use a self-written system to manage the kf2 servers - the command specified here combines the output of all kf2 server logs into one stdout stream. If you want to protect several servers with antiddos, you also need to combine their logs into one stream. Replace this command with yours.

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

