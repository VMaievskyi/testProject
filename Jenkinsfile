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
    }
}
