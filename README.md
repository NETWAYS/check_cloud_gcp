# check_cloud_gcp

Icinga check plugin to check Google Cloud Platform (GCP) resources. At the moment the check only supports the
Compute Engine (GCE) context.

## Usage

### Computing - Instances

When one of the states is non-ok, or a instance is STOPPED, the check will alert.

#### compute instances

Checks all GCP Instances over all zones or multiple GCP Instances in a defined GCP zone.

```
Usage:
  check_cloud_google compute instances [flags]

Flags:
  -f, --filter string        Filter expression that filters resources e.g. '(cpuPlatform = "Intel Broadwell") AND (name != "instance1")'
  -h, --help                 help for instances
      --ignore-api-warning   Disables warning when querying without a zone filter (please do not do that)
  -z, --zone string          GCP Zone name, can include wildcards (e.g. "europe-*")

Global Flags:
  -j, --json-file string   GCP service account key file
```

```
## Zone: europe-west3-c

$ check_cloud_gcp compute instances -z europe-west3-c -j $GOOGLE_APPLICATION_CREDENTIALS
CRITICAL - 2 Instances found - 1 RUNNING - 1 TERMINATED

[OK] "instance1" powerstate=RUNNING size=e2-micro
[CRITICAL] "instance2" powerstate=TERMINATED size=e2-medium


## Zone: all Zones

$ check_cloud_gcp compute instances --filter 'name != "instance-1"' -j $GOOGLE_APPLICATION_CREDENTIALS
CRITICAL - 3 Instances found - 2 RUNNING - 1 TERMINATED

## us-central1-a

[OK] "instance-3" powerstate=RUNNING size=e2-micro

## europe-west3-c

[CRITICAL] "instance-2" powerstate=TERMINATED size=e2-medium
```

More information on [filters](https://pkg.go.dev/google.golang.org/api@v0.43.0/compute/v1#ZoneOperationsListCall.Filter)

#### compute instance

Checks a single GCP Instance

```
check_cloud_gcp compute instance -z europe-west3-c -j $GOOGLE_APPLICATION_CREDENTIALS -n instance1
OK - "instance1" powerstate=RUNNING size=e2-micro
```

## Setting up Access

In order to work correctly you need the correct permissions and configuration within GCP, to grant the plugin proper
read-only access to the resources.

The following step-by-step instructions will help you to setup this configuration.

### Creating a service account

You should create a new service account within the cloud project, and add the proper permissions to it,
name it e.g. "check_cloud_gcp".

* In the Cloud Console, go to the Service accounts page.
* Go to the Service accounts page
* Select a project.
* Click Create service account.
* Enter a service account name to display in the Cloud Console e.g. `check_cloud_gcp`
* Optional: Enter a description of the service account e.g. `monitoring purposes`
* Choose a IAM role to grant the service account the correct permissions on the project `Basic -> Viewer`
* When you are done adding roles, click Continue.
* Optional: In the Service account users role field, add members that can impersonate the service account.
* Optional: In the Service account admins role field, add members that can manage the service account.
* Click Done to finish creating the service account.

### Creating a Service account key

The check itself needs a service account key file which will parse the credentials of the service account:

* In the Cloud Console, go to the Service Accounts page.
* Select a project.
* Click the email address of the service account that you want to create a key for.
* Click the Keys tab.
* Click the Add key drop-down menu, then select Create new key.
* Select JSON as the Key type and click Create.

**Important:** The key can be downloaded only after creation. Also restrict the permissions of the credential file
on disk, so only Icinga can read it.

Then either set environment `GOOGLE_APPLICATION_CREDENTIALS` or pass `--json-file` with the file path.

For more information about authentication see: [Getting started with authentication](https://cloud.google.com/docs/authentication/getting-started)

The full documentation can be found at [Google Cloud Documentation](https://cloud.google.com/docs/)

## License

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
