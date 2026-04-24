# Skill Hub Pro - User Acceptance Test (UAT) Plan

## 1. Overview

| Item | Detail |
|------|--------|
| **Project** | Skill Hub Pro - 技能聚合网站 |
| **Phase** | UAT - 用户验收测试 |
| **PRD Version** | V3.0 |
| **Date** | 2026-04-24 |

## 2. Test Objectives

- Validate that all PRD V3.0 acceptance criteria are met
- Verify the user-facing functionality works as expected in a production-like environment
- Confirm the Smart Router achieves >= 85% Top-1 accuracy
- Ensure the system meets performance and security requirements

## 3. User Personas

| Persona | Role | Goals |
|---------|------|-------|
| **Alice** | Regular Developer | Search skills, execute routing, provide feedback |
| **Bob** | Team Lead | Browse categories, manage API keys, view analytics |
| **Carol** | Admin | Manage skills/categories, review logs, configure system |
| **Dave** | First-time User | Register, explore skills, understand platform |

## 4. Test Scenarios

### 4.1 User Registration & Authentication

| ID | Scenario | Persona | Steps | Expected Result | PRD Ref |
|----|----------|---------|-------|-----------------|---------|
| UAT-01 | New user registers with valid info | Dave | 1. Visit landing page<br>2. Click "Register"<br>3. Enter username/email/password<br>4. Submit | Account created, redirected to dashboard with JWT token | FR-01 |
| UAT-02 | User logs in with credentials | Alice | 1. Click "Login"<br>2. Enter username and password<br>3. Submit | Logged in, sees skill browsing page | FR-01 |
| UAT-03 | User logs in with email | Alice | 1. Enter email instead of username<br>2. Enter password<br>3. Submit | Logged in successfully | FR-01 |
| UAT-04 | User logs in with wrong password | Alice | 1. Enter username<br>2. Enter incorrect password | Error message: invalid credentials | FR-01 |
| UAT-05 | Duplicate username registration | Dave | 1. Enter existing username<br>2. Fill other fields | Error message: username already exists | FR-01 |
| UAT-06 | GitHub OAuth login | Alice | 1. Click "Login with GitHub"<br>2. Authorize on GitHub<br>3. Redirect back | Logged in, new account auto-created if first time | FR-01 |

### 4.2 Skill Browsing & Search

| ID | Scenario | Persona | Steps | Expected Result | PRD Ref |
|----|----------|---------|-------|-----------------|---------|
| UAT-07 | Browse skill list | Bob | 1. Navigate to /skills<br>2. View all skills | Skills displayed with name, stars, description | FR-02 |
| UAT-08 | Filter by category | Bob | 1. Select category "Data Analysis"<br>2. Apply filter | Only data analysis skills shown | FR-02 |
| UAT-09 | Sort by stars | Bob | 1. Click "Sort by Stars" | Skills sorted descending by star count | FR-02 |
| UAT-10 | Search by keyword | Alice | 1. Enter "excel" in search bar<br>2. Press Enter | Skills matching "excel" shown | FR-02 |
| UAT-11 | Paginate results | Bob | 1. Browse skills with many results<br>2. Click page 2 | Page 2 results loaded | FR-02 |
| UAT-12 | View skill detail | Alice | 1. Click on a skill card<br>2. View full information | Skill detail page with full description, repo link, stars | FR-02 |
| UAT-13 | Empty search results | Alice | 1. Search for "zzzzznotexist" | Empty state message displayed | FR-02 |

### 4.3 Smart Router

| ID | Scenario | Persona | Steps | Expected Result | PRD Ref |
|----|----------|---------|-------|-----------------|---------|
| UAT-14 | Route natural language query | Alice | 1. Enter "analyze excel data trends"<br>2. Submit to router | Top matched skill returned with confidence score | FR-03 |
| UAT-15 | Execute matched skill | Alice | 1. After match, click "Execute"<br>2. Wait for result | Skill executed, result displayed | FR-03 |
| UAT-1 | Route accuracy check | Carol | 1. Run 100 test queries<br>2. Compare Top-1 with expected | >= 85% Top-1 accuracy | NFR-01 |
| UAT-17 | Provide feedback on match | Alice | 1. After execution<br>2. Rate match quality (1-5)<br>3. Submit feedback | Feedback recorded, score updated | FR-03 |

