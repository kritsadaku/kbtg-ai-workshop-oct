# PRD: LBK Points — Wallet, Transfer & Shopping

---

## 1) Overview

LBK Points is a lightweight loyalty wallet for (a) tracking point balance and history, (b) transferring/receiving points between users (incl. QR), and (c) redeeming points for goods with SMS receipts. The product ships as a responsive web app.

### Goals

* Give members a clear balance and recent activity view.
* Make point transfers and QR-based “request points” fast and error-free.
* Enable shopping/redeeming with LBK Points end-to-end (catalog → cart → checkout) and send receipts via SMS.
* Reduce support by making history and receipts self-service.

### Non-Goals (v1)

* Cash/credit card payments.
* Coupons, multi-currency, or returns/refunds.
* Complex customer CRM; only “attach customer to order” metadata.

---

## 2) Users & Primary Use Cases

* **Member / Shopper**: checks balance, scans QR to pay/transfer, redeems points for products.
* **Cashier / Staff** (optional role): creates orders, attaches a customer, collects points, sends receipt by SMS.
* **Admin** (future): reporting & inventory; out of scope for v1 UI.

**Key use cases**

1. See balance and latest transactions.
2. Transfer points to a contact with quick amount presets.
3. Request/collect points: set an amount, show/share QR; payer scans and pays.
4. Shop with points: browse products → add to cart → confirm → pay with LBK Points → SMS receipt.
5. View full transaction history and details.
6. Manage basic profile & settings.

---

## 3) Scope & Feature List (mapped to UI flow)

### 3.1 Wallet Dashboard

* Display current **tier** (e.g., Gold) and **LBK balance**.
* Shortcuts: **Shopping**, **Transfer**, **Receive/QR**, **Profile/Settings**.
* **Recent activity** (± amounts, status); **View all** → history.

### 3.2 Transaction History

* Tabs or filters: **All / Earn / Spend / Transfer**.
* Transaction detail page (id, timestamp, type, counterparty, amount, status, notes).

### 3.3 Transfer Points

* Select recipient (search or recent list).
* Enter amount + **quick presets** (100 / 500 / 1000 / …).
* Validation (min/max, balance, daily cap).
* Confirm screen; show success/failure state and new balances.

### 3.4 Receive / Request via QR

* Enter requested amount.
* Generate **QR Code** that encodes pay-to account & amount.
* Options: **Save image / Share** (OS sheet).
* Pending state until paid; success state once settled.

### 3.5 Shopping (Catalog → Cart → Checkout)

* Product catalog: image, name, price (in points), **Add to cart**.
* Cart: items, qty, subtotal (points).
* **Attach customer** (name, phone) for the order.
* Checkout: confirm, **Pay with LBK Points**.
* Post-payment: order summary + **Send receipt via SMS** (required in v1).

### 3.6 Profile / Settings

* View/update basic profile fields (name, email, phone; editing rules per policy).
* Show membership tier & account identifiers.
* Link to history and legal (ToS/Privacy).

---

## 4) Functional Requirements

### 4.1 Authentication & Authorization

* Email/phone login (MVP). *(If SSO/Azure AD is required for staff backoffice, add as non-blocking integration.)*
* Roles: **member**, **staff** (optional), **admin** (future).

### 4.2 Points Ledger

* Double-entry style records for **earn**, **spend**, **transfer in/out**.
* Idempotent operations; every UI action produces a single immutable ledger entry (plus reversal only by admin tool, out of scope v1 UI).

### 4.3 Transfers & QR

* **Transfer**: balance check, lock, post, unlock; failure rolls back.
* **QR request**: generate payload `{recipient, amount, memo?, expiry}`; mark as paid when the payer completes.

### 4.4 Shopping & Orders

* Products with `id, name, description, price_points, active, stock?`.
* Cart computes subtotal; checkout debits points and creates `order` + `payment` record.
* **Receipt via SMS**: send order summary link + amount in points.

### 4.5 Notifications

* On success: in-app toast; optional SMS only for receipts (shopping flow).
* On failure: clear error copy (insufficient balance, invalid phone, expired QR).

