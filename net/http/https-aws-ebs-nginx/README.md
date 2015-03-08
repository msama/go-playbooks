# Redirecting requests to HTTPS using a reverse proxy

In this playbook the GO service is deployed with AWS EBS using Docker. EBS uses nginx as a reverse proxy between incoming requests and Docker. To handle HTTPS requests and to redirect from http to https it is a simple matter of configuring EBS and nginx to do that for us.

The folder `.ebextensions` contains nginx additional configuration files which are applied in alphabetical order.

First enable incoming requests from TSL on port 443. Notice thay this example is configuring an instance insde a VPC. The configuration would be slightly different  if it wasnt.

```
Resources:
  sslSecurityGroupIngress: 
    Type: AWS::EC2::SecurityGroupIngress
    Properties:
      GroupId: {Ref : AWSEBSecurityGroup}
      IpProtocol: tcp
      ToPort: 443
      FromPort: 443
      CidrIp: 0.0.0.0/0
```

Then install a new nginx config file redirecting http to https through Docker:

```
files:
  /etc/nginx/conf.d/ssl.conf:
    mode: "000755"
    owner: root
    group: root
    content: |
      # HTTPS Server
      server {
        listen 80;
        return 301 https://$host$request_uri;
      }
      server {
        listen 443;
        server_name localhost;
        
        ssl on;
        ssl_certificate /etc/pki/tls/certs/server.crt;
        ssl_certificate_key /etc/pki/tls/certs/server.key;
        
        location / {
          proxy_pass http://docker;
          proxy_http_version 1.1;
          
          proxy_set_header Connection "";
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
      }
```

Then also install SSL certificate and private key:

```
files:
  /etc/pki/tls/certs/server.crt:
    mode: "000400"
    owner: root
    group: root
    content: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
      
  /etc/pki/tls/certs/server.key:
    mode: "000400"
    owner: root
    group: root
    content: |
      -----BEGIN RSA PRIVATE KEY-----
      ...
      -----END RSA PRIVATE KEY-----
```

### Concatenating certificates

If you are using a signed certificate it is always a good practice to have your server distributing also the intermediate certificates by concatenating them. Please note that the behaviour is different between browsers and mobile devices. Chrome for intance fetch all the missing intermediate certificates while Android fails with a "Root certificate not found" exception. Also notice that is normally considered a good practice not to include the root certificate because not necessary.

```
cat yourdomain.crt any_intermediate_ca.crt root.crt > ssl-bundle.crt
```
Further reading:
* [Concatenating certificates in NGINX](https://support.comodo.com/index.php?/Default/Knowledgebase/List/Index/37/certificate-installation)
* [COMODO's knowledgebase on certificate installation](https://support.comodo.com/index.php?/Default/Knowledgebase/Article/View/789/37/certificate-installation-nginx)

### testing your installation

It is always a good idea to verify your installation and to doublecheck if your certificate is being correctly served. Online there are a variety of free web tools that can verify the security of your domain. Examples are:
* https://www.ssllabs.com/ssltest/analyze.html
* https://sslcheck.globalsign.com/en_US