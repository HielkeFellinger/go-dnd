# go-dnd


# **DEPRECATED** Now: https://github.com/HielkeFellinger/dramatic_gopher
### Lessons have been learned; the number of unknown unknows has been reduced. Now this V-1 code will be left as is, see the above link for a newer and cleaner V0


[![Build and Test](https://github.com/HielkeFellinger/go-dnd/actions/workflows/go.yml/badge.svg)](https://github.com/HielkeFellinger/go-dnd/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/HielkeFellinger/go-dnd/graph/badge.svg?token=JXQX5TZOXE)](https://codecov.io/github/HielkeFellinger/go-dnd)

- Requires PostgreSql DB (Use `docker-compose.yml`)
- Run by `docker-compose.yml` or manual building source. Develop mode via [air](https://github.com/air-verse/air) 

This will, most likely, be nothing more than a simple gui+engine for a ""D&D""/TTRPG oneshot, 
Build for entertainment (On LAN), with friends and getting familiar with Go.
Do NOT look at this software as an example of anything.

- Some Security features will be used; but will not be sufficient for more than "localhost" LAN use
- A Clean separation html/js/css will be done later if at all
- Most libraries used are chosen on nothing more than "Looks cool, I want to try that" or 
"I want to build something like that myself"

Art:
- Maps are made by me using [inkarnate](https://inkarnate.com) and are free to be used, as a whole.
  - Individual elements are still property of inkarnate
- Black and white maps are made using [dungeonscrawl](https://www.dungeonscrawl.com/)
- AI "Art" is used to temp. fill in some holes (player tokens) and will be removed asap
- System parsing will be done by https://expr-lang.org