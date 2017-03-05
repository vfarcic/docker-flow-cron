# Docker Flow Cron Requirements

The following set of items are required for *Docker Flow Cron*.

## A Cluster Running in Docker Swarm Mode

*Docker Flow Cron* creates Docker Swarm services for each scheduled job execution.

## SH

Operating system used to host Docker Swarm must be able to run `/bin/sh` commands.

## The service running on one of the Swarm manager nodes

*Docker Flow Cron* requires interaction with one of the Swarm managers to schedule services that run as scheduled jobs.
