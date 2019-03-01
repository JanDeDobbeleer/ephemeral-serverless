# ephemeral-serverless: automatically delete your old Tweets with AWS Lambda

**ephemeral** is a Twitter timeline grooming program that runs for pretty much free on AWS Lambda. The code is forked from Vicky Lai's excellent [ephemeral](https://github.com/vickylai/ephemeral) Lambda, and extended to use [serverless](https://serverless.com/) for deployments.

You can use ephemeral-serverless to automatically delete all tweets from your timeline that are older than a certain number of hours that you can choose. For instance, you can ensure that your tweets are deleted after one week (168h), or one day (24h).

The program will run once for each execution based on the trigger/schedule you set in [serverless.yml](./serverless.yml). It will delete up to 200 expired tweets (per-request limit set by Twitter's API) each run.

# Twitter API

You will need to [create a new Twitter application and generate API keys](https://apps.twitter.com/). The program assumes the following environment variables are set using [ssm](https://serverless.com/blog/serverless-secrets-api-keys/) (except for `MAX_TWEET_AGE`):

```
TWITTER_CONSUMER_KEY
TWITTER_CONSUMER_SECRET
TWITTER_ACCESS_TOKEN
TWITTER_ACCESS_TOKEN_SECRET
MAX_TWEET_AGE
```

`MAX_TWEET_AGE` expects a value of hours, such as: `MAX_TWEET_AGE = 72h`

Optionally, you can whitelist certain tweets and save them from deletion by setting the variable `WHITELIST` with the tweet's ID as the value. Find the ID as the string of numbers at the end of the tweet's URL, for example:

https://twitter.com/jandedobbeleer/status/`1101422074718011392`

Set one ID to whitelist, or multiple IDs using the separator `:` like so:

```
WHITELIST = 1101422074718011392:1052942396034609152
```

# build.ps1

This handy bash script is included to help you build your function code using PowerShell.

# License
Copyright (C) 2018-2019 Vicky Lai

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
