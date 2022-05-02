# classeviva

Classeviva is a Go library and CLI tool to access the popular school portal https://web.spaggiari.eu.

## Grades

Text output:

```text
$ classeviva grades list --limit 3
+------------+-------+-----------------+-------------------------------+
| DATE       | GRADE | SUBJECT         | NOTES                         |
+------------+-------+-----------------+-------------------------------+
| 2022-04-27 | 9     | ARTE E IMMAGINE |                               |
| 2022-04-22 | 7+    | COMPORTAMENTO   | comportamento della settimana |
|            | 7     | SCIENZE         |                               |
+------------+-------+-----------------+-------------------------------+%
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

## Agenda

Text output:

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

```text
$ classeviva agenda list --until 2022-04-27 --limit 2
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
