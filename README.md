# classeviva

Classeviva is a Go library and CLI tool to access the popular school portal https://web.spaggiari.eu.

## OpenAPI Specification

An [OpenAPI specification](openapi.yaml) is available that documents the Classeviva API endpoints. This specification can be used to:

- **View interactive API documentation** - Use 3rd party tools like [Swagger Editor](https://editor.swagger.io/) or [Redocly](https://redocly.com/)
- **Generate API clients** in your preferred programming language using tools like [OpenAPI Generator](https://openapi-generator.tech/)
- **Understand API endpoints** and request/response formats

The specification is based on the unofficial community documentation at [Classeviva-Official-Endpoints](https://github.com/Lioydiano/Classeviva-Official-Endpoints).

**Note**: This is an unofficial specification as the Classeviva API is not publicly documented by Spaggiari.

## Quick Start

### Install

```sh
brew install zmoog/homebrew-classeviva/classeviva
```

```sh
brew info zmoog/homebrew-classeviva/classeviva
==> classeviva: 0.2.1

Installed
/opt/homebrew/Caskroom/classeviva/0.2.1 (6.8MB)
  Installed on 2025-10-30 at 23:37:52
From: https://github.com/zmoog/homebrew-classeviva/blob/HEAD/Casks/classeviva.rb
==> Name
classeviva
==> Description

==> Artifacts
classeviva (Binary)
```

### Add your kids

```sh
$ classeviva profile add joelmiller
Username: joelmiller
Password: ******
Profile 'joelmiller' added successfully.
```

```sh
$ classeviva profile list
+------------+------------+---------+
| PROFILE    | USERNAME   | DEFAULT |
+------------+------------+---------+
| joelmiller | joelmiller | *       |
+------------+------------+---------+

$ classeviva grades list
+------------+-------+-----------------------------------------------------------+----------------------------------+
| DATE       | GRADE | SUBJECT                                                   | NOTES                            |
+------------+-------+-----------------------------------------------------------+----------------------------------+
| 2025-10-24 | 9     | SCIENZE MOTORIE E SPORTIVE                                |                                  |
|            | 7-    | STORIA                                                    |                                  |
| 2025-10-22 | 8½    | SCIENZE NATURALI (BIOLOGIA, CHIMICA, SCIENZE DELLA TERRA) | Verifiche genetica moderna e DNA |
+------------+-------+-----------------------------------------------------------+----------------------------------+%     
```

## Authentication

Classeviva supports multiple authentication methods with a priority chain, making it easy to manage credentials for multiple students.

### Authentication Priority

Credentials are resolved in the following order (highest priority first):

1. **CLI flags**: `--username` and `--password`
2. **Profile-based**: `--profile` flag or `default_profile` from config file
3. **Environment variables**: `CLASSEVIVA_USERNAME` and `CLASSEVIVA_PASSWORD`

### Profile Management (Recommended)

For families with multiple students, profile-based authentication is the recommended approach.

#### Setup Profiles

Add profiles for each student:

```shell
# Add first student profile
$ classeviva profile add older-kid
Username: student1_username
Password: ********

# Add second student profile
$ classeviva profile add younger-kid
Username: student2_username
Password: ********
```

#### Set Default Profile

Set a default profile to use when `--profile` is not specified:

```shell
$ classeviva profile set-default older-kid
```

#### List Profiles

View all configured profiles:

```shell
$ classeviva profile list
PROFILE      DEFAULT
older-kid    *
younger-kid
```

#### Show Profile Details

Display profile information (credentials are hidden):

```shell
$ classeviva profile show older-kid
Profile: older-kid
Username: student1_username
Password: ********
```

#### Remove Profile

Delete a profile:

```shell
$ classeviva profile remove younger-kid
```

### Using Profiles

Once profiles are configured, use them with any command:

```shell
# Use default profile
$ classeviva grades list --limit 10

# Use specific profile
$ classeviva --profile younger-kid grades list --limit 10

# Override with CLI flags (highest priority)
$ classeviva --username user --password pass grades list
```

### Configuration File

Profiles are stored in `~/.classeviva/config.yaml`:

```yaml
profiles:
  older-kid:
    username: "student1_username"
    password: "student1_password"
  younger-kid:
    username: "student2_username"
    password: "student2_password"
default_profile: "older-kid"
```

**Note**: The config file has restrictive permissions (0600) to protect credentials.

### Environment Variables (Legacy)

For backward compatibility, you can still use environment variables:

```shell
export CLASSEVIVA_USERNAME="your_username"
export CLASSEVIVA_PASSWORD="your_password"
$ classeviva grades list
```

### Identity Caching

Authentication tokens are cached per-profile in `~/.classeviva/identity-{profile}.json` to minimize API calls. Tokens are automatically refreshed when expired.

## CLI Commands

### Version

Display the application version:

```shell
$ classeviva version
Classeviva CLI v0.0.0 (123) 2022-05-08 by zmoog
```

JSON output:

```shell
$ classeviva version --format json
{
  "version": "v0.0.0",
  "commit": "123",
  "date": "2022-05-08",
  "builtBy": "zmoog"
}
```

### Grades

List student grades with optional limit:

```text
$ classeviva grades list --limit 3
+------------+-------+-----------------+-------------------------------+
| DATE       | GRADE | SUBJECT         | NOTES                         |
+------------+-------+-----------------+-------------------------------+
| 2022-04-27 | 9     | ARTE E IMMAGINE |                               |
| 2022-04-22 | 7+    | COMPORTAMENTO   | comportamento della settimana |
|            | 7     | SCIENZE         |                               |
+------------+-------+-----------------+-------------------------------+
```

JSON output:

```shell
$ classeviva grades list --limit 1 --format json
[
  {
    "subjectDesc": "ARTE E IMMAGINE",
    "evtDate": "2022-04-27",
    "decimalValue": 9,
    "displayValue": "9",
    "color": "green",
    "skillValueDesc": " "
  }
]
```

### Agenda

List agenda items (homework, events) with optional date range and limit:

```text
$ classeviva agenda list --limit 2
+---------------------------+---------------------------+---------+----------------------+-----------------------------------------------------------------------+
| BEGIN                     | END                       | SUBJECT | TEACHER              | NOTES                                                                 |
+---------------------------+---------------------------+---------+----------------------+-----------------------------------------------------------------------+
| 2022-05-02T09:00:00+02:00 | 2022-05-02T10:00:00+02:00 |         | PESANDO MARGHERITA   | Inizio interrogazioni di inglese (1º turno)                           |
| 2022-05-03T00:00:00+02:00 | 2022-05-03T23:59:59+02:00 |         | AVANZATO PAOLA CARLA | Link per colloqui                                                     |
+---------------------------+---------------------------+---------+----------------------+-----------------------------------------------------------------------+
```

JSON output:

```shell
$ classeviva agenda list --until 2022-04-27 --limit 2 --format json
[
  {
    "evtId": 546249,
    "evtCode": "AGNT",
    "evtDatetimeBegin": "2022-05-02T09:00:00+02:00",
    "evtDatetimeEnd": "2022-05-02T10:00:00+02:00",
    "notes": "Inizio interrogazioni di inglese (1º turno)",
    "authorName": "PESANDO MARGHERITA",
    "subjectDesc": ""
  },
  {
    "evtId": 578930,
    "evtCode": "AGNT",
    "evtDatetimeBegin": "2022-05-03T00:00:00+02:00",
    "evtDatetimeEnd": "2022-05-03T23:59:59+02:00",
    "notes": "Link per colloqui prof. AVANZATO",
    "authorName": "AVANZATO PAOLA CARLA",
    "subjectDesc": ""
  }
]
```

### Noticeboards

List school announcements and circulars:

```text
$ classeviva noticeboards list
+---------------------+------+---------------------------------------+
| PUBLICATIONDATE     | READ | TITLE                                 |
+---------------------+------+---------------------------------------+
| 2022-04-28T10:30:00 | true | Comunicazione assemblea di istituto   |
| 2022-04-25T15:45:00 | false| Circolare n. 123 - Uscita anticipata  |
+---------------------+------+---------------------------------------+
```

JSON output:

```shell
$ classeviva noticeboards list --format json
[
  {
    "pubId": 12345,
    "cntTitle": "Comunicazione assemblea di istituto",
    "readStatus": true,
    "pubDT": "2022-04-28T10:30:00",
    "evtCode": "CF",
    "cntValidInRange": true,
    "cntStatus": "active",
    "cntCategory": "General",
    "cntHasAttach": true,
    "attachments": [
      {
        "fileName": "comunicazione.pdf",
        "attachNum": 1
      }
    ]
  }
]
```

Download noticeboard attachments:

```shell
$ classeviva noticeboards download --publication_id 12345 --output-filename ./downloads
+----------------------------------+
| FILE                             |
+----------------------------------+
| ./downloads/12345-documento.pdf  |
+----------------------------------+
```
