# Public key accounts

`publickey` user API database allows you to use a [Beacon](https://www.walletbeacon.io/) compatible public key authentication scheme for machine to machine communication. It's not a real database and no data is stored anywhere. So you can't change any account properties.

To use it simply set `user_api.account_database.connection_string` to `publickey` in the Dendrite config file.