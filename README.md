# vke-agent

`vke-agent` is a simple command-line tool used for setting up Kubernetes clusters. With this tool, you can quickly provision both master and worker nodes.

## Getting Started

Follow these steps to get started with using `vke-agent`.

### Prerequisites

The following prerequisites are required for the project to be used:

- `curl` command
- A user account with `sudo` privileges

### Installation

Follow these steps to build and run the project:

1. Clone this repository:

   ```bash
   git clone https://github.com/vmindtech/vke-agent
   cd vke-agent
  

2. Build the project
   ```bash
    env GOOS=linux GOARCH=amd64 go build vke-agent.go
    ```
3. Use the vke-agent command to create Kubernetes nodes:
   ```bash
    env GOOS=linux GOARCH=amd64 go build vke-agent.go
    ```
   You can use the above command to create the master node. For worker nodes, you can use the following commands:
      ```bash
    ./vke-agent -agentialize=false -rke2AgentType="agent" -rke2Token="your-token" -serverAddress="https://your-adress:9345" -kubeversion="v1.28.2+rke2r1"  -tlsSan="your-loadbalancer-adress"

    ```

# License
This project is licensed under the MIT License - see the LICENSE.md file for details.

