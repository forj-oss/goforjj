pipeline {
    agent any

    stages {
        stage('Build goforjj') {
            steps {
                withEnv(["DOCKER_JENKINS_HOME=${env.DOCKER_JENKINS_MOUNT}"]) {
                    sh('''set +x ; source ./build-env.sh
                    build.sh''')
                }
            }
        }
        stage('Build genapp') {
            steps {
                withEnv(["DOCKER_JENKINS_HOME=${env.DOCKER_JENKINS_MOUNT}"]) {
                    sh('''cd genapp ; set +x ; source ./build-env.sh
                    build.sh''')
                }
            }
        }
        stage('Tests') {
            steps {
                sh('''set +x ; source ./build-env.sh
                go test goforjj goforjj/genapp''')
            }
        }
    }

    post {
        success {
            deleteDir()
        }
    }
}
