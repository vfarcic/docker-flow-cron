# Docker Flow Cron

This project is in the **design phase**.

## High Level Design

The project should be able to:

- [X] schedule jobs using cron syntax
- [X] run jobs as Swarm services
- [ ] schedule jobs through HTTP requests
- [ ] retrieve the list of scheduled jobs through an HTTP request
- [ ] remove scheduled jobs
- [ ] use `constraint`, `reserve-cpu`, `reserve-memory`, and `label` arguments
- [ ] rerun failed jobs
- [ ] send notifications when a job fails
- [ ] provide dashboard for job monitoring (explore the option to extend Portainer or Grafana)
- [ ] provide UI for administration (explore the option to extend Portainer)
- [ ] be fault tolerant
- [X] runs jobs detached from the cron so that they are unaffected in case of a failure
- [ ] after cron rescheduling, it does not start jobs that are already running
- [ ] time out jobs
- [ ] get status of past jobs (status)
- [ ] host pre-made jobs (e.g. docker prune, Mongo backup, etc)
- [ ] remove job services after a defined period

## Assumptions

* Docker experimental features are enabled while `docker service logs` is in it
* System has `/bin/sh`
* Binary is running on a manager node
