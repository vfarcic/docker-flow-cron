import java.text.SimpleDateFormat

pipeline {
  agent {
    label "test"
  }
  options {
    buildDiscarder(logRotator(numToKeepStr: '2'))
    disableConcurrentBuilds()
  }
  stages {
    stage("build") {
      steps {
        script {
          def dateFormat = new SimpleDateFormat("yy.MM")
          currentBuild.displayName = dateFormat.format(new Date()) + "." + env.BUILD_NUMBER
        }
        sh "docker image build -t vfarcic/docker-flow-cron ."
        sh "docker tag vfarcic/docker-flow-cron vfarcic/docker-flow-cron:beta"
        withCredentials([usernamePassword(
          credentialsId: "docker",
          usernameVariable: "USER",
          passwordVariable: "PASS"
        )]) {
          sh "docker login -u $USER -p $PASS"
        }
        sh "docker push vfarcic/docker-flow-cron:beta"
        sh "docker image build -t vfarcic/docker-flow-cron-test -f Dockerfile.test ."
        sh "docker push vfarcic/docker-flow-cron-test"
        sh "docker image build -t vfarcic/docker-flow-cron-docs -f Dockerfile.docs ."
      }
    }
    stage("test") {
      environment {
        HOST_IP = "build.dockerflow.com"
        DOCKER_HUB_USER = "vfarcic"
      }
      steps {
        sh "docker-compose -f docker-compose-test.yml run --rm unit"
      }
    }
    stage("release") {
      when {
        branch "master"
      }
      steps {
        withCredentials([usernamePassword(
          credentialsId: "docker",
          usernameVariable: "USER",
          passwordVariable: "PASS"
        )]) {
          sh "docker login -u $USER -p $PASS"
        }
        sh "docker tag vfarcic/docker-flow-cron vfarcic/docker-flow-cron:${currentBuild.displayName}"
        sh "docker push vfarcic/docker-flow-cron:${currentBuild.displayName}"
        sh "docker push vfarcic/docker-flow-cron"
        sh "docker tag vfarcic/docker-flow-cron-docs vfarcic/docker-flow-cron-docs:${currentBuild.displayName}"
        sh "docker push vfarcic/docker-flow-cron-docs:${currentBuild.displayName}"
        sh "docker push vfarcic/docker-flow-cron-docs"
      }
    }
    stage("deploy") {
      when {
        branch "master"
      }
      agent {
        label "prod"
      }
      steps {
        sh "helm upgrade -i docker-flow-swarm-listener helm/docker-flow-swarm-listener --namespace df --set image.tag=${currentBuild.displayName}"
      }
    }
  }
  post {
    always {
      sh "docker system prune -f"
    }
    failure {
      slackSend(
        color: "danger",
        message: "${env.JOB_NAME} failed: ${env.RUN_DISPLAY_URL}"
      )
    }
  }
}