### 4.4 Admin Functions

| ID | Scenario | Persona | Steps | Expected Result | PRD Ref |
|----|----------|---------|-------|-----------------|---------|
| UAT-18 | Add new skill | Carol | 1. Go to Admin > Skills<br>2. Click "Add Skill"<br>3. Fill form and submit | New skill appears in listing | FR-04 |
| UAT-19 | Edit existing skill | Carol | 1. Find a skill<br>2. Click "Edit"<br>3. Update fields and save | Changes reflected immediately | FR-04 |
| UAT-20 | Delete a skill | Carol | 1. Select skill<br>2. Click "Delete"<br>3. Confirm deletion | Skill removed from listing | FR-04 |
| UAT-21 | Manage categories | Carol | 1. Go to Admin > Categories<br>2. Add/Edit/Delete category | Categories updated | FR-04 |
| UAT-22 | View router logs | Carol | 1. Go to Admin > Logs<br>2. View recent routing activity | Logs displayed with queries, matched skills, duration | FR-04 |

### 4.5 Performance & Reliability

| ID | Scenario | Persona | Steps | Expected Result | PRD Ref |
|----|----------|---------|-------|-----------------|---------|
| UAT-23 | Search response time | Alice | 1. Execute 10 searches<br>2. Measure response time | P95 < 3s | NFR-02 |
| UAT-24 | Router match response time | Alice | 1. Execute 10 match requests<br>2. Measure response time | P95 < 5s | NFR-02 |
| UAT-25 | System uptime | All | 1. Monitor over 24 hours | 99.9% uptime | NFR-03 |

### 4.6 Security

| ID | Scenario | Persona | Steps | Expected Result | PRD Ref |
|----|----------|---------|-------|-----------------|---------|
| UAT-26 | Unauthenticated admin access | Dave | 1. Try to access /admin without token | 401 Unauthorized | NFR-04 |
| UAT-27 | Non-admin access to admin API | Alice | 1. Login as regular user<br>2. Call admin API | 403 Forbidden | NFR-04 |
| UAT-28 | SQL injection attempt | Hacker | 1. Search with SQL injection payload | Query safely handled, no SQL error exposed | NFR-04 |
| UAT-29 | XSS attempt | Hacker | 1. Submit XSS payload in search | Payload sanitized, no script execution | NFR-04 |

## 5. UAT Execution

### 5.1 Environment
- **URL**: [Production Staging URL]
- **Test Data**: Pre-seeded dataset of 50+ skills across 8+ categories
- **User Accounts**: Pre-created accounts for each persona

### 5.2 Entry Criteria
- [ ] All unit tests pass (5.2)
- [ ] Integration tests pass (5.3)
- [ ] Router accuracy >= 85% (5.4)
- [ ] No critical security issues (5.6)

### 5.3 Exit Criteria
- [ ] All UAT-01 to UAT-29 executed and pass
- [ ] All critical bugs fixed and verified
- [ ] Router accuracy meets PRD requirement
- [ ] Sign-off obtained from stakeholders

### 5.4 Test Results Summary

| Category | Total | Pass | Fail | Blocked | Pass Rate |
|----------|-------|------|------|---------|-----------|
| Registration & Auth | | | | | |
| Skill Browsing | | | | | |
| Smart Router | | | | | |
| Admin Functions | | | | | |
| Performance | | | | | |
| Security | | | | | |
| **Total** | | | | | |

## 6. Bug Reporting

Found bugs should be filed using the BUG template (see `tests/bug_report_template.md`).
Critical and High severity bugs must be resolved before UAT sign-off.

## 7. Sign-off

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Owner | | | |
| QA Lead | | | |
| Dev Lead | | | |

---
*Document Version: 1.0*
