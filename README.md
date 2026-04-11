```
           .     .                          
      ..   8"=,,88,   _.                    
       8""=""8'  "88a88'                    
  .. .;88m a8   ,8"" "8       ▄▄▄▄ ▄▄ ▄▄ ▄▄  ▄▄▄ 
   "8"'  "88"  A"     8;     ██▀▀▀ ██▄██ ██ ██▀██
    "8   8,  8,       "8     ▀████ ██ ██ ██ ██▀██
     8,  "8, "8,    ___8,      ▄▄▄▄▄ ▄▄     ▄▄▄  ▄▄▄▄  ▄▄  ▄▄▄▄ ▄▄▄▄▄▄              
     "8,  "8, "8mm""""""8m.    ██▄▄  ██    ██▀██ ██▄█▄ ██ ███▄▄   ██   
      "8,am888i"'   ,mm"     ▄ ██    ██▄▄▄ ▀███▀ ██ ██ ██ ▄▄██▀   ██   
      ,8"  _8"  .m888"                                   
     ,88P"""""I888888        A modular, scalable e-commerce system
     "'         "I888        designed to simulate a real-world
                  "I8_              production architecture —
      ,mmeem.m""i, I8""  ,mmeem,'.     combining modern frontend,
     m""    . "8.8 I8  ,8"   .  "88     backend, security, and AI
    i8  . '  ,mi""8I8 ,8 . '  ,8" 88     intelligence layers.
    88.' ,mm""    "8I88"m,,mm'"    8     
    "8_m""         "I8   ""'            (c) AETHERGEAR
     "8             I8_             
```

## Project Structure

```
chia-florist/
│
├── control-panel        # Admin dashboard (management & monitoring UI)
├── docs                 # Documentation, diagrams, API specs
├── e-commerce           # Customer-facing frontend (Next.js / Vue)
├── intelligence-layer   # AI services (FastAPI)
├── service-core         # Core backend services (Golang)
├── README.md
```

## Architecture Overview

This project is built with a **modular architecture**, with each layer has their own architecture and framework:

* **E-commerce Layer**
  Frontend application for customers (product browsing, checkout, etc.)

* **Service Core**
  Backend API handling business logic, transactions, and data processing

* **Control Panel**
  Admin interface for monitoring, analytics, and management

* **Intelligence Layer**
  AI-powered services (recommendation, analytics, anomaly detection)

* **(Optional Future) WAF Layer**
  Security layer for filtering, monitoring, and protecting requests

## Tech Stack

| Layer           | Technology       |
| --------------- | ---------------- |
| Frontend        | Next.js / Vue    |
| Backend         | Golang           |
| AI / ML         | FastAPI (Python) |
| Database        | PostgreSQL, Supabase       |
| DevOps (future) | Docker / Cloud   |

## Objective

- [ ] Deliver a **working modular system before June**
- [ ] Enable **parallel development across multiple teams**
- [ ] Provide a **real-world backend + AI integration experience**
- [ ] Build a **portfolio-ready system with production mindset**
- [ ] Ensure all development is **tracked via GitHub Issues**

## Development Workflow

* All tasks must be tracked in **GitHub Issues**
* Each module can be developed independently
* Start with **mock data → integrate later**
* Focus on:

  * API contracts first
  * Clear separation of concerns
  * Incremental delivery

## Philosophy

This project is not just about building features.

It’s about learning how to:

* design scalable systems
* collaborate across domains
* think in architecture, not just code

## Notes for Developers

* Don’t over-engineer early — **make it work first**
* Keep modules loosely coupled
* Use mock data when blocked
* Refactor as the system evolves

## Contribution

* Pick tasks from **Issues**
* Keep PRs small and focused
* Communicate across modules when needed

## Let’s Build Something Real

You’re not just building an app.
You’re building a system.

Ship fast. Learn faster. Improve continuously.

```
                  .-.'  '.-.
               .-(   \  /   )-.
              /   '..oOOo..'   \
      ,       \.--.oOOOOOOo.--./
      |\  ,   (   :oOOOOOOo:   )
     _\.\/|   /'--'oOOOOOOo'--'\
     '-.. ;/| \   .''oOOo''.   /
     .--`'. :/|'-(   /  \   )-'
      '--. `. / //'-'.__.'-;
        `'-,_';//     ,  /|
             '((     .|\/./_
               \\ |\: .'`--.
(❁´◡`❁)       \\, .' .--'
(c) chia.florist ))'_,-'`
                //-'
grow • build   // 
• bloom       |/
```