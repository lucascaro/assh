# assh

assh is a command line tool for ssh-ing into all boxes in an auto scaling group.


## Install

```
go install github.com/lucascaro/assh
```

## Usage

run with `-h` for usage instructions.

```
assh -h
```
## Examples

```
# ssh into all instances in the ASG my-asg-name using the ec2-user user and
# the ssh key in ~/.ssh/id-rsa.pem
assh my-asg-name

# Specify a user name
assh -u my-user my-asg-name

# Specify a key file
assh -i ~/.ssh/my-other-key.pem my-asg-name
```
