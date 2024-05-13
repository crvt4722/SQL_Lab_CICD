pipeline {
    agent any
    environment {
        DOCKER_REGISTRY = 'index.docker.io'
        DOCKER_CREDENTIALS_ID = 'dckr_pat_eh56hEXLPuNVQv7q09Lg7vOCROQ'
        REPO_NAME = 'crvt4722'
        IMAGE_NAME = 'sql_lab_server'
    }
    stages {
        stage('Build and Push Image') {
            steps {
                script {
                    def dockerImage = docker.build("${REPO_NAME}/${IMAGE_NAME}:${env.GIT_TAG}")
                    docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
                        dockerImage.push()
                    }
                }
            }
        }
    }
}

