#docker

## What is orchestration
- Scenario - 1
	- Let's say we have 5 machines (hypervisors)
	- Each server has :
		- 16 GB RAM, 8 cores, 1 TB disk
	- Each week, your team requests:
		- one VM with X RAM, Y CPU, Z disk
	- Scheduling = deciding which machine to use for each VM
	- Difficulty: easy!
- Scenario - 2
	- 1000+ machines and counting, in multiple data centers
	- each server has different resources
	- multiple times a day, a different team asks for:
		- up to 50 VMs with different characteristics
	- Scheduling = deciding which machine to use for each VM
	- Difficulty: ???
- Scenario - 3
	- you have machines (physical and/or virtual)
	- you have containers
	- you are trying to put containers on machines