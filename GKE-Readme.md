
## Getting Started with GKE

Here's how to create and install pre-requisites for the Layer7 project

### It's best to run this on a fresh GKE (Google Kubernetes Engine) cluster

  - You can usually run this using the free credits you'll receive when you open a new account.

### Pre-requisites
* Google SDK
[Installation documentation](https://cloud.google.com/sdk/install)

* A GKE Cluster
   - Google SDK command to create a working context
   
   ```code
   $ gcloud init
   .... a bunch of questions and confirmations ....
   $ gcloud config list
  [compute]
  region = us-west1
  zone = us-west1
  [core]
  account = <This.is.you@yourdomain.com>
  disable_usage_reporting = False
  project = my_project

  Your active configuration is: [default]
    $ gcloud container clusters create l7-project --machine-type=n1-standard-4 --num-nodes=2
    ... a bunch of stuff .... then
    kubeconfig entry generated for l7-project.
    NAME        LOCATION  MASTER_VERSION  MASTER_IP      MACHINE_TYPE   NODE_VERSION    NUM_NODES  STATUS
    l7_project  us-west1  1.14.10-gke.24  104.198.x.y    n1-standard-4  1.14.10-gke.24  9          RUNNING

   ```

 - A Cluster Role Binding 
   - For this demo, we are using the cluster-admin role,
   - Locally for production, obviously you would restrict this as per your requirements
   ```code
   $ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user [your-registered-email@yourdomain]
   ```