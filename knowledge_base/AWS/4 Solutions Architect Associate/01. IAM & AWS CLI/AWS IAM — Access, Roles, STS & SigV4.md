#iam #claude-notes
## How Humans and Applications Access AWS

- **AWS Management Console** — protected by password + MFA
- **AWS CLI** — protected by access key + secret, signs requests via SigV4 under the hood
- **AWS SDK** (for application code) — protected by access key + secret, or preferably an IAM Role when running on AWS compute. SDK resolves credentials via the **credential provider chain** (env vars → credentials file → IMDS / container endpoint → ...), so the same code works locally and on AWS without changes
- **Direct AWS API** — everything above is ultimately a wrapper around raw HTTPS calls to AWS API endpoints, signed using **Signature Version 4 (SigV4)**. Useful to know when making raw calls via curl or Postman

---

## Principal

The entity that is **authenticated and making a request** in AWS. In a policy, the principal specifies **who is allowed or denied** to perform actions.

|Principal Type|Example|
|---|---|
|AWS Account (root)|`arn:aws:iam::123456789:root`|
|IAM User|`arn:aws:iam::123456789:user/vishal`|
|IAM Role|`arn:aws:iam::123456789:role/my-role`|
|AWS Service|`lambda.amazonaws.com`|
|Federated identity|via SAML / OIDC / Cognito|
|Everyone|`*`|

> Principal only appears in **resource-based policies** (S3 bucket policies, KMS key policies, role trust policies). Identity-based policies don't have a Principal field — the principal is implicit, it's whoever the policy is attached to.

---

## IAM Roles — Where Can They Be Attached?

Any AWS service that needs to **do something on your behalf** gets an IAM Role. The question to always ask is, **"what does this service need to do, and does it need permission to do it?"** If yes, it gets a role.

|Category|Services|
|---|---|
|**Compute**|EC2 (via Instance Profile), Lambda, ECS (Task Role), EKS (via IRSA), Fargate, Elastic Beanstalk, App Runner, Batch|
|**Orchestration**|Step Functions, EventBridge, MWAA (Managed Airflow)|
|**Data & Analytics**|Glue, EMR, Athena, Kinesis Firehose, Redshift, QuickSight|
|**ML**|SageMaker (execution role)|
|**Storage & Integration**|S3 Replication, Lambda@Edge|
|**Developer & Automation**|CodeBuild, CodePipeline, CodeDeploy, CloudFormation|
|**Other**|API Gateway, IoT Core, Systems Manager (SSM)|

> Every service role ultimately involves an STS call somewhere in the chain. What differs is **who makes that call** — the platform, the metadata service, or your code.

---

## IAM Role — Two Policies

Every IAM Role has two policies:

**Trust Policy** — who is allowed to assume this role:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

