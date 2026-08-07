How computer-A communicate to computer-B(server), in early phase  
Every Server Has an Identity: IP Address

Every computer connected to the Internet has a unique identifier called an IP Address—like 192.168.1.5 or 142.250.183.14—which tells you the exact machine you want to talk to.

But it's difficult for humans to remember numbers. Imagine this:

- You want to visit Google
- You must remember: 142.250.183.14
- For Facebook, another number
- For Amazon, another number

 First Solution: HOSTS.TXT

To solve this problem, an early solution called HOSTS.TXT was introduced:

- A simple text file (like an Excel sheet)
- Stored mappings like:


142.250.183.14    google
157.240.18.35     facebook


When a user typed a name like google, the computer:

1. Looked into the local HOSTS.TXT file
2. Found the corresponding IP address
3. Connected to that server

But this won’t so long, As the Internet grew rapidly, this approach became impossible to manage:


- Millions of servers were added
- File size kept increasing
- Name conflicts occurred (same name for different IPs)
- Constant manual updates were required

 Birth of DNS (Domain Name System)

To solve this scalability problem, DNS was invented.

Instead of one giant list, DNS introduced a distributed, hierarchical system.

DNS Hierarchy
- Root (.) → Locates top-level domains
- TLD (.com, .org, .edu) → Directs to domain owners
- Domain (google, amazon)→ Points to servers
- Subdomain (maps, mail)→ Specific services

Each level only knows the next level.