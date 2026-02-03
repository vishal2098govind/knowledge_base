#cryptography 

### Course Objectives:
- How to 
	- **use** cryptography correctly
	- and **reason about** **security**

### Cryptography Is Everywhere
#### Secure communication:
- **Web traffic** - is protected using HTTPS (SSL/TLS) 
- **Wireless traffic** - 
	- **Wifi traffic** - protected by WPA2
	- **Cellphone traffic** - protected using GSM
	- **Bluetooth traffic**
#### Encrypting files on disk:
- Protected by EFS, TrueCrypt
- So that even if the disk is stolen, the data/files are still secured / not compromised
#### Content Protection
#### User authentication
#### .... and much much more use cases

### Secure Communication in Web traffic - SSL/TLS
- Use HTTPS, which uses SSL/TLS under the hood
- The goal of SSL/TLS:
	- make sure that as the data travels across the network, an attacker 
		- first of all, can't eavesdrop on this data
		- second of all, can't modify (tamper) the data while it's in the network
- i.e. the goal of SSL/TLS
	- no eavesdropping
	- no tampering

#### Secure Socket Layer (SSL) / TLS
- SSL/TLS is the protocol that is used to secure the web traffic
- it has two main parts
	- **Handshake Protocol**
		- Establish a shared secret key
		- using **public-key cryptography**
		- both Alice and Bob know about the key, but attacker has no idea of that key
	- **Record Layer**
		- once we have the shared secret key, we can use this key to communicate securely by properly encrypting data between them
### Protecting Files On Disk
- Even if a disk is stolen, an attacker can't actually read the contents in the file
- If an attacker tries to modify the data on disk, the data in the file, while the data is on disk, it will be detected when Alice tries to decrypt the file, and thus could be ignored.
- So we have both confidentiality and integrity for files stored on disk
- Storing encrypted files on disk is analogous secure communication between Alice today and Alice tomorrow

### Symmetric Encryption Systems
- These are the building blocks for securing traffic
```lt
m -> message to be sent
k -> shared secret key, that only Alice and Bob know, attacker doesn't know
E -> Encryption algorithm function
D -> Decryption algorithm function
c -> cipher text

On Alice's end
E(k, m) = c
i.e. Encryption algorithm when given 'k' and 'm', returns a cipher text 'c'

On Bob's end
D(k, c) = m
i.e. Decryption algorithm when given 'k' and `c`, returns the original message 'm'
```

#### Use cases of Symmetric Encryption Systems
##### Single use key: (one time key)
- Key is only used to encrypt one message
- e.g. encrypted email: new key generated for every email
##### Multi use key: (many time key)
- Key used to encrypt multiple messages
- e.g. encrypt files in a file system - same key used to encrypt many files
- Need more machinery than for one-time key