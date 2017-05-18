# Docker Flow Cron

TODO: Move unfinished to issues

TODO: Write an introduction and a reference to [cron.dockerflow.com](cron.dockerflow.com).

This project is in the **design phase**.

## High Level Design

The project should be able to:

- [X] schedule jobs using cron syntax
- [X] runs jobs detached from the cron so that they are unaffected in case of a failure
- [X] run jobs as Swarm services
- [X] schedule jobs through HTTP requests
- [X] use `constraint`, `reserve-cpu`, `reserve-memory`, and `label` arguments
- [X] rerun failed jobs
- [X] retrieve the list of scheduled jobs
- [X] retrieve job details
- [X] retrieve job executions
- [X] retrieve job execution logs (`docker service logs`)
- [X] remove scheduled jobs
- [X] be fault tolerant
- [ ] use a combination of a job and and an index as a service name
- [ ] after cron rescheduling, it does not start jobs that are already running
- [ ] send notifications when a job fails
- [ ] provide dashboard for job monitoring (explore the option to extend Portainer or Grafana)
- [ ] time out jobs
- [ ] provide UI for administration (explore the option to extend Portainer)
- [ ] host pre-made jobs (e.g. docker prune, Mongo backup, etc)
- [ ] prune old executions (services)
- [ ] create a CLI
- [ ] provide commands that will allow users to filter services with Docker client
- [ ] update scheduled jobs
- [ ] prune old services
- [ ] support stacks

## Tasks

- [ ] develop PoC
- [ ] document (use cron.AddJob as the base for the rules)
- [ ] release

## Assumptions

* Docker experimental features are enabled while `docker service logs` is in it
* System has `/bin/sh`
* Binary is running on a manager node

<a href='https://ko-fi.com/A655LRB' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://az743702.vo.msecnd.net/cdn/kofi2.png?v=0' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>