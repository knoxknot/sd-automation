pipeline {
  agent any
  options {
    buildDiscarder(logRotator(numToKeepStr:'5', artifactNumToKeepStr: '5', artifactDaysToKeepStr: '7'))
  }
  environment {
    IMAGE_TAG="sh returnStdout: true, script:git tag -l | tee tags | tail -n1 tags"
    IMAGE_PATH="knoxknot/boardingapi"
  }
  stages {
    stage ('Initialize') {
      steps {
        checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'github', url: 'https://github.com/knoxknot/sd-automation.git']]])
      }
    }
    stage ('Build App'){
      steps {
            sh 'printenv'
            echo "Building App..."
            ws("${env.WORKSPACE}/application/") {
                sh "/usr/local/go/bin/go env | export GO111MODULE=on"
                sh "/usr/local/go/bin/go mod download"
                sh "/usr/local/go/bin/go build . -o boardingapi"
            }
        }
    }
    stage ('Build Image') {
      steps{
        echo "Building Image..."
         ws("${env.WORKSPACE}/application/") {
            sh 'docker build -t $IMAGE_PATH:$IMAGE_TAG .'
        }
      }
    }
    
    stage ('Push Image') {
      steps {
       echo "Pushing Image to Repository..."
        ws("${env.WORKSPACE}/application/") {
          withCredentials([usernamePassword(credentialsId: 'regcred', passwordVariable: 'regpassword', usernameVariable: 'regusername')]) {
            sh 'docker login -u ${regusername} -p ${regpassword}'
            sh 'docker push $IMAGE_PATH:$IMAGE_TAG'
          }
        }
      }
    }

    stage ('Deploy Image to Server') {
      steps {
       echo "Deploying Image to Server..."
        ws("${env.WORKSPACE}/infrastructure/") {
            sh 'kubectl apply -f kubernetes/'
        }
      }
    }
  }
  post {
    success {
      echo "Pipeline ${currentBuild.fullDisplayName}:${env.BUILD_URL} with ${env.BUILD_NUMBER} Completed Successfully"
    }
    failure {
      echo "Pipeline ${currentBuild.fullDisplayName}:${env.BUILD_URL} with ${env.BUILD_NUMBER} Failed to Complete"
    }
  }
}