# Documentation for Orchestrating the Infrastructure and Deploying the Application

#### Provisioning the Infrastructure
---
I orchestrated a virtual machine with an ubuntu operating system for this task. Please install python3, ansible, vagrant and virtualbox within your terminal. 

<b> Steps </b>
  - <code> ssh-keygen -t rsa -b 4096 -C "sd-automation key"</code> and type "/home/samuel/.ssh/sd-automation_key" on prompt to create an ssh key pair.
  - <code> vagrant up --provision </code> # to provision and configure the vm with all the requiste technologies for the task.
  1. 

#### Import the CSV file into MongoDB
---
<b> Steps </b>
  - The data should be imported automatically on configuration of mongodb if not then run the next command.
  - Run `mongoimport --db boarding --collection people --file titanic.csv --type csv --headerline` please ensure the file is in current directory
  - Log into the mongo shell with `mongo` from your terminal. Run `db.people.find({})` to display the imported data of all persons on board or `db.people.findOne({"_id" : ObjectId("5e1d694701abfe84ba65edb2")})`  to display details of a particular person. When done type  `exit` to quit the mongo shell and return to the terminal.

#### Deploying the Application
---
<b> Steps </b>
  - cd directory into "kubernetes" and run `kubectl apply -f .` to run all kubernetes manifest file




