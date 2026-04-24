# Bug Report Template

## Bug ID: `<PROJECT-XXX>`

---

### 1. Basic Information

| Field | Value |
|-------|-------|
| **Project** | Skill Hub Pro |
| **Module** | `[auth/skills/router/admin/ui]` |
| **Severity** | `[Critical/High/Medium/Low]` |
| **Priority** | `[P0/P1/P2/P3]` |
| **Status** | `[New/Assigned/In Progress/Resolved/Verified/Closed]` |
| **Reported By** | |
| **Reported Date** | `YYYY-MM-DD` |
| **Environment** | `[dev/staging/production]` |

### 2. Description

**Summary**: (One-line description of the issue)

**Detailed Description**:
```
(Detailed explanation of what the bug is)
```

### 3. Steps to Reproduce

1. 
2. 
3. 

### 4. Expected Behavior

```
(What should happen)
```

### 5. Actual Behavior

```
(What actually happens)
```

### 6. Screenshots / Logs

```
(Paste relevant logs, screenshots, or API responses)
```

### 7. Test Data

```
(Query, payload, or conditions used to reproduce)
```

### 8. Root Cause Analysis

*(To be filled by developer)*

```
Root cause:
Fix approach:
Affected files:
```

### 9. Fix Verification

| Check | Result |
|-------|--------|
| Unit test added | [Yes/No] |
| Integration test updated | [Yes/No] |
| Regression test passed | [Yes/No] |
| Code review completed | [Yes/No] |

### 10. Related PR / Commit

- PR: `#XXX`
- Commit: `<commit-hash>`

---

## Severity & Priority Definitions

| Severity | Definition |
|----------|------------|
| **Critical** | System crash, data loss, security breach, core feature completely broken |
| **High** | Major feature broken, no workaround available |
| **Medium** | Feature partially broken, workaround exists |
| **Low** | Cosmetic issue, minor UI glitch, documentation error |

| Priority | Definition |
|----------|------------|
| **P0** | Fix immediately (blocking release) |
| **P1** | Fix before next release |
| **P2** | Fix in current sprint if possible |
| **P3** | Backlog / future sprint |

---

## Regression Test Checklist

After bug fix, verify these related areas are not broken:

- [ ] Skill listing (with/without filters)
- [ ] Skill search (keyword, category)
- [ ] Router match endpoint
- [ ] Router execute endpoint
- [ ] User registration / login
- [ ] Admin CRUD operations
- [ ] API key management
- [ ] Authentication (JWT, GitHub OAuth)
- [ ] Rate limiting
