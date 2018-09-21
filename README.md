# Surtr

Surtr is a cron job that will terminate the oldest node in your kubernetes cluster.
Having really old nodes in your cluster isn't ideal, it can make you vulnerable to malicious attacks and can sometimes hide failure states in your applications.
By making sure no node ever lives too long you both remove the potential security risk and expose flaws in your cluster/application design, if taking out one node causes a load of errors you know you need to improve the resilience of your services.

# Usage
Surtr is intended to be ran as cron-job, see the [examples](examples) directory for the cron and rbac definitions.
Surtr needs to be able to query the Kubernetes api for a list of the nodes, it finds the oldest node, checks it's older than the min age specified and then sends an ec2 termination request for that node.

# AWS Permissions
Surtr just needs `ec2:TerminateInstances`

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "",
            "Action": [
                "ec2:TerminateInstances"
            ],
            "Effect": "Allow",
            "Resource": [
                "*"
            ]
        }
    ]
}

```

# Flags
```
--help                   Show context-sensitive help (also try --help-long and --help-man).
--kubeconfig=KUBECONFIG  Path to kubeconfig.
--older-than=OLDER-THAN  age of nodes to terminate
--debug                  Debug mode
 ```
