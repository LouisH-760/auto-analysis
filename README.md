# *auto-analysis*
Automated binary / memory dump analysis



## How to RUN the program

### **Prerequirement**

- Install [Docker](https://www.docker.com)
- Install [Node.js](https://nodejs.org/en) for the Web Client
- Install [Go](https://go.dev)
- Download the code

### **Run a sample**

#### **Launch the environment**

To analyse a sample you need to be in the *auto-analyser* folder, and to place the sample within this folder, as docker cannot copy files outside it's **environment** (the folder from which you are launching any docker command).

The following command will create a bare docker environment, when ran from the root of the repository. Run the docker-builder with the `-h` option to get an overview of which additional tools can be bundled in.

```sh 
go run ./docker-build/ -sample="path_to_sample" -modules="./modules" --run
```

Again, this container will only include the modules scripts found in `./modules`, to include the module dependencies we need to specify them in the command with `-module_name`.

Following modules/tools are available for now:

- **volatility2**
- **rizin**
- **yara**
- **clamav**
- **diec**

#### **run the modules**

to run a module, you can use a http request tool like postman, where you will run the following request
```HTTP
POST /run HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "name":"some_name",
    "script":"path to module (eg. debug/moddebug.py)",
    "arguments":"arguments (encoded as base64 if the module you want to use requires it)"
}
```

The response will be a JSON object, containing status information and the output of the module encoded with base64.

