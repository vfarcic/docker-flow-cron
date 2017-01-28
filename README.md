# Docker Flow Cron

This project is in the **design phase**.

## High Level Design

The project should be able to:

- [ ] schedule jobs using cron syntax
- [ ] run jobs as Swarm services
- [ ] remove scheduled jobs
- [ ] use `constraint`, `reserve-cpu`, `reserve-memory`, and `label` arguments
- [ ] rerun failed jobs
- [ ] send notifications when a job fails
- [ ] provide dashboard with job statuses (explore the option to extend Portainer)
- [ ] provide UI for administration (explore the option to extend Portainer)
- [ ] be fault tolerant
- [ ] time out
- [ ] get status of past jobs (status)
- [ ] host pre-made jobs (e.g. docker prune, Mongo backup, etc)
- [ ] remove job services after a defined period

## Assumptions

* Docker experimental features are enabled while `docker service logs` is in it
* System has `/bin/sh`
* Binary is running on a manager node
