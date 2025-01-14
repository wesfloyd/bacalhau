# Running locally with the 'devstack' command

The `devstack` command of bacalhau will start a 3 node cluster alongside isolated ipfs servers.

This is useful to kick the tires and/or developing on the codebase.  It's also the tool used by some of the tests.

## Pre-requisites

 * x86_64 linux host
    * Ubuntu 20.0+ has most often been used for development and testing
    * Note: Mac M1 (ARM64) compatible builds are not yet supported at this time. Please consider development on a hosted alternative, such as [Gitpod](https://gitpod.io/#https://github.com/filecoin-project/bacalhau)
 * Go >= 1.17
 * IPFS v0.11
 * [Docker Engine](https://docs.docker.com/get-docker/)
 * A build of the [latest Bacalhau release](https://github.com/filecoin-project/bacalhau/releases/)


### IPFS Installation Instructions

```
# Install IPFS *v0.11* specifically (due to issues in v0.12) via https://docs.ipfs.io/install/command-line/#official-distributions

wget https://dist.ipfs.io/go-ipfs/v0.11.0/go-ipfs_v0.11.0_linux-amd64.tar.gz
tar -xvzf go-ipfs_v0.11.0_linux-amd64.tar.gz
cd go-ipfs
sudo bash install.sh
cd -
```

## (Optional) Building Bacalhau from source

```
sudo apt-get update && sudo apt-get install -y make gcc zip
sudo snap install go --classic
wget https://github.com/filecoin-project/bacalhau/archive/refs/heads/main.zip
unzip main.zip
cd bacalhau-main
go build

```



## Start the cluster

```bash
make devstack
```

This will start a 3 node bacalhau cluster connected with libp2p.

Each node has it's own ipfs server isolated using the `IPFS_PATH` environment variable and it's own JSON RPC server isolated using a random port.

Once everything has started up - you will see output like the following:

```bash
-------------------------------                                     
environment                                                         
-------------------------------                                     

IPFS_PATH_0=/tmp/bacalhau-ipfs1110685378                            
JSON_PORT_0=41081                                                   
IPFS_PATH_1=/tmp/bacalhau-ipfs919189468                             
JSON_PORT_1=41057                                                   
IPFS_PATH_2=/tmp/bacalhau-ipfs490124113                             
JSON_PORT_2=41347
```

## New Terminal Window
* Open an additional terminal window to be used for data submission to the local IPFS instances and and job submission to the 3 node devestack Bacalhau cluster.
* Copy and paste the IPFS and JSON port variables into the new terminal window.

## Add files to IPFS

Each node has it's own `IPFS_PATH` value which points to a path on the local filesystem.  This allows to use the ipfs cli to test adding files to one or multiple nodes.  This is especially useful when you want to test self selection of a job based on whether the cid is *local* to that node.

To add a file to only one of ipfs node within the devstack cluster, execute the `ipfs add` in the following manner:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
```
*Note: the CID is saved as an environment variable so that it can be referenced in the job submission step.

## Set a json rpc port

Each node has it's own `--jsonrpc-port` value.  This means you can use the `go run .` cli in isolation from the other 2 nodes.

For example - to view the current job list from the perspective of only one of the 3 nodes:

```bash
# Note: replace 12345 this with the correct port from the output
export NODE1_JSONRPC_PORT=12345
go run . --jsonrpc-port=$NODE1_JSONRPC_PORT list
```

## Submit a simple job

This will submit a simple job to a single node:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
go run . --jsonrpc-port=$JSON_PORT_0 submit --cids=$cid --commands="grep kiwi /ipfs/$cid"
go run . --jsonrpc-port=$JSON_PORT_0 list
```

After a short while - the job should be in `complete` state.

```
kai@xwing:~/projects/bacalhau$ go run . --jsonrpc-port=$JSON_PORT_0 list
JOB       COMMAND                  DATA                     NODE                     STATE     STATUS                                                               OUTPUT                                         
63b0a80e  grep kiwi /ipfs/QmRy...  QmRyDNzrxwcL4ENNGyKL...  QmcMKp2NQm2QQf7nRFjK...  complete  Got job results cid: QmRZa9mCrjMgMtaaZZTAEBRdHCVJR3WjoncsEuZU9qBpzv  QmRZa9mCrjMgMtaaZZTAEBRdHCVJR3WjoncsEuZU9qBpzv
```

Copy the job id into a variable:

```bash
JOB_ID=63b0a80e
```

Now we can list the results:

```bash
go run . --jsonrpc-port=$JSON_PORT_0 results list $JOB_ID
```

This will show the following:

```
kai@xwing:~/projects/luke/bacalhau$ go run . --jsonrpc-port=$JSON_PORT_0 results list 63b0a80e
NODE                                            IPFS                                                                 RESULTS                                                                      DIFFERENCE  CORRECT 
QmcMKp2NQm2QQf7nRFjKgTaknsdvmFsp4zjJHvAoP9CvRu  https://ipfs.io/ipfs/QmRZa9mCrjMgMtaaZZTAEBRdHCVJR3WjoncsEuZU9qBpzv  ~/.bacalhau/results/63b0a80e/QmcMKp2NQm2QQf7nRFjKgTaknsdvmFsp4zjJHvAoP9CvRu           0  ✅      
```

The results from the job are stored in the `~/.bacalhau/results/<JOB_ID>/<NODE_ID>` directory.

We can see the files that were output by the job here:

```bash
ls -la ~/.bacalhau/results/63b0a80e/QmcMKp2NQm2QQf7nRFjKgTaknsdvmFsp4zjJHvAoP9CvRu
```

## run 3 node job

Now let's run a job across all 3 nodes.  To do this, we need to add the cid to all the IPFS servers so the job will be selected to run across all 3 nodes:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
IPFS_PATH=$IPFS_PATH_1 ipfs add -q ./testdata/grep_file.txt
IPFS_PATH=$IPFS_PATH_2 ipfs add -q ./testdata/grep_file.txt
```

Then we submit the job but with `--concurrency` and `--confidence` settings:

```bash
go run . --jsonrpc-port=$JSON_PORT_0 submit --cids=$cid --commands="grep pear /ipfs/$cid" --concurrency=3 --confidence=2
go run . --jsonrpc-port=$JSON_PORT_0 list
```

We can see that all 3 nodes have produced results by getting the job id and running:

```bash
go run . --jsonrpc-port=$JSON_PORT_0 results list <JOB_ID>
```

## run 3 node job with bad actor

Now let's restart the devstack but this time with one of the three nodes in `bad actor` mode.  This bad node will not run the job and instead just sleep for 10 seconds.

ctrl+c on the running dev-stack and re-run with:

```bash
make devstack-badactor
```

Copy and paste the environment section into your other terminal and then let's submit another job to the 3 nodes:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
IPFS_PATH=$IPFS_PATH_1 ipfs add -q ./testdata/grep_file.txt
IPFS_PATH=$IPFS_PATH_2 ipfs add -q ./testdata/grep_file.txt
go run . --jsonrpc-port=$JSON_PORT_0 submit --cids=$cid --commands="grep pear /ipfs/$cid" --concurrency=3 --confidence=2
go run . --jsonrpc-port=$JSON_PORT_0 list
```

This time - when you list the results, you will see that our bad actor has been caught because their memory trace is substantially different from the others.
