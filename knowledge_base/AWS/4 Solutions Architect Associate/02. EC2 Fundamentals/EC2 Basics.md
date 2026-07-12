#ec2 #aws #solutions-architect-udemy

## Amazon EC2
• EC2 is one of the most popular of AWS’ offering
• EC2 = Elastic Compute Cloud = Infrastructure as a Service
• It mainly consists in the capability of :
	• Renting virtual machines (EC2)
	• Storing data on virtual drives (EBS)
	• Distributing load across machines (ELB)
	• Scaling the services using an auto-scaling group (ASG)
• Knowing EC2 is fundamental to understand how the Cloud works

## EC2 sizing & configuration options 
- Operating System (OS): Linux, Windows or Mac OS - AMI
- How much compute power & cores (CPU) `c`
- How much random-access memory (RAM) `r` or `u`
- How much storage space: `i` or `s`
    - Network-attached (EBS & EFS) 
    - hardware (EC2 Instance Store) 
- Network card (ENI): speed of the card, Public IP address 
- Firewall rules: security group 
- Bootstrap script (configure at first launch): EC2 User Data 
### EC2 User Data 
- It is possible to bootstrap our instances using an EC2 User data script. 
- bootstrapping means launching commands when a machine starts 
- That script is only run once at the instance **first start** 
- EC2 user data is used to automate boot tasks such as: 
	- Installing updates 
	- Installing software 
	- Downloading common files from the internet 
	- Anything you can think of 
	- The EC2 User Data Script runs with the root user

## EC2 Instance Type Families
[[ec2_instance_type_decoder.html]]

```js
{
    "a": {
        "expansion": "ARM / Graviton1",
        "official": false,
        "deprecated": true,
        "name": "ARM-based (Graviton1)",
        "desc": "AWS's first Graviton (ARM) instance. Predates the 'g' suffix convention. Succeeded by m6g, c6g, r6g etc.",
        "use": "Legacy ARM workloads; prefer m/c/r + g suffix for new deployments"
    },
    "t": {
        "expansion": "Turbo / throttleable",
        "official": false,
        "name": "Burstable performance",
        "desc": "Baseline CPU + burst credits. Cheap. Good for low-traffic apps, dev/test.",
        "use": "Web servers, small DBs, CI workers"
    },
    "m": {
        "expansion": "Medium / main",
        "official": false,
        "name": "General purpose",
        "desc": "Balanced CPU, memory, and network. The safe default.",
        "use": "App servers, gaming backends, enterprise apps"
    },
    "c": {
        "expansion": "Compute",
        "official": true,
        "name": "Compute optimised",
        "desc": "High CPU-to-memory ratio. For CPU-heavy workloads.",
        "use": "Batch processing, media encoding, ML inference, HPC"
    },
    "r": {
        "expansion": "RAM",
        "official": true,
        "name": "Memory optimised",
        "desc": "High RAM-to-CPU ratio. For in-memory workloads.",
        "use": "In-memory DBs, real-time analytics, SAP HANA"
    },
    "x": {
        "expansion": "Extreme memory",
        "official": false,
        "name": "Memory intensive",
        "desc": "Extreme RAM. SAP and big in-memory data sets.",
        "use": "SAP HANA, Apache Spark, large-scale analytics"
    },
    "z": {
        "expansion": "Zero compromise (high freq)",
        "official": false,
        "name": "High freq + memory",
        "desc": "Fastest Intel Xeon + large memory. Per-core licensing.",
        "use": "Electronic design automation, trading, gaming"
    },
    "p": {
        "expansion": "Performance (GPU)",
        "official": false,
        "name": "GPU — training/HPC",
        "desc": "NVIDIA GPUs for parallel compute and deep learning training.",
        "use": "ML training, HPC, rendering"
    },
    "g": {
        "expansion": "GPU / Graphics",
        "official": true,
        "name": "GPU — graphics/ML",
        "desc": "NVIDIA GPUs balanced for inference and graphics.",
        "use": "ML inference, video transcoding, game streaming"
    },
    "inf": {
        "expansion": "Inferentia",
        "official": true,
        "name": "AWS Inferentia",
        "desc": "Custom AWS Inferentia chips purpose-built for ML inference.",
        "use": "High-throughput, low-cost ML inference"
    },
    "trn": {
        "expansion": "Trainium",
        "official": true,
        "name": "AWS Trainium",
        "desc": "Custom AWS Trainium chips for deep learning training.",
        "use": "Large-scale model training"
    },
    "f": {
        "expansion": "FPGA",
        "official": true,
        "name": "FPGA",
        "desc": "Field-programmable gate arrays for custom hardware acceleration.",
        "use": "Genomics, financial analytics, video processing"
    },
    "vt": {
        "expansion": "Video transcoding",
        "official": true,
        "name": "Video transcoding",
        "desc": "Xilinx FPGAs tuned for real-time video transcoding.",
        "use": "Live and stored video encoding"
    },
    "d": {
        "expansion": "Dense storage",
        "official": true,
        "name": "Dense storage (HDD)",
        "desc": "High-density HDD local storage. Massive sequential throughput.",
        "use": "Distributed file systems, data warehousing, Hadoop"
    },
    "h": {
        "expansion": "High disk throughput",
        "official": false,
        "name": "High disk throughput",
        "desc": "HDD-backed with high throughput and lower cost.",
        "use": "MapReduce, distributed file systems"
    },
    "i": {
        "expansion": "I/O",
        "official": true,
        "name": "Storage optimised (NVMe)",
        "desc": "NVMe SSD local storage. Very high IOPS.",
        "use": "NoSQL DBs, in-memory caches with persistence, OLAP"
    },
    "im": {
        "expansion": "I/O + memory balanced",
        "official": false,
        "name": "Storage optimised (balanced)",
        "desc": "Balance of compute and high NVMe SSD storage.",
        "use": "High I/O applications, real-time analytics"
    },
    "is": {
        "expansion": "I/O + storage heavy",
        "official": false,
        "name": "Storage optimised (4x NVMe)",
        "desc": "Extra NVMe SSD vs im family.",
        "use": "Data-intensive workloads needing massive local storage"
    },
    "hpc": {
        "expansion": "High performance computing",
        "official": true,
        "name": "High performance computing",
        "desc": "EFA networking, high-core CPUs for tightly coupled HPC jobs.",
        "use": "CFD, weather modelling, seismic analysis"
    },
    "mac": {
        "expansion": "macOS",
        "official": true,
        "name": "macOS (Apple silicon / Intel)",
        "desc": "Bare-metal Mac mini or Mac Pro for Apple development.",
        "use": "iOS/macOS build and test pipelines"
    },
    "u": {
        "expansion": "Ultra memory",
        "official": false,
        "name": "High memory (bare metal)",
        "desc": "Terabytes of RAM for in-memory DBs.",
        "use": "SAP HANA scale-up, in-memory analytics"
    }
}
```