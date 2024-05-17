# UniFi DDNS for Cloudflare via AWS Lambda
UniFi DDNS doesn't support Cloudflare in the UI. This repo is a proxy to fwd UniFi DDNS updates to Cloudflare.\
With typical usage, free tier should be plenty. 

## 1. Prepare Cloudflare
1. Find your API key:
   - https://dash.cloudflare.com/profile/api-tokens > Global API Key
2. Find domain zone id:
   - domain's Overview page, right sidebar, `Zone ID` value
3. Find DNS record id:
   - add/Edit an A record with a temporary ip, ie `home.mydomain.com`=`127.0.0.1`
   - navigate cloudflare dashboard > Manage Account > Audit Log
   - expand most recent log, copy the `Resource ID` value

## 2. Get zip package for AWS
1. Download `cloudflare-ddns.zip` from https://github.com/frifox/cloudflare-ddns/releases
2. Or build it yourself:
```shell
git clone github.com/frifox/cloudflare-ddns
cd cloudflare-ddns
GOOS=linux GOARCH=arm64 go build -o bootstrap -tags lambda.norpc -ldflags "-s -w"
zip cloudflare-ddns.zip bootstrap
```

## 3. Create AWS Lambda
1. Goto AWS Lambda (https://console.aws.amazon.com/lambda/home)
2. Click `Create Function`:
   - select `Author from scratch`
   - name your function, ie `cloudflare-ddns`
   - set Runtime `Amazon Linux 2023`
   - set Architecture `arm64`
   - under `Advanced settings`:
     - check `Enable function URL`
     - set Auth Type `NONE`
   - `Create function`
3. Under `Code` tab:
   - `Upload from` > `.zip file` > `Upload`
   - find your `cloudflare-ddns.zip`
   - `Save`
4. Under `Configuration` tab:
   - `Function URL` > copy the url, will be used in UniFi later

## 4. Config your lambda func via ENV
1. Under `Configuration` tab
   - `Environment variables` > `Edit` > `Add environment variable`
   - add Key/Values:
     - `APP_USER` = pick a username for your lambda app
     - `APP_PASS` = pick a password for your lambda app
     - `CF_EMAIL` = login email for cloudflare
     - `CF_KEY` = cloudflare global api key
     - `CF_ZONE_ID` = cloudflare domain's zone id
     - `CF_RECORD_ID` = cloudflare domain's A record id

## 5. Config Unifi
1. Go to https://unifi.ui.com/ > Your console
2. Go to Settings > Internet > Primary (or Secondary)
3. Under `Dynamic DNS`, click `Create New Dynamic DNS`:
   - Service = `custom`
   - Hostname = your FQDN, ie`home.mydomain.com`
   - Username = `APP_USER` value
   - Password = `APP_PASS` value
   - Server = `{your_aws_lambda_function_url.aws}/?hostname=%h&ip=%i`
     - IMPORTANT: remove the `https://` prefix from your aws lambda function url
   - `Save`

## 6. Force DDNS to update now
1. Go to https://unifi.ui.com/ > Your console
2. Go to UniFi Devices > your router/gateway > on right sidebar: Settings > Debug
4. Once you get the shell, find inadyn conf file location & force it to update:
```shell
# ps x | grep inadyn.conf
2668706 ?        S<     0:00 /usr/sbin/inadyn -n -s -C -f /run/ddns-ppp0-inadyn.conf

# inadyn -n -1 --force -f /run/ddns-ppp0-inadyn.conf
inadyn[2908558]: In-a-dyn version 2.9.1 -- Dynamic DNS update client.
inadyn[2908558]: Update forced for alias home.mydomain.com, new IP# 123.123.123.123
inadyn[2908558]: Updating cache for home.mydomain.com
```