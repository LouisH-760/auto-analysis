# Modules

Modules are the main way of automating analysis. They are uploaded and ran within the docker container used for analysis. Due to them being able to run code without much in the way of controls, they should be reviewed for safety before being used.
## Module metadata

Module metadata is described in a {json} format, including a name, short description, context of executions and the type of expected inputs as well as outputs. They can also describe setup commands to be ran within the docker container
### Inputs / outputs
## Module logic

Module logic should be contained within a single python script

### Execution context
### Available tools

# Run process

The main program responsible for controlling the setup, teardown and execution of scripts is split in two parts: a control interface ran on the user machine, and an implant in the docker container exposing a port to allow control of the different script modules and setup tasks

## External interface

## Docker implant

## communication

Communication is achieved using JSON payloads in HTTP requests. Once the docker environment is ready, the workflow can be started. The control interface sends a request to the implant telling it which script to start and giving it the required parameters, which are passed over to the script using STDIN. Once execution completes, the implant sends a request to the UI with the results of the script obtained from STDOUT and STDERR. Data is once again passed in a JSON format that the script or ui can parse to extract results.