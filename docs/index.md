# Running Docker Flow Cron In a Swarm Cluster

TODO: Review

Docker Swarm services are designed to the long lasting processes that, potentially, live forever. Docker does not have a mechanism to schedule jobs based on time interval or, to put it in other words, it does not have the ability to use cron-like syntax for time-based scheduling.

*Docker Flow Cron* is designed overcome some of the limitations behind Docker Swarm services and provide cron-like time-based scheduling while maintaining fault tolerance features available in Docker Swarm.
