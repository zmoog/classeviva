# classeviva

Classeviva is a Go library and CLI tool to access the popular school portal https://web.spaggiari.eu.

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
