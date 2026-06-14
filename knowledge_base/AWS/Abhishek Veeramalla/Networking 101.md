SG vs NACL
- SG works at **instance** level
- NACL works at **subnet** level
- NACL - Network Access Control List
- SG - Security Group
- NACL can be used to control multiple SG and override Allow rules in SG
- Order of NACL evaluation -> ascending order of Rule number
- Flow -> 
	- Inbound flow: Internet-Gateway -> NACL -> SG
	- Outbound flow: SG -> NACL -> NAT ->