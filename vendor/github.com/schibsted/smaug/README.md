# Smaug

Smaug is a credentials provider for [mesos2iam](https://github.com/schibsted/mesos2iam/)

**Build**

```
make build
```

##### **Run**

```
./smaug --credentials-repository-file /tmp/my-roles.ini
```

## Roles definition

You can define roles using the following .ini file:

```
[roles]
ca82d854-6bc2-4f50-ba0c-8bfbb24cb1ef = arn:aws:iam::my-aws-account:role/testSmaug
```
