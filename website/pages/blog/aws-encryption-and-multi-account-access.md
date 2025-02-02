---
title: Encryption in AWS and Multi-Account Access
tag: security
date: 2022/09/29
description: >-
  How to encrypt in AWS given a multi-account environment.
author: jsonkao
---

import { BlogHeader } from "../../components/BlogHeader"

<BlogHeader/>

## *Dance like nobody is watching. Encrypt like everyone is.*

 - *Werner Vogels, Amazon CTO*

As AWS [outlined at their 2022 Re:Inforce Security Conference](https://www.youtube.com/watch?v=PPunA7tPMyk&t=3062s) and [mentioned by Werner Vogels at an AWS Summit in 2019](https://youtu.be/vWfkbGF6fiA?t=4339), *encrypt everything* and *encrypt like everyone is [watching]*.  

![**AWS Summit Series 2019 - Santa Clara: Keynote featuring Werner Vogels**](/images/blog/aws-encryption-and-multi-account-access/encrypt-like-watching.png)
*AWS Summit Series 2019 - Santa Clara: Keynote featuring Werner Vogels*

In this blog post, we’ll focus on how to encrypt everything in multi-account AWS environments and how to make encryption decisions with your unique environment and data security needs in mind.

In follow-up blog posts of this series on encryption, we’ll follow up with information on how CloudQuery can help, as well as further posts deep-diving into encryption and data security in the cloud.

## Multi-Account Access and Encryption in AWS

For symmetric encryption, AWS offers [2 primary services](https://docs.aws.amazon.com/crypto/latest/userguide/awscryp-service-toplevel.html): **[AWS Key Management Service (KMS)](https://docs.aws.amazon.com/crypto/latest/userguide/awscryp-service-toplevel.html) and [CloudHSM](https://aws.amazon.com/cloudhsm/)**. When AWS KMS was [first announced in 2014](https://aws.amazon.com/blogs/aws/new-key-management-service/), it was launched to support encrypting data at rest for S3, EBS, and Redshift. Now, KMS supports multiple different types of keys including symmetric encryption keys, asymmetric keys for encryption or signing, and HMAC keys to generate and verify HMAC tags. [KMS now supports many more services](https://aws.amazon.com/kms/features/#AWS_Service_Integration).  Some of those services are available for direct access from other AWS accounts, such as S3, SQS, Secrets Manager, and more.

AWS provides a table to describe AWS KMS keys and information about how they can be used and their features as shown below. 

![[AWS Table for Customer Keys and AWS Keys](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html)](/images/blog/aws-encryption-and-multi-account-access/aws-kms-table.png)
*[AWS Table for Customer Keys and AWS Keys](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html)*

We’re going to enrich that table with more detail regarding multi-account access and management while expanding the table with more information about different types of customer managed keys with different key material origins.  We’ve differentiated between the 3 types of Customer Managed Keys: Keys with External Key material, Keys backed by a Custom Key Store (CloudHSM), and lastly, keys with AWS-Provided Key material.  

**Expanded Table:**

| Key Type | Multi-account Access | Can view metadata | Can manage KMS Key | Used Only for my AWS account | Automatic Rotation | Pricing |
| --- | --- | --- | --- | --- | --- | --- |
| Customer Managed Key: External Key Material | Yes | Yes | Yes | Yes (1)  | No | Monthly & Per-Use Fee |
| Customer Managed Key: Custom Key Store | Yes | Yes | Yes | Yes (1) | No | Monthly & Per-Use Fee |
| Customer Managed Key: AWS-Provided Key Material | Yes | Yes | Yes | Yes | Optional | Monthly & Per-Use Fee |
| AWS Managed Key | No | Yes | No | Yes | Required | Per-use fee |
| AWS Owned Key | Varies | No | No | No | Varies | Varies |

### Example Scenario for Multi-Account Access in AWS

In advanced use cases, enterprise cloud workloads may be split up by infrastructure or by project.  For example, one account may host data such as a data lake account and another account may host compute resources.  In these multi-account AWS environments, cross-account access to resources can be necessary to reduce complexity and to reduce the need for data and resource duplication.

With cross-account access to resources, encryption also plays a role in how cross-account access can be granted to users and applications originating from a different account.  The type of KMS Key chosen can affect how cross-account setup can be done and in some cases, make it more complex to manage.

We’ll walk through a setup where cross-account access may be desired.

![Cross-Account Access in AWS to an Encrypted S3 Bucket](/images/blog/aws-encryption-and-multi-account-access/cross-account-diagram.png)
*Cross-Account Access in AWS to an Encrypted S3 Bucket*

As an example, we will look at an IAM role in a Compute AWS Account that needs access to an encrypted S3 bucket in the Data AWS Account.  One method to grant access is to configure the following components if we’re using a customer managed KMS Key:

- IAM Role in the Compute AWS Account with corresponding policies that grant access to the S3 Bucket, Objects, and KMS Key in the Data AWS Account.
- S3 Bucket with a bucket policy in the Data AWS Account that grants access to the IAM Role from the Compute AWS Account.
- KMS Key with a Key Policy in the Data AWS Account that grants access to the IAM Role from the Compute AWS Account.

Let’s revisit the encryption key table that we expanded upon earlier in this post with the additional column for multi-account access.  What happens when we try to use an AWS-managed KMS Key for this cross-account use case?  Note - this is different than the Amazon S3-managed key (SSE-S3) option, which if there’s interest, we’ll do a deep dive into the encryption and security settings available for AWS S3 including S3-managed keys, bucket keys, and more.

![AWS Managed key: aws/s3](/images/blog/aws-encryption-and-multi-account-access/aws-managed-key-s3.png)
*AWS Managed key: aws/s3*

We’ll modify the encryption setting for the S3 bucket in the example data AWS Account to use the aws/s3 AWS managed KMS key and then upload a new object which will use that aws/s3 encryption setting and encrypt the object with the aws/s3 AWS managed KMS key.

![AWS Default Encryption Settings for S3 Bucket](/images/blog/aws-encryption-and-multi-account-access/s3-default-encryption.png)
*AWS Default Encryption Settings for S3 Bucket*

Now, after ensuring that we’ve properly configured the IAM policies attached to the IAM role in the example compute AWS account and the S3 bucket policy attached to the bucket in the data AWS account to permit for access for the IAM role to the S3 bucket and objects, we’ll try to retrieve this new object encrypted with the aws/s3 AWS managed key.  We get an AccessDenied error with the additional context that either the key doesn’t exist or we’re not allowed to access the key.

![AccessDenied error for cross-account S3 GetObject](/images/blog/aws-encryption-and-multi-account-access/not-allowed-to-access-key.png)
*AccessDenied error for cross-account S3 GetObject*

In this case, due to the key being an AWS Managed KMS Key, we’re unable to access the KMS Key due to the following KMS Key Policy.

```json
{
    "Version": "2012-10-17",
    "Id": "auto-s3-2",
    "Statement": [
        {
            "Sid": "Allow access through S3 for all principals in the account that are authorized to use S3",
            "Effect": "Allow",
            "Principal": {
                "AWS": "*"
            },
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt",
                "kms:ReEncrypt*",
                "kms:GenerateDataKey*",
                "kms:DescribeKey"
            ],
            "Resource": "*",
            "Condition": {
                "StringEquals": {
                    "kms:CallerAccount": "123412341234",
                    "kms:ViaService": "s3.us-east-1.amazonaws.com"
                }
            }
        },
        {
            "Sid": "Allow direct access to key metadata to the account",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::123412341234:root"
            },
            "Action": [
                "kms:Describe*",
                "kms:Get*",
                "kms:List*"
            ],
            "Resource": "*"
        }
    ]
}
```

Let’s take a deeper look at this key policy and what this statement does.  Since they’re both Allow statements, these 2 statements function as a logical **or** statement for access.  This Key Policy has 2 main statements:

- “Allow access through S3 for all principals in the account that are authorized to use S3”

```json
{
            "Sid": "Allow access through S3 for all principals in the account that are authorized to use S3",
            "Effect": "Allow",
            "Principal": {
                "AWS": "*"
            },
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt",
                "kms:ReEncrypt*",
                "kms:GenerateDataKey*",
                "kms:DescribeKey"
            ],
            "Resource": "*",
            "Condition": {
                "StringEquals": {
                    "kms:CallerAccount": "123412341234",
                    "kms:ViaService": "s3.us-east-1.amazonaws.com"
                }
            }
        },
```

There are 2 main conditions on this statement that limit the access to this key, despite the “AWS”: “*” Principal block.

- “kms:CallerAccount”: “123412341234”
- “kms:ViaService”: “s3.us-east-1.amazonaws.com”

By combining the kms:CallerAccount condition with a Principal element that specifies all AWS identities, this statement specifies all identities in an AWS account, in this case the account 123412341234.  The kms:ViaService then limits the usage of an AWS KMS Key to requests from a specified AWS service, in this case “s3.us-east-1.amazonaws.com.”  Thus, this key can only be used by all identities in the 123412341234 account with requests from “s3.us-east-1.amazonaws.com.”

- “Allow direct access to key metadata to the account”

```json
{
            "Sid": "Allow direct access to key metadata to the account",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::123412341234:root"
            },
            "Action": [
                "kms:Describe*",
                "kms:Get*",
                "kms:List*"
            ],
            "Resource": "*"
        }
```

With the Principal in this statement being “AWS”: “arn:aws:iam::123412341234:root”, the account principal, the statement doesn’t give any IAM users or roles permissions to use the KMS key.  Instead, this allows for delegation of permissions.  Thus, if an IAM role or user has the above permissions to Describe, Get, and List the key, the IAM role or user can Describe,Get, and List the key.  This is because key policies require explicit permissions on the key policy (different from other resource policies) to grant permissions.  This allows for the key metadata access to be managed via IAM policies.

These 2 statements combined do not permit for our cross-account access example, which is a request that originates outside of the 123412341234 account.  Thus, our request is denied access to the KMS key as shown in the AccessDenied error message.  Furthermore, we’re unable to modify the key policy for AWS provided Managed KMS Keys.

Due to the nature and setup of the AWS provided AWS Managed KMS Keys, which are created, managed, and used on the customer’s behalf by an [AWS service integrated by AWS KMS](https://aws.amazon.com/kms/features/#AWS_Service_Integration), we do not recommend using these keys for resources used in cross-account workloads.  Customers are unable to change the properties of AWS managed keys, create KMS key grants, rotate them, change their key policies, or schedule them for deletion.  Direct access to the AWS Managed KMS Key via grants or KMS Key policies from outside the account hosting the KMS key is not possible.  Instead, indirect access such as cross-account role assumption into the account would need to be used. 

### KMS Key Access

Access for both usage and management of KMS Keys can be governed by a couple different mechanisms.  At a high level, there are IAM Policies, KMS Key Grants, and KMS Key Policies (Resource-Based Policies) where a combination of them could grant access.  We will do a deep dive into these concepts and complexities in a later post.

```json
{
    "Id": "cloudquery-sample-cmk-policy",
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Enable IAM User Permissions",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::123412341234:root"
            },
            "Action": "kms:*",
            "Resource": "*"
        },
        {
            "Sid": "Allow use of the key",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::123412341234:role/cloudquery-role"
            },
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt",
                "kms:ReEncrypt*",
                "kms:GenerateDataKey*",
                "kms:DescribeKey"
            ],
            "Resource": "*"
        },
        {
            "Sid": "Allow attachment of persistent resources",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::123412341234:role/cloudquery-role"
            },
            "Action": [
                "kms:CreateGrant",
                "kms:ListGrants",
                "kms:RevokeGrant"
            ],
            "Resource": "*",
            "Condition": {
                "Bool": {
                    "kms:GrantIsForAWSResource": "true"
                }
            }
        }
    ]
}
```

For the [Sample KMS Key Policy](https://docs.aws.amazon.com/kms/latest/developerguide/key-policy-default.html) above, there are decisions made by AWS to reduce the risk of the key becoming unmanageable, add specific management abilities to IAM entities within the account, and gives the AWS account full access to this KMS key.  In a later post, we’ll cover specific recommendations and techniques to balance security and manageability of KMS Keys and policies.

Despite KMS Key Policies being a specific type of resource policy, KMS Key policies function slightly differently than the typical resource policy.  For typical resource policies, Identity-based policies and resource-based policies are both permissions policies and are evaluated together within a single AWS account.  With KMS, the KMS key policy **must** grant access for access to work, even if the corresponding permissions are on identity policies.

With Cross-Account Access, there are a couple of mechanisms to grant access to a KMS Key and the encrypted resources.

- Direct IAM Access via IAM Policies such as a KMS Key Policy and Identity-based Policies.
- Indirect IAM Access via Services (Deputized Access)
- IAM access via IAM entities within the same account (role assumption or other access).

In our cross-account example above, if an AWS Managed KMS Key is used, cross-account access would have to be granted via IAM Access via other IAM entities via cross-account role assumption and not direct access via the KMS key policy.

## Conclusion

- For resources that may be shared across multiple accounts, use Customer Managed KMS Keys.  For most use cases, CloudQuery recommends using AWS-provided key material as AWS KMS supports automatic key rotation for symmetric encryption KMS keys with key material that AWS KMS creates.
- For resources that require encryption with FIPS 140-2 Level 3, has specific audit requirements, or cannot be stored in a shared environment, CloudQuery recommends either using CloudHSM directly or AWS KMS with a custom key store backed by AWS CloudHSM clusters.

We will follow this post shortly with more encryption blog posts in the series.  Up next will be a post explaining how CloudQuery can help determine the encryption and data security posture of your cloud environments.  We’ll publish that post shortly after the release of CloudQuery v1.

If you have comments, feedback on this post, follow-up topics you’d like to see, or would like to talk about CloudQuery or cloud security - email us at security@cloudquery.io or come chat with us on [Discord](https://www.cloudquery.io/discord)!

## References

[AWS Cryptographic Services and Tools](https://docs.aws.amazon.com/crypto/latest/userguide/awscryp-service-toplevel.html)

[AWS Key Management Service Developer Guide: AWS KMS Concepts](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html)

[AWS Key Management Service: How AWS Services use AWS KMS](https://docs.aws.amazon.com/kms/latest/developerguide/service-integration.html)

[AWS White paper: Organizing your AWS Environment Using Multiple Accounts](https://docs.aws.amazon.com/whitepapers/latest/organizing-your-aws-environment/organizing-your-aws-environment.html)

[AWS Key Management Service Developer Guide: Condition Keys for AWS KMS](https://docs.aws.amazon.com/kms/latest/developerguide/policy-conditions.html)