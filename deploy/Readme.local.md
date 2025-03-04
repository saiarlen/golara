## Local Deployment steps

* Download `Docker` and install
* git pull the repo.
* Change the git branch. 
* setup `.env` and `.denv.yaml` from the example files.
* Note: We are using docker so db local url in denv is `host.docker.internal`
* Initally run `docker-compose up --build` :: Every time if any docker conif files updated or .mod or .sum files updated then rebuild the app. other changes air will handle the reset or just reset the docker.
* To start the container just run `docker-compose up`
* To exit the container just `ctl^c`
* to stop the container run `docker-compose down`
* For run the terminal migrations or any commands use docker desktop container terminal.



#### Tips:

* Install vscode docker extensions for ease of use.
