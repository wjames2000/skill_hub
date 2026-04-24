# Frontend UI Walkthrough Checklist

## Setup
- [ ] Build frontend without errors: `npm run build`
- [ ] TypeScript compilation passes: `npm run typecheck`
- [ ] All existing tests pass: `npm run test`

## 1. Landing Page

### 1.1 Layout & Responsiveness
- [ ] Page loads within 3s on 4G connection
- [ ] Header with logo, navigation, and auth buttons visible
- [ ] Hero section with tagline and CTA button
- [ ] Footer with links and copyright
- [ ] Responsive on mobile (375px), tablet (768px), desktop (1440px)
- [ ] No horizontal scroll on any breakpoint

### 1.2 Navigation
- [ ] "Skills" link navigates to /skills
- [ ] "Login" button navigates to /login
- [ ] "Register" button navigates to /register
- [ ] Logo links to homepage
- [ ] Active nav item is highlighted

## 2. Authentication Pages

### 2.1 Login Page (/login)
- [ ] Form with username/email and password fields
- [ ] "Login with GitHub" button visible
- [ ] Validation: empty fields show error messages
- [ ] Validation: invalid email format caught
- [ ] Successful login redirects to /skills
- [ ] Failed login shows error toast/message
- [ ] "Register here" link navigates to /register
- [ ] Loading state during API call (spinner/disabled button)

### 2.2 Register Page (/register)
- [ ] Form with username, email, password, confirm password
- [ ] Password strength indicator (if implemented)
- [ ] Validation: username min length
- [ ] Validation: email format
- [ ] Validation: password min length
- [ ] Validation: passwords match
- [ ] Successful registration redirects to /skills
- [ ] Duplicate username shows error
- [ ] "Login here" link navigates to /login

### 2.3 GitHub OAuth
- [ ] "Login with GitHub" redirects to GitHub authorization
- [ ] After GitHub auth, user is redirected back and logged in
- [ ] First-time GitHub user gets auto-registered

## 3. Skill Browsing

### 3.1 Skill List Page (/skills)
- [ ] Skills displayed as cards with name, stars, description
- [ ] Category filter dropdown with all categories
- [ ] Sort options: by stars, by name, by date
- [ ] Pagination controls (prev/next, page numbers)
- [ ] Search bar with placeholder text
- [ ] Search triggers on Enter or search icon click
- [ ] Loading skeleton/spinner shown during fetch
- [ ] Empty state when no results ("No skills found")
- [ ] Error state with retry button
- [ ] Search results update within 1s of typing (debounced)

### 3.2 Skill Card
- [ ] Card shows name, display name, stars, category badge
- [ ] Stars display correctly (e.g., "⭐ 100")
- [ ] Click on card navigates to skill detail page
- [ ] Card has hover effect (elevation/scale)
- [ ] Category badge color matches category

### 3.3 Skill Detail Page (/skills/:id)
- [ ] Full skill information displayed
- [ ] Repository link opens in new tab
- [ ] Execute button visible (if authenticated)
- [ ] Back to list link
- [ ] 404 page for non-existent skill ID
- [ ] Loading state during fetch
- [ ] Error state for network failure

### 3.4 Category Sidebar/Filters
- [ ] Clicking category filters skills correctly
- [ ] Active category is visually distinct
- [ ] "All" option resets filter
- [ ] Category count shown (if available)

## 4. Smart Router UI

### 4.1 Query Input
- [ ] Text input with placeholder "Describe what you need..."
- [ ] Submit button (search icon or "Find Skill")
- [ ] Input validation: empty query disabled
- [ ] Loading state while matching
- [ ] Character limit indication (if applicable)

### 4.2 Match Results
- [ ] Matched skill(s) displayed with name and confidence score
- [ ] Score displayed as percentage or badge (e.g., "95% match")
- [ ] Multiple matches shown in priority order
- [ ] Confidence bar/indicator
- [ ] Execute button on matched skill
- [ ] No match state: "No matching skill found"
- [ ] Error state with retry option

### 4.3 Execution Results
- [ ] Result displayed in a readable format
- [ ] Loading spinner during execution
- [ ] Session ID shown
- [ ] Duration displayed

