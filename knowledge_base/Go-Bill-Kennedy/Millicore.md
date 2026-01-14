#go #k8s #capacity-planning

## what is millicore unit or 'm' while expressing cpu limit

In Kubernetes (and most cloud/container contexts), **`millicore` (written as `m`)** is a **unit of CPU** that lets you express **fractions of a CPU core** precisely.

---

## The core idea (literally)

- **1 CPU core = 1000 millicores**
    
- **1 millicore = 0.001 CPU core**
    

So:

|CPU value|Meaning|
|---|---|
|`1000m`|1 full CPU core|
|`500m`|Half a CPU core|
|`250m`|Quarter CPU|
|`100m`|10% of a CPU|
|`50m`|5% of a CPU|

If you see `cpu: "500m"`, it simply means:  
👉 _this container can use up to half a CPU core._

---

## Why millicores exist (why not just decimals?)

Millicores solve **three practical problems**:

### 1️⃣ Precise scheduling

Containers often need **small, predictable CPU slices**.

Instead of:

`cpu: 0.1`

Kubernetes uses:

`cpu: 100m`

This avoids floating-point ambiguity and makes scheduling math deterministic.

---

### 2️⃣ Fair CPU sharing

On a node with 1 core:

- Pod A → `500m`
    
- Pod B → `500m`
    

Each gets **50% CPU time** when both are busy.

If Pod A uses only `200m`, Pod B can temporarily use more—**CPU is shared, not reserved unless limited**.

---

### 3️⃣ Hardware-agnostic abstraction

Millicores work the same whether your node has:

- 1 core
    
- 4 cores
    
- 64 cores
    

Example on a **4-core node**:

`cpu: "1000m"`

This means **1 full core**, not the entire machine.

---

## Requests vs Limits (important distinction)

In Kubernetes, you usually see CPU defined in two places:

`resources:   requests:     cpu: "250m"   limits:     cpu: "500m"`

### 🔹 `requests.cpu`

- Minimum CPU guaranteed
    
- Used by the scheduler to place pods
    

### 🔹 `limits.cpu`

- Maximum CPU allowed
    
- Enforced by Linux CFS (Completely Fair Scheduler)
    

If a container hits its CPU limit:  
👉 it gets **throttled**, not killed.

---

## What actually enforces this?

Under the hood:

- Kubernetes converts `m` → **CPU shares / quotas**
    
- Linux **CFS scheduler** enforces fair usage
    
- Time slicing is done in **CPU time**, not wall time
    

This is why CPU limits cause **throttling**, not OOM kills (memory limits do that).