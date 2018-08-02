This is a short procedure to build genapp from sources:

  ```
  export GOPATH=~/go					# Set the go path
  ```
  ```
  mkdir ~/go/src && cd ~/go/src				# Create src folder and enter in it
  ```
  ```
  git clone git@github.com:forj-oss/goforjj.git		# Clone goforjj repository
  ```
  ```
  cd /goforjj/genapp					# Enter in genapp folder
  ```
  ```
  source ./build-env.sh					# Source environment
  ```
  ```
  build.sh						# Build
  ```
  ```
  go install						# Now you have genapp build on ~/go/bin
  ```