### 4.4 Feedback Component
- [ ] Star rating (1-5) or thumbs up/down
- [ ] Feedback submitted confirmation
- [ ] Feedback recorded per match/execution

## 5. Admin Pages

### 5.1 Admin Layout
- [ ] Sidebar with navigation: Skills, Categories, Users, Logs, Settings
- [ ] Only accessible to admin users (redirect if not admin)
- [ ] Responsive sidebar (collapsible on mobile)

### 5.2 Skill Management
- [ ] Table/list of all skills with CRUD buttons
- [ ] "Add Skill" button opens form (modal or page)
- [ ] Form validation for required fields
- [ ] Edit pre-fills form with existing data
- [ ] Delete shows confirmation dialog
- [ ] Success/error toast after operations
- [ ] Search/filter within admin list

### 5.3 Category Management
- [ ] List of categories with sort order
- [ ] Add/Edit/Delete categories
- [ ] Drag-and-drop reorder (if implemented)

### 5.4 Router Logs
- [ ] Log table with columns: query, matched skill, strategy, duration, timestamp
- [ ] Sortable columns
- [ ] Search by query text
- [ ] Pagination

### 5.5 API Keys
- [ ] List of user API keys
- [ ] Generate new key button
- [ ] Copy key to clipboard
- [ ] Revoke/delete key with confirmation

## 6. Common Elements

### 6.1 Navigation Bar
- [ ] Shows user avatar/name when logged in
- [ ] Dropdown menu: Profile, API Keys, Logout
- [ ] Admin link visible only to admin users
- [ ] Mobile hamburger menu
- [ ] Active page highlighted
- [ ] Notification badge (if implemented)

### 6.2 Toast/Notification System
- [ ] Success toasts (green)
- [ ] Error toasts (red)
- [ ] Warning toasts (yellow)
- [ ] Auto-dismiss after 5s
- [ ] Manual dismiss button
- [ ] Stacked notifications

### 6.3 Loading States
- [ ] Skeleton loaders for cards/lists
- [ ] Spinner for inline actions
- [ ] Disabled buttons during API calls
- [ ] Progress bar for long operations

### 6.4 Empty States
- [ ] "No skills found" for empty search
- [ ] "No results" for empty filters
- [ ] "No logs" for empty log view
- [ ] Illustration or helpful message

### 6.5 Error States
- [ ] Network error with retry button
- [ ] Server error with user-friendly message
- [ ] 404 page with navigation back
- [ ] 403 page for unauthorized access

## 7. Accessibility (a11y)

- [ ] All images have alt text
- [ ] Form inputs have associated labels
- [ ] Keyboard navigation works (Tab, Enter, Escape)
- [ ] Focus indicators visible
- [ ] Color contrast meets WCAG AA standards
- [ ] Screen reader friendly (aria attributes)
- [ ] Skip to content link

## 8. Performance

### 8.1 Page Load (Lighthouse)
- [ ] First Contentful Paint (FCP) < 2s
- [ ] Largest Contentful Paint (LCP) < 2.5s
- [ ] Time to Interactive (TTI) < 3s
- [ ] Cumulative Layout Shift (CLS) < 0.1
- [ ] Lighthouse score > 80 for all categories

### 8.2 Runtime
- [ ] Smooth scrolling (60fps)
- [ ] No UI jank during data loading
- [ ] Infinite scroll/pagination works smoothly
- [ ] Images lazy-loaded

## 9. Cross-browser Testing

- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Edge (latest)
- [ ] Mobile Chrome (Android)
- [ ] Mobile Safari (iOS)

## 10. Bug Check Summary

| Check | Result | Notes |
|-------|--------|-------|
| Build succeeds | [Pass/Fail] | |
| All tests pass | [Pass/Fail] | |
| All pages load correctly | [Pass/Fail] | |
| All forms validate correctly | [Pass/Fail] | |
| All modals/dialogs work | [Pass/Fail] | |
| Responsive on 3 breakpoints | [Pass/Fail] | |
| Accessible (keyboard+screen reader) | [Pass/Fail] | |
| Cross-browser compatible | [Pass/Fail] | |

## 11. Issues Found

| # | Page | Issue | Severity | Status |
|---|------|-------|----------|--------|
| | | | | |

---
*Template Version: 1.0*
