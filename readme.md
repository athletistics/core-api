# ⚡Athletistics - Core API

> Bachelor thesis
>
> Nicolas Fabrice Baeriswyl
>
> HEIA-FR, 2023

## 🚩 Introduction

This repository contains the core API as well as the Kubernetes configurations for the Athletistics project. More
details can be found in the main repository's technical report.

## Setup

To start this API, you first need to create a .env file pointing towards an Infisical instance with the following
environment variables.

```dotenv
POSTGRES_USER=[postgres login]
POSTGRES_PASSWORD=[postgres login]
POSTGRES_DB=[postgres database name]
DB_HOST=[postgres DNS name]
ZITADEL_INSTANCE=[Zitadel URL]
ZITADEL_KEY_PATH=[some file name to temporarily store the key file]
ZITADEL_API_KEY=[Content of the key file for an API project]
TOTP_ENDPOINT=[URL of the TOTP service]
ZITADEL_INSTANCE_GRPC=[Zitadel GRPC URI]
ZITADEL_MGMT_KEY_PATH=[some file name to temporarily store the key file}
ZITADEL_MGMT_API_KEY=[Content of the key file for a service user]
TOTP_SECONDS_LIMIT=[time in seconds]
METRICS_USR=[login for the metrics path]
METRICS_PWD=[login for the metrics path]
```

Locally, the following `.env` file is to be created.

```.env
INFISICAL_TOKEN=
INFISICAL_API_URL=https://your-infisical-instance.ch/api
ENV=dev
```

If not existing, create a docker network called `athletistics`.

```shell
docker network create -d bridge athletistics
```

Then, start the compose file by executing the following command at the project's root.

```shell
docker compose up --build
```

## License

```
Athletistics
Copyright (C) 2023 Nicolas Fabrice Baeriswyl

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
```