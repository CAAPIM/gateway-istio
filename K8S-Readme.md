
```FIXME: @JT correct doc```
```FIXME: @GV add steps for minikube, will involve metallb to use lbs```

# Local or Small Installation on a minikube system E.G. on a laptop

* You will need to have minikube installed. Kubernetes [Installation page](https://kubernetes.io/docs/tasks/tools/install-minikube/)

* This requires Docker as well. [Choose your OS](https://docs.docker.com/install/)

* ``` FIXME: This Step needs to source prebuilt containers```

* #### Local Setup for Mac & Linux Based hosts
    * For the initial setup, go to the `minikube` folder and type `make create`. 
    * This will create (by default) a minikube cluster that reserves 10G of RAM and 4 CPUs.
    * You can type `make` by itself to see your options. 
    * This will take several minutes


