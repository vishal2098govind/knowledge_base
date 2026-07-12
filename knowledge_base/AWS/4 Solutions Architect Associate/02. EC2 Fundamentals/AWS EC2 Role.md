#aws-ec2 #solutions-architect-udemy #iam-role

### Configuring an EC2 instance with AWS CLI `aws configure` to gain permissions to use AWS services

**THIS IS A REALLY REALLY REALLY BAD IDEA**

Because, if we run aws configure and enter our personal access key id and secret access key on an ec2 instance, then anyone else in our aws account can SSH into that instance and simply run commands using AWS CLI and use the AWS services they want to and worst they can retrieve the values of this secret

Instead, attach EC2 instance with an IAM Role with required permissions or policies