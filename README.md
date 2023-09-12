# Debugging Go: from zero to Kubernetes

This repo contains a web application that takes a text and creates an image or a GIF with the input text in ASCII Art.

## How to start a debugging session on this application

Here are the steps to reproduce the steps to debug this application as seen in the slides. 
The slides are located at the following points:

- [Here](https://docs.google.com/presentation/d/1rr2giYpo9XfnjfccaVf8vjsxjdgs83ME_ZtrzYiiSIs) for the slides presented at WeAreDevs World Congress 2023.
- [Here](https://docs.google.com/presentation/d/1wzF9MPisfZWi16sRiSt1CwkTbLScNhAEKe-g39zRubM) for the slides presented at ContainerDays 2023

### In your local environment

1. Download `dlv` from <https://github.com/go-delve/delve>
2. Run `dlv debug main.go` from the root of this project

The debugger is started at this point. You can interact with it: start with the `help` command to know which commands to run.

Example of debug commands:
(dlv) break byBlink.go:27
[...]
(dlv) continue

Then on another server, on the root of the project, run `make curl-jeans` and the window with the debugger will have stopped at the breakpoint inside byBlink.go, line 27.

### Inside the docker container

You will need the `docker` cli for this.

1. Read the dockerfile for the application (./deploy/docker/Dockerfile).
2. Read the dockerfile for the debugger (./deploy/docker/debug.Dockerfile)
3. Run `make docker-up` to start the application container
4. Run `make docker-debug` to open a shell on the debugger container
5. Run `ps` inside this shell
6. Take the value of the PID of the application (gif-app)
7. Run `dlv attach` with the PID taken from the step 6

The debugger is started at this point. You can follow the same steps of debugging in the local environment to debug the application in docker.

### Inside the kubernetes cluster

The cluster is a `KinD` cluster, you need to download `kind`. Also you'll need to download `helm` for installing all the needed resources in the cluster.

1. Run `make cluster-up` to start the cluster
2. Run `make app-install` to install the application in the cluster
3. Run `make app-debug-k8s POD=$POD` to start the debugger container via the ephemeral container and open a shell. $POD should be replaced with the name of the pod to debug
4. Run `ps` inside this shell
5. Take the value of the PID of the application (gif-app)
6. Run `dlv attach` with the PID taken from the step 5

The debugger is started at this point. You can follow the same steps of debugging in the local environment except run `make curl-jeans-k8s IP=$IP`; $IP should be replaced with the external IP of the service _gif-app_ that you can get with `kubectl get service gif-app`

If there are issues with the installation of the metallb installation, run `docker inspect network kind` and update the first two parts of the ip address pool in the file `deploy\gif-app\templates\loadbalancer.yaml` and restart with `make app-uninstall app-install` and then from step 3 on.

## How to use the endpoints

For this web application three endpoints can be used:

1. /free: POST request transforming a text into a gif using the parameters passed in the JSON body
2. /byBlink: GET request transforming a text into a blinking gif with blue background and yellow text
3. /goCol: GET request transforming a text into a banner gif with light blue background and white text (Go's colors)

You can provide the text to convert to a GIF via a query parameter named `text`, for example: `curl localhost:8080/free?text=help -d '{"gifType":"alt"}'`. This is valid for all endpoints.

### Specifics for the '/free' endpoint

As `/free` accepts POST requests, it accepts a json body with the following inputs:

```go
struct {
    Delay      int     `json:"delay"`
    Figlet     string  `json:"figlet"`
    FontSize   float64 `json:"fontSize"`
    FontPath   string  `json:"fontPath"`
    GifType    string  `json:"gifType"`
    XPtFactor  float64 `json:"xPtF"`
    YPtFactor  float64 `json:"yPtF"`
    BgColorHex string  `json:"bgHex"`
    FgColorHex string  `json:"fgHex"`
}
```

## Few final notes

1. Breakpoints can be set in delve before starting it. See <https://github.com/go-delve/delve/discussions/3189>
2. When the debugger stops at a breakpoint it halts all the goroutines. In a kubernetes environment when working with probes, e.g. liveness and readiness probe, this is harmful because in the case of the liveness probe as it will not get any response it will cause the restart of the pod and end the current debugging session. This can be overcome with the copy of the pod (see Makefile's `app-debug-k8s-copy` for the how to create a copy of the pod with kubectl debug), but the copy of the pod has no label and receives no traffic to be debugged. See <https://github.com/kubernetes/kubectl/issues/1108> for more information.
