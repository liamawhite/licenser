---
name: Maintainability Workstream
about: Create an issue to add to the maintainability workstream project.
title: ''
labels: bug
assignees: ''

---

|Resource|Link or Approver|
|---|---|
|Product Design|A link to the PRD/User Story Doc or the approving member of the product team|
|Engineering Design|A link to the design doc or the approving engineer|
|Github Project|Each workstream has its own Github project to track implementation progress|

## Description
Align user experience of installing the data plane with that of the control and management plane. Remove our reliance on Istio Operator upstream by consuming it as a binary and not as a library.

## Customer Impact
Current experience is disjointed. We use our own API for the control and management plane but an obscured Istio one for data. This UX results in subtle differences that can confuse users.