**Permission Policy** — what this role is allowed to do:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:PutObject"],
      "Resource": "arn:aws:s3:::my-bucket/*"
    }
  ]
}
```

---

## Trust Policy Conditions

Conditions narrow down who can assume a role. All conditions within a single block are **ANDed** — all must be true.

**1. Restrict to a specific account**

```json
"Condition": {
  "StringEquals": {
    "aws:PrincipalAccount": "111122223333"
  }
}
```

**2. Restrict to a specific IAM role/user ARN**

```json
"Condition": {
  "ArnLike": {
    "aws:PrincipalArn": "arn:aws:iam::111122223333:role/deployment-role"
  }
}
```

**3. Require MFA + recency**

```json
"Condition": {
  "Bool": { "aws:MultiFactorAuthPresent": "true" },
  "NumericLessThan": { "aws:MultiFactorAuthAge": "3600" }
}
```

**4. Restrict by source IP**

```json
"Condition": {
  "IpAddress": {
    "aws:SourceIp": ["203.0.113.0/24", "198.51.100.0/24"]
  }
}
```

**5. Restrict to specific VPC / VPC Endpoint**

```json
"Condition": {
  "StringEquals": {
    "aws:SourceVpc": "vpc-0a1b2c3d4e5f",
    "aws:SourceVpce": "vpce-0a1b2c3d4e5f6a7b8"
  }
}
```

**6. ExternalId — confused deputy prevention**

Critical for third-party integrations. Without the correct `ExternalId`, the AssumeRole call is rejected even if the principal is correct:

```json
"Condition": {
  "StringEquals": {
    "sts:ExternalId": "my-secret-external-id-xyz987"
  }
}
```

Vendor must pass this in their AssumeRole call:

```json
{
  "RoleArn": "arn:aws:iam::123456789012:role/vendor-role",
  "RoleSessionName": "vendor-session",
  "ExternalId": "my-secret-external-id-xyz987"
}
```

**7. Restrict max session duration**

```json
"Condition": {
  "NumericLessThanEquals": {
    "sts:DurationSeconds": "900"
  }
}
```

**8. ABAC — restrict by principal tags**

```json
"Condition": {
  "StringEquals": {
    "aws:PrincipalTag/Team": "backend",
    "aws:PrincipalTag/Environment": "production"
  }
}
```

**9. Restrict by region**

```json
"Condition": {
  "StringEquals": {
    "sts:RequestedRegion": "us-east-1"
  }
}
```

|Condition Key|Purpose|
|---|---|
|`aws:PrincipalAccount`|Restrict to specific account|
|`aws:PrincipalArn`|Restrict to specific user/role ARN|
|`aws:MultiFactorAuthPresent`|Require MFA|
|`aws:MultiFactorAuthAge`|Require recent MFA|
|`aws:SourceIp`|Restrict by IP range|
|`aws:SourceVpc`|Restrict to specific VPC|
|`aws:SourceVpce`|Restrict to specific VPC Endpoint|
|`sts:ExternalId`|Confused deputy prevention for third parties|
|`sts:DurationSeconds`|Restrict max session duration|
|`aws:PrincipalTag/*`|ABAC — restrict by principal's tags|
|`sts:RequestedRegion`|Restrict which region the call comes from|

---

## Complete Flow: IAM Role → STS → Signed API Call

### Step 1: Role is attached to the service

Role with trust policy + permission policy attached to Lambda:

```bash
aws lambda create-function \
  --function-name my-function \
  --role arn:aws:iam::123456789012:role/my-lambda-role
```

### Step 2: Lambda is invoked, AWS runtime calls STS AssumeRole

```http
POST https://sts.amazonaws.com/ HTTP/1.1

Action=AssumeRole
&RoleArn=arn:aws:iam::123456789012:role/my-lambda-role
&RoleSessionName=my-function-execution-abc123
&DurationSeconds=3600
```

- `RoleSessionName` is auto-generated, useful for CloudTrail audit logs
- `DurationSeconds` defaults to 1 hour for Lambda

### Step 3: STS validates and returns temporary credentials

```json
{
  "Credentials": {
    "AccessKeyId": "ASIAXXXXXXXXXXX12345",
    "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
    "SessionToken": "FwoGZXIvYXdzEJr...long-token...Tuw==",
    "Expiration": "2026-04-08T11:00:00Z"
  }
}
```

- `AccessKeyId` starting with `ASIA` → always means temporary STS credentials. Permanent IAM user keys start with `AKIA`
- `SessionToken` encodes the role, expiry, and session metadata
- After `Expiration`, credentials are completely useless even if intercepted

### Step 4: SDK picks up credentials via credential provider chain

```
1. Hard-coded credentials in code?          → no
2. AWS_ACCESS_KEY_ID env var?               → no
3. ~/.aws/credentials file?                 → no
4. ECS container credentials endpoint?      → no
5. Lambda runtime environment variables?    → YES ✓
   AWS_ACCESS_KEY_ID     = ASIAXXXXXXXXXXX12345
   AWS_SECRET_ACCESS_KEY = wJalrXUtnFEMI/...
   AWS_SESSION_TOKEN     = FwoGZXIvYXdzEJr...
```

> The STS call already happened before your code started running — made by the AWS runtime, not your SDK. By the time the SDK walks the credential provider chain, credentials are already sitting in env vars. The SDK's job here is purely **credential discovery**, not credential generation.

### Step 5: SDK builds the canonical request

```
GET
/my-file.txt

host:my-bucket.s3.amazonaws.com
x-amz-content-sha256:e3b0c44298fc1c149afb...
x-amz-date:20260408T100000Z
x-amz-security-token:FwoGZXIvYXdzEJr...

host;x-amz-content-sha256;x-amz-date;x-amz-security-token

e3b0c44298fc1c149afb...
```

The request is called **canonical** because it is a strict normalization of the raw HTTP request into a **single, unambiguous, authoritative form**. The same HTTP request could be expressed in multiple equivalent ways (different header casing, different query string order, trailing slashes), and HMAC is extremely sensitive — even one character difference produces a completely different signature.

Normalization rules:

- HTTP method → always uppercase (`GET` not `get`)
- URI path → always normalized (no `..`, no double slashes)
- Query string → sorted alphabetically (`a=1&b=2` not `b=2&a=1`)
- Header names → always lowercase (`host` not `Host`)
- Headers → sorted alphabetically
- Header values → whitespace trimmed
- Body → always SHA256 hashed, never included raw

### Step 6: SDK signs via SigV4

```
# Hash the canonical request
SHA256(canonical_request) = "3b4c5d6e..."

# Build string to sign
AWS4-HMAC-SHA256
20260408T100000Z
20260408/us-east-1/s3/aws4_request
3b4c5d6e...

# Derive signing key (scoped per date/region/service)
kDate    = HMAC-SHA256("AWS4" + SecretAccessKey, "20260408")
kRegion  = HMAC-SHA256(kDate,    "us-east-1")
kService = HMAC-SHA256(kRegion,  "s3")
kSigning = HMAC-SHA256(kService, "aws4_request")

# Final signature
Signature = HMAC-SHA256(kSigning, string_to_sign)
```

Attach to Authorization header:

```http
Authorization: AWS4-HMAC-SHA256
  Credential=ASIAXXXXXXXXXXX12345/20260408/us-east-1/s3/aws4_request,
  SignedHeaders=host;x-amz-content-sha256;x-amz-date;x-amz-security-token,
  Signature=fe5f80f77d5fa3beca...
```

The signing key is derived fresh per date/region/service — a key valid for S3 today cannot sign a DynamoDB request tomorrow.

### Step 7: S3 validates the request

1. Extracts `AccessKeyId` from Authorization header. `ASIA` prefix → temporary credentials, expects a session token
2. Calls STS internally to validate `SessionToken` — is it valid, non-expired, which role does it belong to?
3. Independently recomputes the canonical request and signature, compares with the one in the header
4. Checks IAM permissions — does `my-lambda-role` have `s3:GetObject` on this bucket?
5. Any explicit Deny? Any bucket policy blocking this?
6. All checks pass → response returned

### Step 8: Credentials expire, cycle repeats

- **Lambda**: AWS runtime automatically calls STS again before expiry, injects fresh credentials
- **EC2**: IMDS at `http://169.254.169.254/latest/meta-data/iam/security-credentials/<role>` handles rotation, SDK polls it automatically

---

## Who Makes the STS Call — by Service

|Pattern|Services|Who calls STS|
|---|---|---|
|**Pre-injected credentials**|Lambda, ECS, Fargate, App Runner|AWS platform, before your code starts. Credentials in env vars / container endpoint|
|**On-demand via IMDS**|EC2|IMDS calls STS and caches result. SDK calls IMDS, never STS directly|
|**SDK calls STS directly**|EKS (IRSA)|SDK calls `AssumeRoleWithWebIdentity`, exchanges Kubernetes service account JWT for AWS credentials|
|**You call STS explicitly**|Cross-account, Federation, CLI|You call `AssumeRole` manually with RoleArn, ExternalId, SessionName etc.|

---

## Full Flow Summary

```
Attach Role to Lambda (trust policy + permission policy)
          ↓
Lambda invoked
          ↓
AWS runtime → POST STS AssumeRole (RoleArn, RoleSessionName)
          ↓
STS checks trust policy → returns (ASIA... AccessKeyId, SecretAccessKey, SessionToken, Expiration)
          ↓
Credentials injected into Lambda runtime environment
          ↓
SDK walks credential provider chain → finds credentials in env vars
          ↓
SDK builds canonical request (method + URI + headers + SHA256(body))
          ↓
SDK builds string to sign (algorithm + timestamp + credential scope + hash of canonical request)
          ↓
SDK derives signing key (HMAC chain: secret → date → region → service → aws4_request)
          ↓
SDK computes Signature = HMAC-SHA256(signingKey, stringToSign)
          ↓
SDK attaches Authorization header + x-amz-security-token to request
          ↓
Request hits S3
          ↓
S3 extracts AccessKeyId → calls STS internally to validate session token
          ↓
S3 recomputes signature independently → compares with request signature
          ↓
S3 checks IAM permissions for the assumed role
          ↓
All checks pass → response returned
          ↓
Credentials near expiry → STS called again → fresh credentials injected
          ↓
Cycle repeats
```

---

## Cross-Account Access — STS Extended

### The Problem

Two AWS accounts are completely isolated by default:

```
Account A (123456789012) — your application, Lambda runs here
Account B (999988887777) — your data, S3 bucket lives here
```

Account A's roles have zero access to anything in Account B out of the box.

### The Solution — a role in Account B that Account A is allowed to assume

**Step 1: In Account B, create a role with a trust policy pointing at Account A**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:role/my-lambda-role"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "shared-secret-xyz"
        }
      }
    }
  ]
}
```

**Step 2: In Account B, attach a permission policy to that role**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject"],
      "Resource": "arn:aws:s3:::account-b-bucket/*"
    }
  ]
}
```

**Step 3: In Account A, give `my-lambda-role` permission to call AssumeRole on Account B's role**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Resource": "arn:aws:iam::999988887777:role/cross-account-role"
    }
  ]
}
```

> Both sides must agree — Account B's role must trust Account A, AND Account A's role must be allowed to call AssumeRole on Account B's role. Missing either side = access denied.

### The Full Cross-Account Flow

```
Lambda in Account A starts
        ↓
AWS runtime calls STS in Account A → gets Layer 1 credentials for my-lambda-role
        ↓
Your code explicitly calls STS AssumeRole for Account B's role:

stsClient.AssumeRole({
  RoleArn:         "arn:aws:iam::999988887777:role/cross-account-role",
  RoleSessionName: "lambda-cross-account-session",
  ExternalId:      "shared-secret-xyz"
})
        ↓
STS checks:
  - Does Account B's trust policy allow my-lambda-role from Account A? → yes
  - Does Account A's permission policy allow calling AssumeRole on this ARN? → yes
  - Does ExternalId match? → yes
        ↓
STS returns a NEW set of temporary credentials scoped to Account B's role:
  AccessKeyId:     ASIAYYYY...
  SecretAccessKey: newSecret...
  SessionToken:    newToken...
        ↓
Your code uses THESE credentials to build a new S3 client pointed at Account B
        ↓
S3 request signed with Account B credentials
        ↓
S3 in Account B validates:
  - Signature valid? → yes
  - SessionToken valid, belongs to cross-account-role in Account B? → yes
  - Does cross-account-role have s3:GetObject on this bucket? → yes
        ↓
Response returned to Lambda in Account A
```

### Credential Layering

In cross-account, your Lambda holds **two sets of credentials simultaneously:**

```
Layer 1 — Account A credentials (from runtime, via env vars)
  → used for anything in Account A
  → used to call STS AssumeRole for Account B

Layer 2 — Account B credentials (from explicit AssumeRole call)
  → used only for resources in Account B
  → completely separate STS session
  → separate expiry, separate rotation
```

The mechanism is identical to single-account — AssumeRole → temporary credentials → signed requests. The only difference is you do it twice, and the second time you do it explicitly in code rather than letting the runtime handle it.

### Why ExternalId Matters More in Cross-Account

In single-account, if someone tricks a service into assuming a role, blast radius is limited to your account. In cross-account, a third-party vendor has a role in your account they can assume from their account. If their system gets compromised, an attacker could assume your role from the vendor's account.

`ExternalId` is the defence — even if an attacker gains access to the vendor's AWS account, they don't know your `ExternalId`, so the AssumeRole call fails. This is the **confused deputy problem** and cross-account is exactly where it matters most.

### Common Real-World Patterns

|Pattern|Example|
|---|---|
|Multi-account org|Dev, staging, prod in separate accounts, CI/CD assumes role in each|
|Third-party access|Datadog, New Relic assuming a role in your account to read metrics|
|Shared services account|Central logging or security account that other accounts push to|
|Data mesh|Each team owns their data account, other teams assume roles to query|

---

## IAM Policy Language Reference

|What you need|Where to look|
|---|---|
|Full policy grammar (elements, structure)|https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies.html|
|All `aws:*` global condition keys|https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_condition-keys.html|
|Per-service Actions, Resource ARNs, Condition keys|https://docs.aws.amazon.com/service-authorization/latest/reference/reference_policies_actions-resources-contextkeys.html|
|All condition operators (`StringEquals`, `ArnLike`, etc.)|https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html|
|Formal BNF-style policy grammar|https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html|
|Test policies without deploying|https://policysim.aws.amazon.com/|
