

```sh
 cd Downloads
➜  Downloads
➜  Downloads cat geek-ec2-kp.pem
-----BEGIN RSA PRIVATE KEY-----
<base64-url-encoded private key>
-----END RSA PRIVATE KEY-----%
➜  Downloads
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.226.167.59
^C # here there was a security-group (firewall) issue as it was not allowing any inbound. added my IP against SSH protocol and then it worked
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.226.167.59
The authenticity of host '54.226.167.59 (54.226.167.59)' cant be established.
ED25519 key fingerprint is SHA256:WPHp1jPiPTmWCriZegalziBOTc3W0464HIYAdj0XxUU. # # this is actually the SHA-256 hash of server-host-keypair's public key for human readability
This key is not known by any other names.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '54.226.167.59' (ED25519) to the list of known hosts.
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@         WARNING: UNPROTECTED PRIVATE KEY FILE!          @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
Permissions 0644 for 'geek-ec2-kp.pem' are too open.
It is required that your private key files are NOT accessible by others.
This private key will be ignored.
Load key "geek-ec2-kp.pem": bad permissions
ec2-user@54.226.167.59: Permission denied (publickey,gssapi-keyex,gssapi-with-mic).
➜  Downloads chmod 400 geek-ec2-kp.pem
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.226.167.59
   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/
[ec2-user@ip-172-31-20-136 ~]$ ls
[ec2-user@ip-172-31-20-136 ~]$ pwd
/home/ec2-user
[ec2-user@ip-172-31-20-136 ~]$ exit
logout
Connection to 54.226.167.59 closed.
➜  Downloads echo "I am restarting my ec2 instance now"
I am restarting my ec2 instance now
➜  Downloads echo "I have rebooted my ec2 instance now. let's try connecting again"
I have rebooted my ec2 instance now. lets try connecting again
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.226.167.59
   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/
Last login: Thu Apr  9 19:27:53 2026 from 49.37.115.193
[ec2-user@ip-172-31-20-136 ~]$ ls
[ec2-user@ip-172-31-20-136 ~]$ pwd
/home/ec2-user
[ec2-user@ip-172-31-20-136 ~]$ exit
logout
Connection to 54.226.167.59 closed.
➜  Downloads echo "rebooting didnt give any MITM warning, probably because its the same instance and the same public key of this geek-ec2-kp key-pair thats being used by this instance. lets try to stop and start and then try"
rebooting didnt give any MITM warning, probably because its the same instance and the same public key of this geek-ec2-kp key-pair thats being used by this instance. lets try to stop and start and then try
➜  Downloads echo "now i have stopped and restarted. now a fresh instance is being started. lets see"
now i have stopped and restarted. now a fresh instance is being started. lets see
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.226.167.59
^C
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.196.250.140
The authenticity of host '54.196.250.140 (54.196.250.140)' cant be established.
ED25519 key fingerprint is SHA256:WPHp1jPiPTmWCriZegalziBOTc3W0464HIYAdj0XxUU. # this hash is again showing up the same as it was showing previously. this means that the same sshd server is running and the server-host-keypair is not changed/generated again and it survived instance stop/start since the EBS volume (root) attached is same even after stop/start
This host key is known by the following other names/addresses:
    ~/.ssh/known_hosts:42: 54.226.167.59
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '54.196.250.140' (ED25519) to the list of known hosts.
   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/
Last login: Thu Apr  9 19:29:51 2026 from 49.37.115.193
[ec2-user@ip-172-31-20-136 ~]$ pwd
/home/ec2-user
[ec2-user@ip-172-31-20-136 ~]$ exit
logout
Connection to 54.196.250.140 closed.
➜  Downloads echo "now on restarting, the public ip was changed. so we trusted new host now"
now on restarting, the public ip was changed. so we trusted new host now
➜  Downloads ls ~/.ssh/known_hosts
/Users/govind/.ssh/known_hosts
➜  Downloads cat ~/.ssh/known_hosts | grep 54.196.250.140
54.196.250.140 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDxD6gi+FNijA1n40q0xZy5t8YyPL/CeqX2HTzmLhKnR # this is actually the base64-encoded version of the server's server-host-keypair's public key that was sent during first handshake with this host (before it landed in known_hosts)
➜  Downloads echo "now lets try to terminate the instance and launch new instance with same geek-ec2-kp keypair"
now lets try to terminate the instance and launch new instance with same geek-ec2-kp keypair
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@54.242.80.236
The authenticity of host '54.242.80.236 (54.242.80.236)' cant be established.
ED25519 key fingerprint is SHA256:jrRK6H1dssHviNzrP2zm7ucoZOzG2zPjVprL7Bz26Pw.
This key is not known by any other names.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '54.242.80.236' (ED25519) to the list of known hosts.
   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/
[ec2-user@ip-172-31-17-41 ~]$ pwd
/home/ec2-user
[ec2-user@ip-172-31-17-41 ~]$ exit
logout
Connection to 54.242.80.236 closed.
➜  Downloads echo "i am trying to reproduce MITM warning. so I am using Elastic IP to do that. I now first launch an EC2 instance and also created an Elastic IP and associated with this instance."
i am trying to reproduce MITM warning. so I am using Elastic IP to do that. I now first launch an EC2 instance and also created an Elastic IP and associated with this instance.
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@98.90.74.60
The authenticity of host '98.90.74.60 (98.90.74.60)' cant be established.
ED25519 key fingerprint is SHA256:p8l/VDBRYXPbiHLijdOjdSf8aiRWAblFJH29CWaTDLM.
This key is not known by any other names.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '98.90.74.60' (ED25519) to the list of known hosts.
   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/
[ec2-user@ip-172-31-27-167 ~]$ pwd
/home/ec2-user
[ec2-user@ip-172-31-27-167 ~]$ exit
logout
Connection to 98.90.74.60 closed.
➜  Downloads echo "now I will terminate this instance and use the same Elastic IP to attach it to another instance so that it also gets same IP and same key-pair's public key"
now I will terminate this instance and use the same Elastic IP to attach it to another instance so that it also gets same IP and same key-pair\'s public key
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@98.90.74.60
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The fingerprint for the ED25519 key sent by the remote host is
SHA256:ZeRBDXzEqoTaF4lI2cC1ihBoUST27cXg/AdjQ2cqrUA.
Please contact your system administrator.
Add correct host key in /Users/govind/.ssh/known_hosts to get rid of this message.
Offending ECDSA key in /Users/govind/.ssh/known_hosts:48
Host key for 98.90.74.60 has changed and you have requested strict checking.
Host key verification failed.
➜  Downloads echo "finally we got the MITM warning. this is because the public key sent from new instance is
different as it's base64-encoding is differnet than what is present known-host and its IP address (host name) is same which matches with known_host so SSH thinks it could be MITM who is using same IP address as the original server although not same public key since the man in the middle cannot generate the same public key at all. to get around this we can forget this IP address from the known_hosts and trust it again freshly"
finally we got the MITM warning. this is because the public key sent from new instance is same and its IP address is also same so SSH thinks it could be MITM. to get around this we can forget this IP address from the known_hosts and trust it again freshly
➜  Downloads ssh-keygen -R 98.90.74.60
ssh-keygen   ssh-keyscan
➜  Downloads ssh-keygen -R 98.90.74.60
# Host 98.90.74.60 found: line 47
# Host 98.90.74.60 found: line 48
/Users/govind/.ssh/known_hosts updated.
Original contents retained as /Users/govind/.ssh/known_hosts.old
➜  Downloads ssh -i geek-ec2-kp.pem ec2-user@98.90.74.60
The authenticity of host '98.90.74.60 (98.90.74.60)' cant be established.
ED25519 key fingerprint is SHA256:ZeRBDXzEqoTaF4lI2cC1ihBoUST27cXg/AdjQ2cqrUA.
This key is not known by any other names.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '98.90.74.60' (ED25519) to the list of known hosts.
   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/
[ec2-user@ip-172-31-22-221 ~]$ pwd
/home/ec2-user
[ec2-user@ip-172-31-22-221 ~]$ exit
logout
Connection to 98.90.74.60 closed.
➜  Downloads echo "now we were able to SSH since we forgot the host and again trusted"
now we were able to SSH since we forgot the host and again trusted
➜  Downloads echo "now let me make sure to release the Elastic IP and terminate the instance to save from cost"
now let me make sure to release the Elastic IP and terminate the instance to save from cost
```