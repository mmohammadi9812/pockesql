# PockeSQL
---

## STATUS
I see this project as completed, and it probably won't receive new features

So, consider this project in _maintanance_ mode, where it may or may not receive bug fixes,

but no new features

## Why
[Pocket](https://getpocket.com/), while being a great and (thanks god) free service, is still a third-party service

That means if someday, they decide to charge for their service, or even close (hope not!),

what happens with your data is at best, a rush to download before it's always gone

That's why I, along with a lot of other people, believe that owning data is important

## What

This project helps to own your pocket account data by retrieving your saved articles to sqlite

you may run `update` command time to time to fetch new saved articles or (soft-)delete removed items

run `pockesql --help` to see other commands and flags

P.S. I have not searched about limits of _Pocket_ applications _yet_, so that's why I'm not comfortable sharing my key

It's assumed that you have created an application in [develop portal](https://getpocket.com/developer/apps/) and consumer key is present in env file

copy `.env.example` file to `.env` file and copy credentials as required

### Credit

This tool is a heavily inspired by [pocket-to-sqlite](https://github.com/dogsheep/pocket-to-sqlite) by [Simon Willison](https://simonwillison.net/), to whom I'm very thankful

## License
Copyright 2022 Mohammad Mohamamdi. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

