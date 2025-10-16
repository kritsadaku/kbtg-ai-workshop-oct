# **PRD: FinTrack Personal Finance Manager**

**Version:** 1.0 (UI Scope)
**Date:** \[Current Date]
**Author:** Product Team
**Status:** Ready for UI Generation (3–4 screens)

---

## **1. Objective (UI Scope)**

Design **3–4 core screens** to let users (bank customers) quickly:

1. set basic goals & permissions, 2) view/categorize transactions, 3) set and track budgets/goals, and 4) see a simple monthly report.
   Non-essential features are deferred.

---

## **2. Screen Map & Navigation**

Bottom tab (or top nav) with 4 entries:

1. **Onboarding** → one-time flow (then hidden)
2. **Transactions**
3. **Budgets & Goals**
4. **Insights**

Primary navigation: Tabs; secondary: in-page cards + drilldowns.

---

## **3. Screens & UI Requirements**

### **Screen A — Onboarding (One-time)**

**Goal:** Collect minimal inputs for personalization and permissions.

**Key UI Sections**

* **Welcome Card**: Logo, brief value prop (“Track spending, set budgets, reach goals.”), CTA **“Get Started”**.
* **Financial Goal Quick Pick**: Chips (e.g., *Monthly budget*, *Debt reduction*, *Save for trip*). Optional custom goal (text).
* **Permissions**: Toggle/check for *“Auto-categorize transactions”* with brief explanatory text and link to privacy.
* **Confirmation**: Summary of selections; CTA **“Finish Setup”**.

**States**

* Empty: default selections unchosen.
* Error: permission denied → show tip to enable later in Settings.

**Microcopy**

* Helper text: “You can edit goals anytime.”

---

### **Screen B — Transactions (List & Categorization)**

**Goal:** Let users see and correct categories fast.

**Key UI Sections**

* **Header**: Month selector (left/right chevrons), account filter (All/Account X).
* **Spend Summary Strip**: *Total outflow*, *Top category*, mini donut chart.
* **Transaction List**: each row shows Merchant, Amount, Date, **Category pill** (tap to change), optional notes icon.
* **Category Picker (Bottom Sheet)**: Search, top categories (Food, Transport, Shopping…), **“Create/Edit Rule”** toggle (apply to future).
* **Bulk Actions** (optional, if time): Multi-select → Re-categorize.

**States**

* Empty: “No transactions yet. Come back after your next purchase.”
* Loading shimmer.
* Error toast on save failure.

**Quick Interactions**

* Swipe left/right to recategorize (optional).
* Long-press → “Always categorize as …” (rule).

---

### **Screen C — Budgets & Goals**

**Goal:** Set monthly category budgets and track simple savings goals.

**Key UI Sections**

* **Budget Overview**: Bar of *Spent / Budget* for this month; **“Add Budget”** button.
* **Budget Cards (per category)**: Category name, progress bar, amount spent and remaining; badges at **80%** and **100%** thresholds.
* **Create/Edit Budget (Modal/Sheet)**: Category selector, amount input, monthly reset note.
* **Goals Section**: Cards with goal name, target amount, current saved, progress ring; CTA **“New Goal”**.
* **Create/Edit Goal (Modal/Sheet)**: Name, target amount, due date (optional).

**States**

* Empty: “No budgets yet. Set one to stay on track.” / “Add your first savings goal.”
* Alert: subtle inline warning when ≥80% used.

---

### **Screen D — Insights (Monthly Report)**

**Goal:** Show a simple monthly summary and top insights.

**Key UI Sections**

* **Month Selector** (same pattern as B).
* **Income vs Expense Card**: two numbers + simple bar comparison.
* **Category Breakdown**: Donut chart + top 5 categories list with amounts & %.
* **Trends Mini-Chart**: Simple line showing last 3 months spend.
* **Highlights/Insights**: Short chips like “Dining likely to overspend” (explainable tooltip: “Based on last 3 months”).

**Actions**

* **Export** button: PDF/CSV (stubbed in UI).

**States**

* Empty: “We’ll build your report after we see some activity.”

---

## **4. Minimal Data Model (UI Binding)**

* **Transaction**: `{ id, date, merchant, amount, categoryId, notes?, isRecurring?, ruleApplied? }`
* **Category**: `{ id, name, icon? }`
* **Budget**: `{ id, categoryId, month, amount, spent }`
* **Goal**: `{ id, name, targetAmount, savedAmount, targetDate? }`
* **Insight**: `{ id, type, message, basis }` *(e.g., type=“forecast”, basis=“last\_3\_months”)*

---

## **5. Interaction Rules (UI)**

* **Auto-Categorization**: Show badge “Auto” on newly categorized items; editing a category prompts **“Apply to future?”**.
* **Budget Alerts**: Inline badge at **80%** and **100%** usage on Budget Cards; optional push (out of scope for UI MVP).
* **Drilldowns**: Tapping a category in Insights opens a filtered view of Transactions for that category.

---

## **6. Visual/Components (Design System Hints)**

* Charts: donut (category), line (trend), bar (income vs expense).
* Cards: elevated, rounded corners (md), soft shadows.
* Tokens: use brand primary for highlights, neutrals for lists, semantic warning at 80%/danger at 100%.
* Bottom sheets/modals for quick edits (category picker, add budget/goal).

---

## **7. Non-Functional (UI-Level)**

* **Accessibility**: AA contrast, chart labels with text equivalents.
* **Performance**: Lists virtualized; actions show optimistic toasts.
* **Localization**: Numbers/currency formatting via system locale.

---

## **8. Acceptance Criteria (UI Cut)**

**Onboarding**

* [ ] User can pick at least one goal OR enable auto-categorization and finish setup.
* [ ] Permission state is clearly shown and recoverable later.

**Transactions**

* [ ] List groups by date; category pill is editable in-place.
* [ ] Category picker supports search and “apply to future” rule.
* [ ] Spend summary and mini donut reflect current month & filters.

**Budgets & Goals**

* [ ] User can create/edit a category budget; progress shows live *spent/amount*.
* [ ] Badges appear at **≥80%** and **≥100%** budget usage.
* [ ] User can create a savings goal and see progress.

**Insights**

* [ ] Monthly Income vs Expense displayed; donut shows top categories.
* [ ] Trend line shows at least last 3 months when data exists.
* [ ] Tapping a category opens a filtered Transactions view.

**Export (Stub)**

* [ ] Tapping **Export** reveals choices (PDF/CSV) and shows a success toast (mocked).

---

## **9. Out of Scope (This Cut)**

* Full notification center / scheduling preferences
* Advanced merchant analytics & recurring detection management
* Multi-currency & shared budgets
* Detailed AI explanations beyond short tooltips

---

## **10. Sample Copy**

* Empty Transactions: “No activity yet. Your next purchase will appear here.”
* 80% Budget: “You’ve used 80% of your Dining budget.”
* Goal Created: “Goal saved. Let’s make steady progress!”

---
