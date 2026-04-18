# NetBird skill tree — a beginner's prerequisite map

> **What this is.** A Math-Academy-Way-style prerequisite DAG for learning the [NetBird](https://github.com/netbirdio/netbird) codebase, layered with the pedagogy of the `codebase-study-guide` skill (threshold concepts, Mermaid maps, PRIMM exploration tasks). 15 atomic nodes, each with explicit prereqs, behavior-observable mastery criteria, and a timed drill.
>
> **Who it's for.** You can read Go, but you're new to networking, VPNs, WireGuard, and zero-trust. You want to understand *how NetBird actually works*, not just use it.
>
> **Fun-experiment status.** This file is an exploration — it's not wired into the `go-dojo verify` system (NetBird is a third-party codebase; there are no per-task tests we control). The frontmatter matches `curriculum-v2.md`'s shape so `go-dojo` *could* parse it if we ever want to, but for now treat it as a reading order with mastery gates you self-assess against.

---

## NetBird in one paragraph

NetBird is an open-source mesh VPN. Each machine runs a daemon that brings up a WireGuard interface and connects to a central **management server** for network config + policy. A second service, **signal**, relays encrypted ICE candidates so peers can find each other and establish direct WireGuard tunnels. A third service, **relay**, is a TURN fallback for when NAT makes direct P2P impossible. The magic is that the data plane (WireGuard packets) is end-to-end encrypted between peers — management and signal never see plaintext traffic; they only distribute config and coordinate the handshake.

---

## Threshold concepts (read these three paragraphs before anything else)

These are the ideas that, once you hold them, make every file in the repo click. If you're ever lost, come back here.

1. **Management is the source of truth; peers are stateless reconcilers.** The server owns accounts, peers, policies, groups, routes. A peer logs in and asks "what should my network look like?" and gets a **NetworkMap** — an immutable snapshot. The peer then reconciles its local WireGuard interface, firewall, DNS, and routing to match. Peers don't store policy; they execute. This is why updates propagate live over a long-lived gRPC `Sync()` stream.

2. **Control plane and data plane are deliberately separate.** The **signal service** is a dumb relay that moves encrypted ICE candidates between peers. It never sees WireGuard keys, peer IPs, or traffic. The **relay (TURN) service** sees only encrypted WireGuard packets — it can't eavesdrop either, because WireGuard provides end-to-end encryption. This separation is what makes signal stateless and relay a dumb tunnel. If you try to understand NetBird as "one big service," you'll drown; think of it as *three services that only exchange ciphertext*.

3. **Everything between components uses WireGuard keypairs as its identity.** Client ↔ management messages are encrypted using the client's WireGuard private key and the server's WireGuard public key. Peer ↔ peer messages through signal are encrypted with WireGuard keys too. The same key that builds the VPN tunnel also authenticates RPCs. Once you see this, the line between "crypto" and "protocol" blurs in the right way.

---

## System map

```mermaid
graph TD
    subgraph Peers
        A[Peer A daemon]
        B[Peer B daemon]
    end

    subgraph "Control services"
        M[management :8443<br/>gRPC]
        S[signal :50051<br/>gRPC]
        R[relay :3478<br/>TURN]
    end

    subgraph "Persistence"
        DB[(SQLite / Postgres / MySQL)]
        IDP[OIDC IdP<br/>Okta / Dex / Google]
    end

    A -- "Login + Sync stream" --> M
    B -- "Login + Sync stream" --> M
    M --> DB
    M <-.-> IDP

    A -- "Encrypted ICE offer/answer" --> S
    S -- "relay to peer" --> B
    B -- "Encrypted ICE answer" --> S
    S -- "relay to peer" --> A

    A <== "WireGuard tunnel (direct)" ==> B
    A -. "WireGuard via relay (NAT fallback)" .-> R
    R -. .-> B
```

Four arrows to understand here:
- **Peer → management (solid)**: login + long-lived Sync stream pushing NetworkMap updates.
- **Peer → signal (solid)**: encrypted ICE exchange for NAT traversal.
- **Peer ↔ peer (thick direct)**: the WireGuard data tunnel. All user traffic flows here.
- **Peer → relay (dotted)**: TURN fallback when direct P2P fails.

---

## How mastery works (self-assessed here, since we're not running tests)

| Level | You know it when… |
|---|---|
| `unseen` | Haven't started |
| `learning` | You can answer the node's "explain it back" question *after* reading |
| `proficient` | You can answer it from memory, 3+ days later, without re-reading |
| `automatic` | You beat the drill's target time, and you can trace the relevant flow across files while talking through it |

Review intervals: 3 / 7 / 21 / 60 days after each pass.

---

## The tree

15 nodes across 5 stages. Stages A and B can interleave; C depends on both; D and E depend on C.

### Stage 1 — Prereqs from outside the repo

These aren't NetBird — they're the networking/crypto primitives you need to read NetBird's code without pattern-matching in the dark.

```yaml
---
id: nb.1.1-pubkey-crypto
phase: 1
title: Public-key crypto for networking
prereqs: []
interleave_with: [nb.1.2-wireguard-primer]
mastery_criteria: |
  - Explain asymmetric encryption in 3 sentences: what private vs public keys are for, why you can share the public one
  - Define a nonce and why reusing one is catastrophic
  - Describe a Diffie-Hellman key exchange in plain English — what gets exchanged publicly, what each side derives privately
drill_ids: [crypto-primitives-recall]
review_intervals_days: [3, 7, 21, 60]
remediation: []
diagnostic_checkpoint: false
---
```

**Why this first.** Every other concept in NetBird stands on this. WireGuard, ICE-over-signal encryption, management mTLS — they're all the same primitive reused.

**Explore task (PRIMM).** *Predict:* if a peer's WireGuard private key leaks, what can an attacker do and not do? Write your guess, then ask yourself after nb.1.2.

---

```yaml
---
id: nb.1.2-wireguard-primer
phase: 1
title: WireGuard in one page
prereqs: [nb.1.1-pubkey-crypto]
interleave_with: []
mastery_criteria: |
  - Name the four WireGuard concepts: interface, peer, public/private keypair, allowed IPs
  - Explain why WireGuard is stateless and what the handshake establishes
  - Draw the packet path: app → kernel → WireGuard → encrypt → UDP → remote → decrypt → app
drill_ids: [wireguard-packet-path]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.1.1-pubkey-crypto]
diagnostic_checkpoint: false
---
```

**Reading.** [WireGuard whitepaper](https://www.wireguard.com/papers/wireguard.pdf) §§1–3 (skip the formal security proofs for now). Then `wg-quick` man page.

**Explore task.** On your own machine: `wg genkey | tee privatekey | wg pubkey > publickey` — look at what you created. Now you've made a WireGuard identity. That file is all a peer *is*.

---

```yaml
---
id: nb.1.3-nat-and-stun
phase: 1
title: NAT types + STUN + UDP hole punching
prereqs: []
interleave_with: [nb.1.2-wireguard-primer]
mastery_criteria: |
  - Distinguish full-cone, restricted-cone, symmetric NAT; explain why symmetric breaks P2P
  - Explain what STUN returns and why a peer would ask
  - Articulate the core trick of UDP hole punching
drill_ids: [nat-traversal-explain]
review_intervals_days: [3, 7, 21, 60]
remediation: []
diagnostic_checkpoint: false
---
```

**Why.** NetBird's entire ICE dance exists because NAT is everywhere. Until you feel *why* two home routers can't just open a TCP connection to each other, the signal service looks overcomplicated.

**Explore task.** From two different networks, run `stun stun.l.google.com:19302` (install `stuntman-client`). Compare the mapped addresses. That's what your peer is doing.

---

```yaml
---
id: nb.1.4-grpc-streams
phase: 1
title: gRPC + Protobuf + bidirectional streams
prereqs: []
interleave_with: [nb.1.1-pubkey-crypto]
mastery_criteria: |
  - Explain the four gRPC call types: unary, server-stream, client-stream, bidi-stream
  - Read a .proto file and identify services, messages, and streaming RPCs
  - Articulate when you'd pick a bidi stream over polling
drill_ids: [proto-file-read]
review_intervals_days: [3, 7, 21, 60]
remediation: []
diagnostic_checkpoint: false
---
```

**Reading.** `go doc google.golang.org/grpc`. Then skim [gRPC concepts](https://grpc.io/docs/what-is-grpc/core-concepts/).

**Explore task.** Open `shared/management/proto/management.proto`. Without reading implementation, answer: which RPCs are unary? Which are streaming? Which direction(s)?

---

### Stage 2 — NetBird anatomy (components in isolation)

You can read these in any order; pick whichever directory you're least intimidated by first.

```yaml
---
id: nb.2.1-three-services
phase: 2
title: The three services — management, signal, relay
prereqs: [nb.1.4-grpc-streams, nb.1.2-wireguard-primer]
interleave_with: [nb.2.2-data-model]
mastery_criteria: |
  - For each of management, signal, relay: state its single job in one sentence
  - Explain why signal is stateless and relay sees only ciphertext
  - Name which binary gets built from which top-level directory
drill_ids: [three-services-explain]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.1.4-grpc-streams]
diagnostic_checkpoint: false
---
```

**Reading.** Top-level README, then `management/cmd/`, `signal/cmd/`, `relay/cmd/` — each just 30–50 lines wiring cobra → server.

**Explore task.** Run `find . -maxdepth 2 -name main.go` inside `/tmp/netbird`. Count the mains. What does each turn into?

---

```yaml
---
id: nb.2.2-data-model
phase: 2
title: Account / Peer / SetupKey / Group / Policy
prereqs: [nb.1.4-grpc-streams]
interleave_with: [nb.2.1-three-services]
mastery_criteria: |
  - Draw the ER diagram: Account → Peer, Account → User, Account → SetupKey, Account → Group → Policy
  - Explain the lifecycle of a SetupKey (when it's created, used, retired)
  - Explain a Group — what it contains, why it exists separately from Peer
drill_ids: [data-model-er]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.2.1-three-services]
diagnostic_checkpoint: false
---
```

**Reading.** `management/server/types/account.go`, `management/server/types/peer.go`, `management/server/types/setupkey.go`, `management/server/types/group.go`, `management/server/types/policy.go`.

**Explore task (PRIMM).** *Predict:* if a user is in Group "Engineers" and Policy allows "Engineers → Servers", what happens when their laptop (a peer) is added to the account? Look in `management/server/policy.go` and check yourself.

---

```yaml
---
id: nb.2.3-networkmap
phase: 2
title: NetworkMap — the immutable per-peer view
prereqs: [nb.2.2-data-model]
interleave_with: []
mastery_criteria: |
  - Explain what fields live inside a NetworkMap and why the peer needs each
  - Articulate why NetworkMap is built per-peer (two peers in the same account receive different maps)
  - Identify the builder function and the two main inputs (Account + Peer identity)
drill_ids: [networkmap-field-tour]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.2.2-data-model]
diagnostic_checkpoint: true
---
```

**Reading.** `management/server/types/networkmap.go` (struct definition), `management/server/types/networkmapbuilder.go`, `management/server/types/networkmap_golden_test.go` — the golden tests are excellent self-documentation.

**Explore task.** Read one of the golden fixtures and hand-calculate what NetworkMap a specific peer should get. Then run the test.

---

### Stage 3 — End-to-end flows (how the components actually talk)

Once the pieces make sense individually, trace the wires.

```yaml
---
id: nb.3.1-peer-login
phase: 3
title: Peer login — from setup key to allocated IP
prereqs: [nb.2.3-networkmap, nb.2.1-three-services]
interleave_with: [nb.3.2-sync-stream]
mastery_criteria: |
  - Trace LoginRequest from client/cmd/login.go to management/server/account.go::LoginPeer, naming every file
  - Explain when a JWT is used vs a setup key, and which RPC handles each
  - Describe IP allocation: where the network CIDR comes from and how collisions are avoided
drill_ids: [login-trace]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.2.3-networkmap]
diagnostic_checkpoint: false
---
```

**Reading.** `client/cmd/login.go`, `client/internal/connect.go::loginToManagement`, `shared/management/proto/management.proto` (Login RPC), `management/server/account.go::LoginPeer`, `management/server/setupkey.go`.

**Explore task (PRIMM).** *Predict:* what happens if the same setup key is used from two machines simultaneously? Read the code and confirm. (Hint: usage-limit logic.)

---

```yaml
---
id: nb.3.2-sync-stream
phase: 3
title: Sync() — the long-lived update stream
prereqs: [nb.3.1-peer-login, nb.1.4-grpc-streams]
interleave_with: []
mastery_criteria: |
  - Describe what triggers a new message on the Sync stream (peer added, policy changed, route added, etc.)
  - Explain how the client applies an incoming NetworkMap (WireGuard reconfig + firewall + DNS + routes)
  - Identify what happens if the stream drops: reconnect cadence, re-login requirements
drill_ids: [sync-trigger-trace]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.3.1-peer-login, nb.1.4-grpc-streams]
diagnostic_checkpoint: false
---
```

**Reading.** `management/server/account.go::SyncAndMarkPeer`, `client/internal/connect.go` (the run-loop around lines 200–300), `client/internal/engine.go::Engine.updateNetworkMap`.

**Explore task.** Add a log line to `SyncAndMarkPeer` (don't commit) and run two peers locally. Change a policy. Count how many Sync frames land.

---

```yaml
---
id: nb.3.3-peer-discovery-ice
phase: 3
title: Peer discovery — ICE candidates through signal
prereqs: [nb.1.3-nat-and-stun, nb.2.1-three-services]
interleave_with: [nb.3.4-relay-fallback]
mastery_criteria: |
  - Trace one ICE candidate from local gather → signal SendOffer → peer → SignalExchange.Receive
  - Explain what Encrypted field contains in a signal message and why signal can't read it
  - Describe how the ICE agent picks a winning candidate pair
drill_ids: [ice-candidate-path]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.1.3-nat-and-stun]
diagnostic_checkpoint: false
---
```

**Reading.** `client/internal/peer/conn.go::Conn.establish`, `client/internal/peer/worker_ice.go`, `client/internal/peer/ice/agent.go`, `client/internal/peer/signaler.go`, `signal/server/relay.go`, `shared/signal/proto/signalexchange.proto`.

**Explore task.** Read the EncryptedMessage proto. What's the difference between `Body` (cleartext) and `Encrypted` (ciphertext)? Which fields does signal ever need to look at?

---

```yaml
---
id: nb.3.4-relay-fallback
phase: 3
title: Relay — TURN fallback with HMAC auth
prereqs: [nb.3.3-peer-discovery-ice]
interleave_with: []
mastery_criteria: |
  - Explain when ICE gives up and relay kicks in (symmetric NAT scenario)
  - Describe the HMAC auth token: who issues, who verifies, what it binds
  - Argue why relay is not a privacy regression — what it sees vs doesn't see
drill_ids: [relay-threat-model]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.3.3-peer-discovery-ice]
diagnostic_checkpoint: true
---
```

**Reading.** `relay/server/`, `shared/relay/auth/hmac/`, `client/internal/peer/conn.go` (relay fallback path).

**Explore task.** Forge a bad HMAC token in a test. Confirm relay rejects it.

---

### Stage 4 — Client internals

```yaml
---
id: nb.4.1-engine-reconciliation
phase: 4
title: Engine — the client reconciliation loop
prereqs: [nb.3.2-sync-stream]
interleave_with: [nb.4.2-platform-wireguard]
mastery_criteria: |
  - Identify what Engine owns: WireGuard interface, peer connections, firewall, DNS, routes
  - Explain how Engine decides to add/remove a peer connection when NetworkMap changes
  - Trace one policy-change event from Sync frame arrival to firewall rule update
drill_ids: [engine-ownership]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.3.2-sync-stream]
diagnostic_checkpoint: false
---
```

**Reading.** `client/internal/engine.go` (whole file — it's the heart of the daemon), `client/internal/peer/conn.go::Conn`.

**Explore task.** Draw Engine's state machine. What's the loop? What events advance it?

---

```yaml
---
id: nb.4.2-platform-wireguard
phase: 4
title: Platform-specific WireGuard device bring-up
prereqs: [nb.1.2-wireguard-primer, nb.4.1-engine-reconciliation]
interleave_with: []
mastery_criteria: |
  - Name the three platform paths: Linux (netlink + kernel module), macOS (utun + userspace), Windows (WinTun + userspace)
  - Explain why userspace exists even when kernel WireGuard is available
  - Locate the configure interface for each platform and list one method
drill_ids: [platform-configurer-tour]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.1.2-wireguard-primer]
diagnostic_checkpoint: false
---
```

**Reading.** `client/iface/device/` (TUN abstraction), `client/iface/configurer/` (platform-specific config).

**Explore task.** Pick the platform you're on. Find where the interface MTU is set. Change it in a branch, observe what breaks at different values.

---

### Stage 5 — Server internals & policy enforcement

```yaml
---
id: nb.5.1-policy-enforcement
phase: 5
title: Policy + ACL — from rule to firewall
prereqs: [nb.2.2-data-model, nb.4.1-engine-reconciliation]
interleave_with: [nb.5.2-posture-checks]
mastery_criteria: |
  - Trace a Policy row in the DB all the way to a kernel firewall rule on a peer
  - Explain the Group → Policy → Peer expansion: how a rule "Engineers → Servers" becomes concrete peer-to-peer allowances
  - Explain why firewall enforcement lives on the client, not on a gateway
drill_ids: [policy-expansion-trace]
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.2.2-data-model, nb.4.1-engine-reconciliation]
diagnostic_checkpoint: false
---
```

**Reading.** `management/server/policy.go`, `client/internal/acl/`.

**Explore task (PRIMM).** *Predict:* if you delete a peer from a Group, how fast does the firewall on remote peers update? Read the Sync push path to check.

---

```yaml
---
id: nb.5.2-posture-checks
phase: 5
title: Posture checks — device compliance gating
prereqs: [nb.5.1-policy-enforcement]
interleave_with: []
mastery_criteria: |
  - Describe three posture check types (OS version, geo, process, etc.) and what they evaluate
  - Explain where posture results flow: client → management → policy eval → Sync update to other peers
  - Argue why posture is a server-side decision even though the check runs client-side
drill_ids: []
review_intervals_days: [3, 7, 21, 60]
remediation: [nb.5.1-policy-enforcement]
diagnostic_checkpoint: true
---
```

**Reading.** `management/server/posture_checks.go`, `client/internal/peer/guard/`.

**Explore task.** Write a fake posture check that always fails. Log in from a machine — does the peer get admitted? Where in the code is the gate?

---

## Drill library

Short, timed exercises for automaticity. "Target" is what you're aiming for after you've read the relevant node.

```yaml
drills:
  - id: crypto-primitives-recall
    target_seconds: 180
    prompt: "On paper, in 3 minutes: define asymmetric crypto, nonce, and Diffie-Hellman. No references."
  - id: wireguard-packet-path
    target_seconds: 180
    prompt: "Draw the WireGuard packet path, app to app, naming every layer. 3 min."
  - id: nat-traversal-explain
    target_seconds: 300
    prompt: "Explain why two peers behind symmetric NATs can't connect directly, and what relay does about it. 5 min, no references."
  - id: proto-file-read
    target_seconds: 300
    prompt: "Open shared/management/proto/management.proto. In 5 min, list every RPC with its streaming type (unary/server-stream/client-stream/bidi)."
  - id: three-services-explain
    target_seconds: 180
    prompt: "In 3 minutes: one sentence per service (management, signal, relay) describing its single job. No references."
  - id: data-model-er
    target_seconds: 300
    prompt: "Draw the ER diagram: Account, User, Peer, SetupKey, Group, Policy. Arrows = foreign keys. 5 min."
  - id: networkmap-field-tour
    target_seconds: 300
    prompt: "From memory: list 6 fields inside a NetworkMap and why the peer needs each."
  - id: login-trace
    target_seconds: 600
    prompt: "Trace a Login RPC from client CLI to management DB, naming every file touched. 10 min."
  - id: sync-trigger-trace
    target_seconds: 600
    prompt: "Name 4 events that cause a new Sync frame to fire, and for each, identify the server-side function that emits it."
  - id: ice-candidate-path
    target_seconds: 600
    prompt: "Trace one ICE candidate: local gather → signal SendOffer → remote receive → agent decides. 10 min, with file names."
  - id: relay-threat-model
    target_seconds: 300
    prompt: "What does relay see? What does it not see? What can a compromised relay do? 5 min."
  - id: engine-ownership
    target_seconds: 300
    prompt: "List what Engine owns and what it doesn't, from memory. 5 min."
  - id: platform-configurer-tour
    target_seconds: 300
    prompt: "Name the three platform paths (Linux/macOS/Windows) and locate one configure method each."
  - id: policy-expansion-trace
    target_seconds: 600
    prompt: "Given rule 'Engineers → Servers' and 3 peers (2 Engineers, 1 Server), write the exact firewall rules each peer ends up with."
```

---

## Where to read first — quick lookups by role

| If you care about… | Start at |
|---|---|
| The whole system | nb.2.1 → nb.3.1 → nb.3.2 (two days) |
| Writing a new policy type | nb.2.2 → nb.5.1 |
| NAT traversal bugs | nb.1.3 → nb.3.3 → nb.3.4 |
| Client daemon behavior | nb.4.1 → nb.4.2 |
| Self-hosting it | nb.2.1 + README + `combined/` |

---

## Exit criterion — you've "got" NetBird when:

1. You can redraw the system map from memory.
2. You can trace any of the three flows (login, discovery, relay fallback) across files without looking up the path.
3. You can explain threshold concept #1 (management = source of truth, peers = reconcilers) in a way your non-technical friend understands.
4. You can predict, before reading the code, what a feature change (new policy field, new IdP) would touch — and be right most of the time.

If you hit all four, congratulations: you're no longer a NetBird tourist.
