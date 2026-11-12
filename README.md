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

Configure credentials in a private local terminal, not an agent transcript. **The current `profile add` prompt echoes input**; asterisks below are documentation redaction, not password masking. Use the [safe profile-setup procedure](.agents/skills/classeviva-cli/references/secrets.md#human-only-profile-setup) to enter credentials without echoing them or putting them in command arguments.

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

Classeviva supports multiple authentication methods with a priority chain, making it easy to manage credentials for multiple students. Read the [credential and private-data safety guide](.agents/skills/classeviva-cli/references/secrets.md) before setup or agent use. Profiles store plaintext passwords, identity caches contain authentication tokens, and `--debug` can expose authentication responses.

### Authentication Priority

Credentials are resolved in the following order (highest priority first):

1. **CLI flags**: `--username` and `--password`
2. **Profile-based**: `--profile` flag or `default_profile` from config file
3. **Environment variables**: `CLASSEVIVA_USERNAME` and `CLASSEVIVA_PASSWORD`

### Profile Management (Recommended)

For families with multiple students, profile-based authentication is the recommended approach.

#### Setup Profiles

For private credential entry, follow the [human-only profile setup](.agents/skills/classeviva-cli/references/secrets.md#human-only-profile-setup). The examples below illustrate the flow with redacted passwords:

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

Display profile information (the password is masked; the username is displayed):

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
```

Credential flags are supported, but avoid passing real passwords in arguments: they can appear in process listings, shell history, and agent logs. Prefer an existing profile or trusted secret-manager injection as described in the [safety guide](.agents/skills/classeviva-cli/references/secrets.md).

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

**Note**: New config files are created with restrictive permissions (0600), but passwords are not encrypted. Existing file permissions may need tightening; see the [safety guide](.agents/skills/classeviva-cli/references/secrets.md#storage-and-current-limitations).

### Environment Variables (Legacy)

For backward compatibility, a trusted launcher or secret manager can inject `CLASSEVIVA_USERNAME` and `CLASSEVIVA_PASSWORD`. Avoid literal `export` commands containing real secrets, and remember that a selected/default profile takes precedence over these variables. The CLI does not automatically load `.env` files. See [environment-only automation](.agents/skills/classeviva-cli/references/secrets.md#environment-only-automation) for isolated token-cache handling.

```shell
# Both credential variables have already been injected privately
classeviva grades list
```

### Identity Caching

Authentication tokens are cached per-profile in `~/.classeviva/identity-{profile}.json` to minimize API calls. Tokens are automatically refreshed when expired.

## Agent skill

The portable [Classeviva CLI skill](.agents/skills/classeviva-cli/SKILL.md) lives in `.agents/skills/classeviva-cli/`. It covers grades, agenda, notice review, profile selection, and the authorization required before downloading attachments (which marks notices read).

Agents supporting the project `.agents/skills/` convention can discover it after reload/restart. For other agents, load `SKILL.md` explicitly or copy the entire `classeviva-cli` directory, including `references/`, into a supported skills location. In Pi, it can be invoked as `/skill:classeviva-cli`. Loading a skill does not configure credentials or grant permission to change school data.

The bundled [credential and private-data safety guide](.agents/skills/classeviva-cli/references/secrets.md) documents plaintext storage, private human setup, environment injection, token caches, debug-log risks, and incident handling. **The API blocks requests originating from cloud service providers (CSPs)**; hosted agents and CI runners can receive HTTP 403 with valid credentials. Run authenticated operations from an authorized local machine or approved non-CSP network; see [network restrictions](.agents/skills/classeviva-cli/references/secrets.md#network-restrictions-and-http-403).

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

#### Review notices

Surface notices worth reviewing without changing their read status:

```shell
# Default: unread notices expiring within 7 days, plus any workflow-flagged notices
classeviva noticeboards review

classeviva --profile older-kid noticeboards review --days 14
classeviva noticeboards review --days 14 --format json

# Exclude read notices, even those with workflow flags
classeviva noticeboards review --days 14 --unread

# Today only for the unread-expiry signal; workflow flags are still included
classeviva noticeboards review --days 0
```

By default, a notice is included when **any `needJoin`, `needReply`, `needFile`, or `needSign` flag is true**, **or** when it is unread and `cntValidTo` falls between today and today + `--days`, inclusive. The default is 7 days; negative values are rejected. Date comparisons use `Europe/Rome` calendar dates, including across daylight-saving transitions.

`--unread` narrows the review queue to unread notices only, overriding the default inclusion of read notices with workflow flags. It does not include every unread notice: the workflow-or-approaching-expiry criteria still apply. It is disabled by default and works with both text and JSON output.

Notices matching both signals come first, followed by other workflow-flagged notices, then notices matching only unread expiry. Within each group, earliest expiry comes first, unknown expiry last, and ties preserve server order.

The text view shows publication ID, website notice number, read status, validity end date, attention reasons, and title. Reasons say “Signature flag set” rather than claiming a signature is still pending. Expired flagged notices remain visible as **expired**, not “overdue”; flagged notices with missing or malformed expiry dates show **Unknown**. Unread notices with expired, missing, or malformed expiry dates are not included unless a workflow flag is set.

JSON output is an array of objects containing the original notice under `notice` and derived signals under `attention`:

- `attention.reasons`: `participation_flag`, `reply_flag`, `file_submission_flag`, `signature_flag`, `unread`, and/or `expires_soon`.
- `attention.days_until_expiry`: calendar days remaining, negative for expired notices, or `null` for unknown expiry.

This is a **heuristic review queue**, not a definitive task list: all-false workflow flags do not prove no work is required, reading does not mean completion, and publication expiry may differ from an action deadline in the content. The command only lists notices; it does not mark them read, download attachments, or perform actions. `noticeboards list` retains its existing behavior and output shape.

Download noticeboard attachments:

```shell
$ classeviva noticeboards download --publication_id 12345 --output-filename ./downloads
+----------------------------------+
| FILE                             |
+----------------------------------+
| ./downloads/12345-documento.pdf  |
+----------------------------------+
```
