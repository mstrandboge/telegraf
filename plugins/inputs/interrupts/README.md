# Interrupts Input Plugin

This plugin gathers metrics about IRQs from interrupts (`/proc/interrupts`) and
soft-interrupts (`/proc/softirqs`).

⭐ Telegraf v1.3.0
🏷️ system
💻 all

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# This plugin gathers interrupts data from /proc/interrupts and /proc/softirqs.
[[inputs.interrupts]]
  ## When set to true, cpu metrics are tagged with the cpu.  Otherwise cpu is
  ## stored as a field.
  ##
  ## The default is false for backwards compatibility, and will be changed to
  ## true in a future version.  It is recommended to set to true on new
  ## deployments.
  # cpu_as_tag = false

  ## To filter which IRQs to collect, make use of tagpass / tagdrop, i.e.
  # [inputs.interrupts.tagdrop]
  #   irq = [ "NET_RX", "TASKLET" ]
```

## Metrics

There are two styles depending on the value of `cpu_as_tag`.

With `cpu_as_tag = false`:

- interrupts
  - tags:
    - irq (IRQ name)
    - type
    - device (name of the device that is located at the IRQ)
    - cpu
  - fields:
    - cpu (int, number of interrupts per cpu)
    - total (int, total number of interrupts)

- soft_interrupts
  - tags:
    - irq (IRQ name)
    - type
    - device (name of the device that is located at the IRQ)
    - cpu
  - fields:
    - cpu (int, number of interrupts per cpu)
    - total (int, total number of interrupts)

With `cpu_as_tag = true`:

- interrupts
  - tags:
    - irq (IRQ name)
    - type
    - device (name of the device that is located at the IRQ)
    - cpu
  - fields:
    - count (int, number of interrupts)

- soft_interrupts
  - tags:
    - irq (IRQ name)
    - type
    - device (name of the device that is located at the IRQ)
    - cpu
  - fields:
    - count (int, number of interrupts)

## Example Output

With `cpu_as_tag = false`:

```text
interrupts,irq=0,type=IO-APIC,device=2-edge\ timer,cpu=cpu0 count=23i 1489346531000000000
interrupts,irq=1,type=IO-APIC,device=1-edge\ i8042,cpu=cpu0 count=9i 1489346531000000000
interrupts,irq=30,type=PCI-MSI,device=65537-edge\ virtio1-input.0,cpu=cpu1 count=1i 1489346531000000000
soft_interrupts,irq=NET_RX,cpu=cpu0 count=280879i 1489346531000000000
```

With `cpu_as_tag = true`:

```text
interrupts,cpu=cpu6,irq=PIW,type=Posted-interrupt\ wakeup\ event count=0i 1543539773000000000
interrupts,cpu=cpu7,irq=PIW,type=Posted-interrupt\ wakeup\ event count=0i 1543539773000000000
soft_interrupts,cpu=cpu0,irq=HI count=246441i 1543539773000000000
soft_interrupts,cpu=cpu1,irq=HI count=159154i 1543539773000000000
```
