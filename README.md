# *auto-analysis*
Automated binary / memory dump analysis



## How To RUN

### **Prerequirement**

- Install [Docker](https://www.docker.com)
- Install [Node.js](https://nodejs.org/en) for the Web Client
- Install [Go](https://go.dev)
- Download the code

### **Run a sample**

#### **Launch the environment**

To analyse a sample you need to be in the *auto-analyser* folder, and to copy the sample inside this folder because docker cannot read the path off parent folders.

the following command will create the docker environment with all the tools needed

```sh 
go run ./docker-build/ -sample="path_to_sample" -modules="./modules" --run
```

but it will only include the modules scripts found in `./modules`, to include the modules installation we need to specify the modules in the command with `-module_name`.

Following modules/tools are possible to install:

- **volatility2**
- **rizin**
- **yara**
- **clamav**
- **diec**

#### **run the modules**

to run a module you need to use a http request tool like postman,...
where you will run the following request
```HTTP
POST /run HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "name":"some_name",
    "script":"path to module (debug/moddebug.py)",
    "arguments":"arguments in base64"
}
```

The responce will be the result of the command in base64.

