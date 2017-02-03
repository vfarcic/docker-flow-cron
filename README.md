# Docker Flow Cron

This project is in the **design phase**.

## High Level Design

The project should be able to:

- [X] schedule jobs using cron syntax
- [X] runs jobs detached from the cron so that they are unaffected in case of a failure
- [X] run jobs as Swarm services
- [X] schedule jobs through HTTP requests
- [X] use `constraint`, `reserve-cpu`, `reserve-memory`, and `label` arguments
- [X] rerun failed jobs
- [X] retrieve the list of scheduled jobs through an HTTP request
- [ ] retrieve job details with the list of executions through an HTTP request
- [ ] remove scheduled jobs
- [ ] be fault tolerant
- [ ] after cron rescheduling, it does not start jobs that are already running
- [ ] send notifications when a job fails
- [ ] provide dashboard for job monitoring (explore the option to extend Portainer or Grafana)
- [ ] time out jobs
- [ ] provide UI for administration (explore the option to extend Portainer)
- [ ] host pre-made jobs (e.g. docker prune, Mongo backup, etc)
- [ ] prune old executions (services)
- [ ] create a CLI

## Tasks

- [ ] develop PoC
- [ ] integration tests
- [ ] document (use cron.AddJob as the base for the rules)
- [ ] release
- [ ] reduce testing time or split into groups (e.g. unit tests, integration tests)

## Assumptions

* Docker experimental features are enabled while `docker service logs` is in it
* System has `/bin/sh`
* Binary is running on a manager node
