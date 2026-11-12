# Noticeboard fields and semantic investigation

Field-preservation implementation, review command, and follow-up research for [issue #100](https://github.com/zmoog/classeviva/issues/100), including the [timezone and workflow investigation request](https://github.com/zmoog/classeviva/issues/100#issuecomment-6821039056).

## Sources and method

Primary sources inspected on 2026-11-12:

1. Authenticated `GET /rest/v1/students/{studentId}/noticeboard`. The initial audit contained 63 notices and 88 attachments; the implementation-time refresh contained 64 notices and 89 attachments. These are snapshots of one authorized account, not an exhaustive API contract.
2. `GET /rest/v1/students/{studentId}/noticeboard/attach/{evtCode}/{pubId}/{attachNum}` for the 82 attachments belonging to the 56 **already-read** notices in the refreshed response. The seven attachments on unread notices were deliberately skipped. No mark-as-read or other action endpoint was called.
3. Local extraction of all 82 downloaded PDF/DOCX files. Extraction succeeded for all files; one produced no text and was not OCRed. Keyword searches located signature, participation, reply, and submission instructions; selected matching passages and their surrounding instructions were manually reviewed. Keyword matches alone were not treated as proof of an action requirement.
4. The user-provided website screenshot of a notice with a displayed Deadline matching its API `cntValidTo` date. The website's own network traffic and frontend source were not inspected.

Raw responses, documents, and extracted text were kept outside the repository and deleted after investigation. This note contains only aggregate observations and paraphrased examples: no credentials, student IDs, school names, document titles, email addresses, or external form links. The test fixture is entirely synthetic; its true workflow/change flags are hypothetical, not positive production examples.

## Fields preserved by the client

`spaggiari.Noticeboard` now preserves all 21 top-level fields observed in both snapshots, plus both nested attachment fields. Existing JSON keys and types remain unchanged; the new properties are additive. The existing `list` text table is unchanged; `review` provides a separate attention-oriented view.

| Newly captured Go field | JSON key | Wire type / interpretation |
| --- | --- | --- |
| `ContentID` | `cntId` | Integer content ID, distinct from publication ID `pubId`. |
| `Number` | `cntNum` | String website notice number, distinct from a circular number embedded in the title. |
| `ValidFrom` | `cntValidFrom` | Date-only string, `YYYY-MM-DD`. |
| `ValidTo` | `cntValidTo` | Date-only publication expiry, `YYYY-MM-DD`. |
| `HasChanged` | `cntHasChanged` | Boolean; exact change-tracking semantics unverified. |
| `NeedJoin` | `needJoin` | Boolean; likely participation/adhesion workflow flag. |
| `NeedReply` | `needReply` | Boolean; likely reply workflow flag. |
| `NeedFile` | `needFile` | Boolean; likely file-submission workflow flag. |
| `NeedSign` | `needSign` | Boolean; likely signature/acknowledgment workflow flag. |
| `EventID` | `evento_id` | String; equals the string representation of `cntId` in all sampled notices. |
| `InsertedAt` | `dinsert_allegato` | Raw `YYYY-MM-DD HH:MM:SS` string without an offset; exact provenance unknown. |

The original 10 fields are `pubId`, `pubDT`, `readStatus`, `evtCode`, `cntValidInRange`, `cntStatus`, `cntTitle`, `cntCategory`, `cntHasAttach`, and `attachments`. Every sampled attachment has exactly `fileName` and `attachNum`.

## Timezone: Europe/Rome is plausible, not established for every field

The refreshed response contains 49 `pubDT` timestamps with `+02:00` and 15 with `+01:00`. All 64 explicit offsets match `Europe/Rome` at their respective instants. This supports the comment's CET/CEST hypothesis, but other regions share those offsets, and the offset-bearing publication field does not establish the timezone of another field.

`dinsert_allegato` has **no timezone information**. In the initial 63-record audit, its literal wall-clock value preceded `pubDT` by 0–10 seconds in 61 records, by 124 seconds in one, and by 6h17m57s in one. This is consistent with local insertion followed by publication, but does not prove its exact event or timezone. It is populated even for a notice with no attachments, so “attachment upload time” would be an unjustified label.

`cntValidFrom` and `cntValidTo` contain **only dates**. The screenshot's winter Deadline displays 01:00, while several summer dates display 02:00. UTC-midnight conversion to Europe/Rome would explain that pattern, but frontend conversion was not verified. Neither a timezone nor a deadline hour is added by the client.

If future functionality needs timezone-aware school-local calculations, investigate using the IANA zone `Europe/Rome`, not fixed CET (`+01:00`), so daylight saving time is handled. This change deliberately performs no new date parsing or conversion.

## Workflow flags: content actions are not reliably represented by them

All four `need*` flags and `cntHasChanged` were false on **all 64** refreshed notices. Nevertheless, attachment review found these concrete counterexamples to interpreting false as “no action requested anywhere in the notice”:

- A course-registration notice instructs applicants to complete an external web form and include a completed, signed authorization. The corresponding notice has `needJoin`, `needReply`, `needFile`, and `needSign` all false.
- Another activity-registration notice instructs applicants to email an attached form to a teacher. The attached form includes student/guardian signature fields. All four corresponding flags are false.

These examples support a **portal-workflow metadata** interpretation rather than a general summary of the prose or attached forms. They do not prove whether a flag means a configured capability, an obligation, or an outstanding recipient-specific action. Eligibility, previous completion, or school configuration could also matter. A keyword such as “firma” can also refer to a school official's digital signature rather than a parent's task, so keyword counting is insufficient.

Safe behavior today:

- Preserve each boolean exactly as returned.
- Treat all four names as tentative workflow descriptions, not independently verified requirements.
- Do not use false to hide a notice from an “actions requested in content” view.
- Do not derive or overwrite API flags from document keywords.

To resolve the remaining ambiguity, obtain an authorized positive (`true`) example, inspect the website action associated with it, and compare read-only observations before/after a human independently completes that action. Do not automatically sign, join, reply, or upload to investigate.

## Publication expiry is not the application deadline

In the activity-registration example above, the explicit application cutoff in the attachment precedes `cntValidTo` by more than a month. In the external-form example, the action cutoff and validity end date coincide. Therefore publication validity can coincide with an action deadline, but must not be treated as a universal action deadline.

All 28 expired notices in the refreshed response still have `cntStatus: "active"`; 36 are in range. Use validity information rather than `Status == "active"` to reason about publication expiry. The current list command does not change its filtering behavior.

## Schema and regression coverage

- `openapi.yaml` now documents `cntNum` as a string.
- `dinsert_allegato` is documented with its observed space-separated representation instead of incorrectly declaring RFC 3339 `format: date-time`.
- Workflow, expiry, identifier, and timezone descriptions distinguish observations from hypotheses.
- `adapters/spaggiari/testdata/noticeboards.json` covers both explicit publication offsets, date-only validity values, a leading-zero string notice number, distinct publication/content IDs, hypothetical true flags, false flags, an expired-but-active notice, an empty category, and empty/nonempty attachment arrays.
- Regression tests check typed field access, complete JSON decode/re-encode preservation, and JSON result output through the existing formatter. No production API access is needed to run them.

## Read-only review command

`classeviva noticeboards review --days 14` combines the two independent signals: any workflow flag, or an unread notice with a validity end date from today through 14 calendar days ahead (inclusive). The default window is 7 days, and zero means today only. It does not change `noticeboards list`.

The `commands.ReviewNoticeboards` function accepts the reference time explicitly, calculates Europe/Rome calendar dates, and returns selected notices with their reasons and nullable days until expiry. It does not alter input data. Missing or malformed expiry values do not suppress a workflow flag. Combined signals sort before flags alone, then unread expiry alone; expiry sorts ascending within each group, with unknown dates last and stable ties.

The result formatter keeps the notice and derived attention data separate in JSON. Text labels are deliberately conservative: workflow flags are “set,” not “pending,” and expired publication validity is not called an overdue task. No read-marking, download, or action endpoint is invoked. In line with the caveats above, this is a prioritization heuristic rather than an exhaustive task detector.

Tests cover independent signals, all workflow flags, window boundaries, zero/negative windows, invalid/missing expiry, read and expired flagged notices, stable sorting, winter/summer local midnight, both DST transitions, empty output, JSON field preservation, and argument validation. A synthetic HTTP-transport test exercises profile selection through the Cobra command and allows only a single notice-list GET, rejecting any unexpected operation. See the README for usage and JSON reason codes.
