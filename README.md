# Maven Deployer
Simple program for mass deploying of jar-archive on remote maven repository. 

## Build

For app building you have to install golang tools package in your OS. Check the installation:

```bash
$ go version
go version go1.19.8 linux/amd64
```

Then, clone, build and launch the app's help:
```bash
$ git clone https://github.com/JAkutenshi/mvn-deployer.git
$ cd mvn-deployer
$ go build -o mvn-deployer main.go
$ ./mvn-deployer -h
Usage of ./mvn-deployer:
  -file string
        Required: JSON-file path with JARs' descriptions
  -host string
        Required: GitLab instance's hostname e.g. "gitlab.com"
  -proj string
        Required: GitLab's project ID: number or unique name
  -serv string
        Required: Desired maven server ID in .m2/settings.xml with gitlab-token provided
```

## JAR-archives JSON's format
The app can deploy (upload) all given jar-archives in provided json-file to remote GitLab project's maven repository. JSON's file `example-jars.json` format:

```json
[
   {
        "groupId": "commons-cli",
        "artifactId": "commons-cli",
        "version": "1.2",
        "path": "example-jars/commons-cli-1.2.jar"
    },
   {
        "groupId": "org.apache.commons",
        "artifactId": "commons-math3",
        "version": "3.6.1",
        "path": "example-jars/commons-math3-3.6.1.jar"
    }
]
```

## GitLab repository's server settings
Before launching you should prepare your `~/.m2/settings.xml` and provide maven repository server:

```xml
<server>
  <id>my-server-id</id>
  <configuration>
    <httpHeaders>
      <property>
        <name>Private-Token</name>
        <value>my-secret-gitlab-access-api-token-here</value>
      </property>
    </httpHeaders>
  </configuration>
</server>

```

## Launching example
Then, you can go:
```bash
$ mvn-deployer -file=example-jars.json -host=gitlab.com -proj=1 -serv=my-server-id
2023/08/19 22:34:48 JARs to maven repo uploader starts...
2023/08/19 22:34:48 Successfully load "example-jars.json"
2023/08/19 22:34:48 Unmarshalling json file...
2023/08/19 22:34:48 Unmarshalling json file done! The file is closed
2023/08/19 22:34:48 JARs to maven's repo uploading starts...
2023/08/19 22:34:48 Uploading artifact "commons-cli:commons-cli:1.2" started...
2023/08/19 22:35:00 Uploading artifact "commons-cli:commons-cli:1.2" done!
2023/08/19 22:35:00 Uploading artifact "org.apache.commons:commons-math3:3.6.1" started...
2023/08/19 22:35:06 Uploading artifact "org.apache.commons:commons-math3:3.6.1" done!
2023/08/19 22:35:06 JARs to maven's repo uploading done!
```