### 4.6 Analytics / Audit

* Events: `view_dashboard`, `transfer_attempt/success/failure`, `qr_create`, `qr_paid`, `cart_add`, `checkout_success/failure`, `sms_sent/sms_failed`.
* Exportable transaction list (API).

---

## 5) Data Model (high level)

* **User**: `id, role, name, email, phone, tier, created_at`
* **Wallet**: `user_id, balance_points, updated_at`
* **Ledger**: `id, user_id, type[earn|spend|transfer_in|transfer_out], amount, counterparty_id?, order_id?, transfer_id?, status, created_at`
* **Transfer**: `id, from_user, to_user, amount, status[pending|posted|failed], created_at`
* **QRRequest**: `id, recipient_user, amount, expires_at, status[pending|paid|expired]`
* **Product**: `id, name, price_points, active, stock?`
* **Order**: `id, created_by_user, customer_name, customer_phone, total_points, status`
* **Receipt/SMS**: `id, order_id, phone, status, provider_ref, sent_at`

---

## 6) Integrations

* **SMS provider**: single vendor with webhook for delivery status.
* **(Optional) Identity**: SSO/Azure AD for staff.
* **(Optional) Payment QR spec**: internal format sufficient for v1 (not bank rails).

---

## 7) UX & Content Requirements

* Clear currency: **Points** (LBK) everywhere; no cash symbol.
* Presets for amounts on transfer/request.
* Empty states for dashboard, history, cart.
* Error copy examples:

  * “Insufficient balance. You have 1,240 LBK.”
  * “This QR has expired. Ask the requester to generate a new one.”
* Accessibility: keyboard usable, alt text for QR images.

---

## 8) Security & Compliance

* Server-side validation for all mutations.
* Rate limiting on transfers and QR generation.
* PII (phone, name) stored per privacy policy; redact in logs.
* All transport over TLS; CSRF protection for web.

---

## 9) Performance & Reliability

* P95 page load < 2.5s on 4G.
* Transfers/checkout commit < 2s or show pending with retry.
* Idempotency keys for all write endpoints.

---

## 10) Acceptance Criteria (MVP)

* Users can log in and see **accurate balance** and last 10 transactions.
* **Transfer** succeeds with valid inputs and fails gracefully on errors; ledger & balances reconcile.
* **QR Request** can be created, shared, paid, and marked **paid**; requester sees final state.
* **Shopping**: add ≥1 product, checkout with LBK Points, points deducted, order recorded, **SMS receipt** delivered to provided phone.
* **History**: filters show correct subsets; detail page matches ledger record.
* **Settings**: user can view profile; tier and identifiers visible.
* Unit tests cover core flows; basic design system applied; database schema migrated.

---

## 11) Release Plan

* **v1 (MVP):** Wallet dashboard, history, transfer, QR request, shopping with points, SMS receipt, basic settings.
* **v1.x:** Search in history, product categories, stock control, email receipts.
* **v2:** Refunds/voids, promotions, multi-tenant reporting.

---

## 12) Open Questions

1. Are transfers P2P only within tenant, or cross-tenant allowed?
2. Required QR standard (if any) and maximum amount?
3. SMS sender ID and copy templates approval?
4. Do we allow profile edits (name/phone) or lock to KYC data?

---

## 13) Success Metrics

* T+30d: **≥80%** of active users complete at least one core action (transfer, QR, or checkout).
* Transfer failure rate **<2%** (non-insufficient-funds).
* SMS delivery success **>98%**.
* Support tickets about balance/receipts reduced by **30%** vs baseline.

---

### Appendix: Screen-to-Feature Mapping (from provided flow)

* **Dashboard** → Wallet balance, shortcuts, recent activity.
* **Settings** → Profile/tier + link to history.
* **History** → “View all” list & detail.
* **Transfer** → Recipient select, amount presets, confirm.
* **Receive/QR** → Enter amount, generate/share QR, pay state.
* **Shopping** → Catalog → Cart → **Attach customer** → Checkout → **Pay with LBK Points** → **SMS receipt**.
