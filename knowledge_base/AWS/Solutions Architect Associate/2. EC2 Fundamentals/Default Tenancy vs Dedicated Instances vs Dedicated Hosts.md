#aws-ec2 #ec2-instance-tenancy


Tags: #aws #ec2 #compute #saa

---

## The Core Layered Model

```
Physical Server (Host)
└── Hypervisor (AWS Nitro)
    ├── EC2 Instance (VM 1)
    ├── EC2 Instance (VM 2)
    └── EC2 Instance (VM 3)
```

- **Host** = actual physical machine in an AWS data centre (real CPU, RAM, network card)
- **Instance** = a virtual machine carved out of that physical host by the hypervisor
- By default (shared tenancy), multiple AWS customers' VMs run on the same physical host

---

## Default Tenancy

AWS picks any available host in the AZ, carves a slice for your VM, and other customers may land on the same machine. You have zero visibility into the host.

---

## Dedicated Instances

AWS guarantees the physical host it places your instance on will **only ever run instances from your AWS account**. No other customer's VMs can share that host.

```
Physical Host (AWS chooses, you don't control which)
├── Your Dedicated Instance        ← what you launched
├── Your regular EC2 instance      ← also possible, same account
└── [never] another customer's instance
```

**Key rule:** When AWS needs to place your Dedicated Instance, it can only pick a host that is either:

- **Empty** — no instances running yet
- **Already claimed by your account** — a host that only has your instances on it

Once a host gets claimed for your account, other customers can never use it.

**Limitations:**

- No visibility into which physical host you land on
- No host affinity — stop/start may move the instance to a different (but still account-dedicated) host
- Limited BYOL support

---

## Dedicated Hosts

You provision a **specific, named physical machine**. You can see its socket count, core count, and host ID in the console.

**Key additions over Dedicated Instances:**

- Full visibility into the physical host
- Host affinity — instance comes back to the same host after stop/start
- Full BYOL support (per-socket, per-core, per-VM licenses: Windows Server, SQL Server, SUSE, RHEL, etc.)
- Integrated with AWS License Manager

---

## Comparison Table

|Default|Dedicated Instance|Dedicated Host|
|---|---|---|---|
|Isolated from other accounts|❌|✅|✅|
|Isolated from own account's other instances|❌|❌|✅|
|Control over which physical host|❌|❌|✅|
|Host affinity (same server after stop/start)|❌|❌|✅|
|BYOL licensing|❌|Limited|Full|
|Billing|Per instance|Per instance + $2/hr regional flat fee|Per physical host|

---

## Pricing

- **Dedicated Instances:** $2/hr flat fee per region (regardless of instance count) + ~10% premium per instance over on-demand
- **Dedicated Hosts:** Hourly per physical host based on instance family. Can be cheaper if host is packed efficiently.

---

## Decision Rule

```
Need hardware isolation?
├── No  → Default tenancy
└── Yes → Need BYOL licensing?
          ├── Yes → Dedicated Hosts
          └── No  → Dedicated Instances (simpler, sufficient)
```

---

## Compliance Angle

Some industries require that workloads **cannot share physical hardware with unknown third parties** — HIPAA (healthcare), PCI-DSS (payments), some government contracts.

The concern is physical-level attack vectors like side-channel attacks, where a malicious VM on the same host tries to sniff memory or CPU cache from your VM.

- **Dedicated Instances** satisfy most compliance frameworks — you can tell auditors "no other AWS customer's code runs on our hardware"
- **Dedicated Hosts** go further — you can name the exact host ID, show socket/core counts, demonstrate full documented control. Useful when auditors want hard evidence of physical isolation

---

## Licensing Angle (BYOL)

Some vendors (Microsoft, Oracle, IBM) sell licenses tied to **physical hardware**, not VMs — e.g. "per physical socket" or "per physical core."

**Problem with shared tenancy / Dedicated Instances:** you don't know the physical host's socket/core count, so you can't satisfy a vendor license audit. You end up paying AWS's included license rate.

**Dedicated Hosts solve this:** you can see the exact socket and core count of your host and present that to the vendor (e.g. Microsoft) during an audit. This is **BYOL — Bring Your Own License** — reusing a license you already paid for on-prem instead of paying again on AWS.

**Practical impact:** migrating an on-prem Windows Server or SQL Server workload with existing licenses → use Dedicated Hosts to avoid paying AWS's per-instance Windows/SQL licensing fee on top of compute cost.

---

## Exam Anchors (SAA-C03)

- "existing per-core / per-socket software licenses" → **Dedicated Hosts**
- "same physical server after restart" → **Dedicated Hosts** (host affinity)
- "compliance, no other customer on my hardware" → **Dedicated Instances** is sufficient
- Cannot go from dedicated tenancy back to default shared tenancy