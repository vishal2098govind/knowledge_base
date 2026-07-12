#cpu-architecture #x86_64 #arm64
## x86_64 (also called amd64)

Intel invented the original x86 architecture in 1978, named after their 8086 chip. It evolved over decades (286, 386, 486, Pentium) always staying 32-bit and backward compatible.

AMD extended it to 64-bit in 2003 and called it `amd64`. Intel adopted the same extension. The industry standardized on `x86_64` as the neutral term since it describes the architecture without implying a manufacturer.

**Key point:** Both Intel and AMD chips are x86_64. A binary compiled for x86_64 runs on either without recompilation. They speak the same instruction language.

**Why AWS uses `x86_64` over `amd64`:** AWS offers both Intel and AMD instances. Calling both `amd64` would confuse users on Intel chips. `x86_64` is manufacturer-neutral. The manufacturer is surfaced separately via the instance name suffix: no suffix = Intel, `a` suffix = AMD.

**Who makes x86_64 chips:** Only Intel and AMD, due to decades-old cross-licensing agreements. No one else is allowed to manufacture them.

---

## arm64 (also called aarch64)

ARM stands for Advanced RISC Machine. Designed in the 1980s primarily for embedded and mobile use where power efficiency mattered more than raw speed.

`arm64` means the ARM architecture running in 64-bit mode.

**ARM Holdings** (a UK company) designs the architecture and licenses it to anyone who wants to build chips based on it. They don't manufacture chips themselves.

**Who makes arm64 chips:**

- AWS → Graviton (used in `g` suffix EC2 instances)
- Apple → M1, M2, M3 (Macs), A-series (iPhones)
- Qualcomm → Snapdragon (Android phones, some laptops)
- Samsung → Exynos (some Galaxy phones)
- NVIDIA → Tegra (embedded boards)
- MediaTek → budget Android phones

---

## CISC vs RISC

|x86_64 (CISC)|arm64 (RISC)|
|---|---|---|
|Instruction set|Large, complex|Small, simple|
|Transistors needed|More|Fewer|
|Power consumption|Higher|Lower|
|Heat generated|More|Less|
|Historically strong at|Raw performance|Power efficiency|

**CISC** (Complex Instruction Set Computer): more transistors switching at high frequency = more electrical activity = more heat.

**RISC** (Reduced Instruction Set Computer): simpler instructions, fewer transistors, runs cooler, draws less power.

This is why every smartphone runs ARM, and why Apple's M1 switch in 2020 dramatically improved MacBook battery life. Same principle applies at AWS data center scale, multiplied by millions of servers, making Graviton cheaper to run and cheaper for you to use.

---

## Chip vs CPU

**Chip** is the physical thing. A piece of silicon with microscopic transistors etched onto it, packaged in a small casing. Informal term for "integrated circuit."

**CPU** is a functional role. Central Processing Unit, the component responsible for executing program instructions. Defined by what it does, not what it's made of.

A CPU is always implemented as a chip. But a chip is not always a CPU. GPUs, RAM, network cards, SSD controllers are all chips, none of them are CPUs.

**SoC (System on a Chip):** Apple's M1 is one physical chip containing a CPU, GPU, RAM, neural engine, and other components all on the same piece of silicon. One chip, multiple functional units inside it.

---

## Quick mental model

```
x86_64 / amd64  →  Intel or AMD chips, only ever those two
arm64 / aarch64 →  anyone who licensed from ARM Holdings
```

---

## AWS context

- EC2 AMI architecture filter: `x86_64` vs `arm64`
- Lambda runtime architecture: `x86_64` vs `arm64`
- Docker images: `linux/amd64` (same as x86_64) vs `linux/arm64`
- A binary compiled for x86_64 won't run natively on arm64 and vice versa
- Graviton (`g` suffix instances) = arm64, needs software compiled for arm64
- Most popular software and Docker images support both today