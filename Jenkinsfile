pipeline {
    agent any
    environment {
        DOCKER_CREDENTIALS_ID = 'docker_access_token'  // ID of your Docker Hub credentials in Jenkins
        REPO_NAME = 'crvt4722'  // Docker Hub repository name
        IMAGE_NAME = 'sql_lab_server'  // Image name
    }
    stages {
        stage('Set Tag Name') {
            steps {
                script {
                    if (env.BRANCH_NAME && env.BRANCH_NAME.startsWith('refs/tags/')) {
                        env.TAG_NAME = env.BRANCH_NAME.replace('refs/tags/', '')
                    }
                }
            }
        }
        stage('Build and Push Image') {
            steps {
                script {
                    if (env.TAG_NAME) {
                        // Build the Docker image
                        def builtImage = docker.build("${REPO_NAME}/${IMAGE_NAME}:${env.TAG_NAME}")
                        
                        // Login to Docker Hub
                        docker.withRegistry('https://index.docker.io/v1/', DOCKER_CREDENTIALS_ID) {
                            // Push the Docker image
                            builtImage.push()
                        }
                    } else {
                        echo "No tag found, not building or pushing an image."
                    }
                }
            }
        }
    }
}

