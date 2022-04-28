# classeviva

Classeviva is a Go library and CLI tool to access the popular school portal https://web.spaggiari.eu.

## Grades

Text output:

```shell
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

```shell
$ classeviva agenda list --until 2022-04-27 --limit 2
[
  {
    "evtId": 550756,
    "evtCode": "AGNT",
    "evtDatetimeBegin": "2022-04-26T08:00:00+02:00",
    "evtDatetimeEnd": "2022-04-26T09:00:00+02:00",
    "notes": "PORTARE LETTERATURA",
    "authorName": "DICEMBRE ELISA",
    "subjectDesc": ""
  },
  {
    "evtId": 537508,
    "evtCode": "AGNT",
    "evtDatetimeBegin": "2022-04-26T11:00:00+02:00",
    "evtDatetimeEnd": "2022-04-26T12:00:00+02:00",
    "notes": "Studiare Intermezzo da Cavalleria rusticana. Interrogazione.",
    "authorName": "GRIMALDI ALESSANDRO",
    "subjectDesc": ""
  }
]
```
