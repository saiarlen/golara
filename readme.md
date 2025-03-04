## Golara 

### A Laravel kind golang framework on top of gofiber and gorm

- More to add, any contributions are welcome :)
- Supports migrations and auto refresh and many helper functions are included.
- Docker Support enabled
- Default support mysql

### Install

* create .env file
* cp .denv-exmaple.yaml .denv.yaml
* go build -o ekycapp
* ./ekycapp -subcommand migrate

#### Commands
-    ./ekycapp -subcommand migrate


### Note
* If any .env values changed app rebuild needed
* If any .denv values changed app restart needed



#### Technical Note

* Why panic errors or fmt.Errorf, for panic error no need to tell what to return from the function but not best approach and for fmt.Error need to tell error return from to the function. For the code reference check the sms_provider.go and helpers.go files 


### Package info

* **cast package:** Cast is a library to convert between different go types in a consistent and easy way. ex used in email provider to convert `any` into `int`



#### Developed by saiarlen
