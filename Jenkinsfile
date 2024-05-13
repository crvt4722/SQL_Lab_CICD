pipeline {
    agent any
    environment {
        DOCKER_CREDENTIALS_ID = 'docker_access_token'  // ID of your Docker Hub credentials in Jenkins
        REPO_NAME = 'crvt4722'  // Docker Hub repository name
        IMAGE_NAME = 'sql_lab_server'  // Image name
    }
    stages {
        stage('Prepare Environment') {
            steps {
                script {
  //                  echo 'Pull new code'
//                    sh 'git pull'
                    def tagVersion = sh(script: 'git describe --tags $(git rev-list --tags --max-count=1)', returnStdout: true).trim()
                    if (tagVersion) {
                        env.TAG_NAME = tagVersion
                        echo "Latest tag version: ${env.TAG_NAME}"
                    } else {
                        error("No tags found in the repository.")
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
                        echo "Image pushed: ${REPO_NAME}/${IMAGE_NAME}:${env.TAG_NAME}"
                    } else {
                        echo "No tag specified, not building or pushing an image."
                    }
                }
            }
        }
    }
    post {
        always {
            echo 'Cleaning up workspace'
            cleanWs()
        }
    }
}

