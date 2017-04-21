pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                checkout scm
                sh 'go build'
                echo 'go build'
            }
        }
        stage('Deploy') {
             steps {
                sh 'docker-compose down'
                sh 'docker-compose up -d'
                echo 'Deploying....'
             }
        }
    }
}
