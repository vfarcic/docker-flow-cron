# Usage

TODO: Write

## Put Job

> Adds a job to the cron

The following body parameters can be used to send a *create job* `PUT` request to *Docker Flow Cron*. They should be added to the base address **[CRON_IP]:[CRON_PORT]/v1/docker-flow-cron/job/[jobName]**.

|param           |Description                                                        |Mandatory|Example  |
|----------------|-------------------------------------------------------------------|---------|---------|
|image           |Docker image.                                                      |yes      |alpine   |
|command         |The command that will be executed when a job is created.           |no       |echo "hello World"|
|schedule        |The schedule that defines the frequency of the job execution. TODO: Link to a reference from https://godoc.org/github.com/robfig/cron|yes|@every 15s|
|args            |The list of arguments that can be used with the `docker service create` command.<br><br>`--restart-condition` cannot be set to `any`. If not specified, it will be set to `none`.<br>`--name` argument is not allowed.<br><br>Any other argument supported by `docker service create` is allowed.|no|TODO|

TODO: Example

## Get All Jobs

TODO: Write

## Get Job

TODO: Write

## Delete Job

TODO: Write