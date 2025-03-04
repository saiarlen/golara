## Deploy Process

### Requirements

* Ubuntu 24.04
* Nginx
* Mysql
* go 1.22
* supervisor
* inotify-tools



### Process

* Using Git FTP copy the build file to app server root folder (Note: Only build file is enough to run the app)
* copy the `.denv-example.yaml` to `.denv.yaml` and setup the values
* .env is optional because it auto embed during build in github

#### Step-2

* Provide permissions to the build: `sudo chmod +x ekycapp-randomval`
* Create a symlink: `sudo ln -s ekycapp-randomval ekycapp.symlink`
* Make a copy or create a `storage` folder as per exp-storage structure and provide appropriate permissions.
* Create a `watch_ekycapp.sh` with the reference of exp_watch_ekycapp.sh and provide permissions `sudo chmod +x watch_ekycapp.sh`
* Edit watch_ekycapp.sh file and update `WORK_DIR` path as your app root path.

#### Step-3

* Creat Supervisor files as per examples in deploy folder and modify the paths inside the file as of root path

* #### Steps to create supervisor files
  1. go to `cd /etc/supervisor/conf.d`
  2. Create anyname_app.conf and watch_ekycapp.conf as per examples 
  3. nano the files modify the paths and username
  4. then 
  `sudo supervisorctl reread`
  `sudo supervisorctl update`
  `sudo supervisorctl start program_name`




* #### Dependinces 
 1. ghostscript
 2. imagemagick
 2. poppler-utils
 4. wkhtmltopdf => follow separate installation process below





 #### `Wkhtmltopdf` installation (NO NEED THIS BCZ BINARY ADDED IN xbin FOLDER )
 * Direct installitaion wont enable the qt patch for header and footer so manual install needed
 * `wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.jammy_amd64.deb`
 * `dpkg -i wkhtmltox_0.12.6.1-2.jammy_amd64.deb`
 * if dependance error `sudo apt --fix-broken install` or ` apt-get -f install -y`
 * if this error "root@user:/home/path# wkhtmltopdf --version
bash: /usr/bin/wkhtmltopdf: No such file or directory"
THEN Follow below
* sudo find / -name wkhtmltopdf 2>/dev/null
If the binary is found, note the directory (e.g., /usr/local/bin/wkhtmltopdf).
"Temporarily add it to your current session:"
`export PATH=$PATH:/usr/local/bin/`
`wkhtmltopdf --version`
"If it works, make the change permanent by adding the PATH update to your shell configuration file:"
`echo 'export PATH=$PATH:/usr/local/bin/' >> ~/.bashrc`
`source ~/.bashrc`
