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
                    echo 'Pull new code'
		    sh(script: 'git pull')
                }
                script {
                    def tagVersion = sh(script: 'git tag --sort version:refname | tail -1', returnStdout: true).trim()
                    env.TAG_NAME = tagVersion
                    echo "Tag version: ${env.TAG_NAME}"
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

