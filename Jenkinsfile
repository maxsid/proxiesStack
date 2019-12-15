pipeline {
  agent any
  options {
      buildDiscarder(logRotator(numToKeepStr: '10',
                                daysToKeepStr: '5',
                                artifactNumToKeepStr: '10',
                                artifactDaysToKeepStr: '5'))
  }
  environment {
    COMMIT_HASH = """${sh(
                       returnStdout: true,
                       script: 'git rev-parse --short=12 HEAD'
                       ).trim()}"""
    IMAGE_NAME              = "maxsid/proxies-stack"
    SOURCE_IMAGE_NAME       = "proxies-stack-source"
    BUILD_IMAGE_NAME        = "$IMAGE_NAME:$COMMIT_HASH"
    DOCKERHUB_CREDENTIAL_ID = "dockerhub-maxsid"
  }
  stages {
    stage('Build') {
      steps {
        sh "docker build -f docker/source.Dockerfile -t $SOURCE_IMAGE_NAME ."
        sh "docker build -f docker/Dockerfile -t $IMAGE_NAME ."
      }
    }
    stage('Unit Test') {
      agent { docker { image "$SOURCE_IMAGE_NAME" } }
      steps {
        sh "go test ./..."
      }
    }
    stage('Test') {
      agent {
        docker {
          image 'docker/compose'
          args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
      }
      steps {
        dir("mockTest") {
          script {
            try {
              echo "--- Running Test Environment ---"
              sh   'docker-compose up -d --build'
              sleep 10
              echo "--- Running Test ---"
              sh   "docker-compose exec -T python-tester python3 check.py"
            } catch (err) {
              echo "--- proxies-stack logs ---"
              sh   "docker-compose logs proxies-stack"
              sh "exit 1"
            } finally {
              echo "--- Stopping Test Environment ---"
              sh   "docker-compose down"
            }
          }
        }
      }
    }
    stage("Push Build to DockerHub") {
      steps {
        sh "docker tag $IMAGE_NAME $BUILD_IMAGE_NAME"
        script {
          docker.withRegistry( '', DOCKERHUB_CREDENTIAL_ID ) {
            docker.image("${env.BUILD_IMAGE_NAME}").push()
          }
        }
      }
    }
    stage("Push Release to DockerHub") {
      when { allOf { anyOf { branch "master"; branch "develop" }; tag "release-*" } }
      steps {
        script {
          def tagName = sh(returnStdout: true, script: "git tag --sort version:refname | tail -1").trim().substring(8)
          sh "docker tag $BUILD_IMAGE_NAME $IMAGE_NAME:" + tagName
          docker.withRegistry( '', DOCKERHUB_CREDENTIAL_ID ) {
            docker.image("$IMAGE_NAME:" + tagName).push()
          }
        }
      }
    }
  }
  post {
    always {
      sh "docker images | grep $IMAGE_NAME | tr -s ' ' | cut -d ' ' -f 2 | xargs -I {} docker rmi $IMAGE_NAME:{} 2> /dev/null"
      sh "docker system prune -f"
    }
  }
}