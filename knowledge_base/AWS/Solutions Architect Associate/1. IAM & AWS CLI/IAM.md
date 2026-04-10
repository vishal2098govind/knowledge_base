#iam #aws #solutions-architect-udemy 

### Users, Groups and Policies 
Every AWS Account when created, a root account is created by default

This root account has administrator level access and thus shouldn't be used or shared with anybody else

The only thing it should be used for is to setup initial AWS Account that's it

Instead create users in IAM
One **user** represents one **person** within the organization
Users can be grouped together if it makes sense

Example:
- 6 people organization
![[Pasted image 20260407235316.png]]
- Groups can only contain users, not other groups
- Some users need not belong to any group
- A user can belong to multiple groups

### IAM Permissions
- Users or Groups can be assigned JSON documents called **IAM Policy**
- These policies define **permissions** of users
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "ec2:Describe*",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "elasticloadbalancing:Describe*",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "cloudwatch:ListMetrics",
                "cloudwatch:GetMetricStatistics",
                "cloudwatch:Describe*"
            ],
            "Resource": "*"
        }
    ]
}
```
- **Least privilege principle** : Don't give more permissions than a user needs

>IAM as an entire service is a **Global** Service


## IAM Policies
![[Pasted image 20260408000708.png]]
### IAM Policy Structure
![[Pasted image 20260408000838.png]]
- Components:
	- Version: version of policy language
	- ID (Optional)
	- Statement array (Required)
		- **Sid** - optional statement id
		- **Effect** - whether to allow or deny listed actions
		- **Principal** - AWS Account or user or role or service to which this policy is applied 
		- **Action** - List of actions or permissions this policy allows or denies
		- **Resource** - List of resources to which the actions are applied to
		- **Condition** - List of conditions when this policy is in effect (optional)

| Principal Type     | Example                               |
| ------------------ | ------------------------------------- |
| AWS Account (root) | `arn:aws:iam::123456789:root`         |
| IAM User           | `arn:aws:iam::123456789:user/vishal`  |
| IAM Role           | `arn:aws:iam::123456789:role/my-role` |
| AWS Service        | `lambda.amazonaws.com`                |
| Federated identity | via SAML / OIDC / Cognito             |
| Everyone           | `*`                                   |
**Principal only appears in resource-based policies** (like S3 bucket policies, KMS key policies, trust policies on roles).
**Identity-based policies** (attached to users/roles) don't have a Principal field because the principal is _implicit_ — it's whoever the policy is attached to.

Examples:
```json
{
  "Version": "2012-10-17",
  "Id": "default",
  "Statement": [
    {
      "Sid": "lambda-605bb3db-2723-4483-b83d-97ed0bb29461",
      "Effect": "Allow",
      "Principal": {
        "Service": "s3.amazonaws.com"
      },
      "Action": "lambda:InvokeFunction",
      "Resource": "arn:aws:lambda:us-east-1:103001974:function:sagemaker-callback",
      "Condition": {
        "StringEquals": {
          "AWS:SourceAccount": "1030041974"
        },
        "ArnLike": {
          "AWS:SourceArn": "arn:aws:s3:::ams-data-platform-ml-data-prod"
        }
      }
    },
    {
      "Sid": "103002841974_event_permissions_from_ams-data-platform-ml-data-prod_for_sagemaker-callback",
      "Effect": "Allow",
      "Principal": {
        "Service": "s3.amazonaws.com"
      },
      "Action": "lambda:InvokeFunction",
      "Resource": "arn:aws:lambda:us-east-1:1030041974:function:sagemaker-callback",
      "Condition": {
        "StringEquals": {
          "AWS:SourceAccount": "1030041974"
        },
        "ArnLike": {
          "AWS:SourceArn": "arn:aws:s3:::ams-data-platform-ml-data-prod"
        }
      }
    }
  ]
}
```

## IAM MFA
- MFA = password + security device owned
- Main benefit of MFA:
	- If a password is stolen or hacked, the account is not compromised
- MFA Device Options
	- Virtual MFA device
		- Google Authenticator (phone only)
		- Authy (phone only)
	- Universal 2nd Factor (U2F) **security key**
		- 3rd Party device
	- **Hardware** key Fob MFA device - By 3rd Party
	- Hardware Key Fob MFA Device for AWS GovCloud (US) - By 3rd Party

## How humans and applications access AWS
- **AWS Management Console**
    - Protected by password + MFA
- **AWS CLI**
    - Protected by access key + secret
    - Signs requests via SigV4 under the hood
- **AWS SDK** (for application code)
    - Protected by access key + secret, or preferably an **IAM Role** when running on AWS compute (EC2, Lambda, ECS)
    - When a role is attached, the SDK calls **AWS STS** via the **instance metadata service (IMDS)** automatically to get short-lived credentials, no manual AssumeRole needed
    - SDK resolves credentials via the **credential provider chain** (env vars → credentials file → IMDS role → ...), so the same code works locally and on AWS without changes
    - Roles are safer than access keys, credentials are ephemeral and auto-rotated
- **Direct AWS API**
    - Everything above is ultimately a wrapper around raw **HTTPS calls to AWS API endpoints**
    - Requests are signed using **Signature Version 4 (SigV4)**
    - Useful to know when making raw calls via curl or Postman

**IAM Roles — Where can they be attached?**

Any AWS service that needs to _do something on your behalf_ gets an IAM Role. The question to always ask is, **"what does this service need to do, and does it need permission to do it?"** If yes, it gets a role.

**Under the hood, this always involves STS.** When a service assumes a role, it calls **STS AssumeRole** internally, which returns temporary credentials — an access key, secret, and a **session token** — valid for a short duration. The service then uses these to sign its API requests. You never see this happen, AWS handles it transparently.

So the flow is always:

```
Service → STS AssumeRole → temporary credentials → signed API calls
```

| Category                   | Services                                                                                                                             |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **Compute**                | EC2 (via Instance Profile)<br>Lambda<br>ECS (Task Role)<br>EKS (via IRSA)<br>ECS Fargate<br>Elastic Beanstalk<br>App Runner<br>Batch |
| **Orchestration**          | Step Functions<br>EventBridge<br>MWAA (Managed Airflow)                                                                              |
| **Data & Analytics**       | Glue<br>EMR<br>Athena<br>Kinesis Firehose<br>Redshift<br>QuickSight                                                                  |
| **ML**                     | SageMaker (execution role)                                                                                                           |
| **Storage & Integration**  | S3 Replication<br>Lambda@Edge                                                                                                        |
| **Developer & Automation** | CodeBuild<br>CodePipeline<br>CodeDeploy<br>CloudFormation                                                                            |
| **Other**                  | API Gateway<br>IoT Core<br>Systems Manager (SSM)                                                                                     |


### IAM Permission, IAM Policy, IAM Role, IAM User, IAM User Groups

```
IAM Permission (actions based on service) - low level service level actions/apis

IAM Policy - collection of IAM permissions or actions or apis along with 
	- Principal (optional)
		  - principal is implicit in case of identity based policy (IAM Role or IAM User)
		  - principal is explicitly added in case of resource based policies
		  - can be: 
		    - User arn 
		    - or Role arn 
		    - or Service name
			- or AWS account
	- Resources
		  - 
	- Conditions (optional)
		  - can be:
		    - 
```

### IAM Security Tools
- IAM Credentials Report (account level)
	- list of account's users and status of their various credentials
- IAM Access Advisor (user level)
	- shows service permissions granted to a user and when those services were last accessed
	- can use this to revise policies by reducing permissions not being used and follow **principal of least privilege**

### IAM Best Practices
- Don’t use the root account except for AWS account setup
- One physical user = One AWS user
- Assign users to groups and assign permissions to groups
- Create a strong password policy
- Use and enforce the use of Multi Factor Authentication (MFA)
- Create and use Roles for giving permissions to AWS services
- Use Access Keys for Programmatic Access (CLI / SDK)
- Audit permissions of your account using IAM Credentials Report & IAM Access Advisor
- Never share IAM users & Access Keys