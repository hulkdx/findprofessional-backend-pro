# Temporary Notes on how to create payments:
How to design idempotency for payments


- client will call pro microservice for booking/payments
- booking_holds record is created
- payment microservice is called
- payment microservice calls stripe microservice
- payment microservice return responses back to pro microservice with payment_intent_id
- final bookings record is created inside pro microservice
- final response is back to the client

mermaid flowchart (WIP):
```mermaid
sequenceDiagram
    autonumber
    participant C as Client (App)
    participant Pro as Pro Microservice
    participant DB as Postgres
    participant Pay as Payments MS
    participant S as Stripe MS

    C->>Pro: POST /professional/{id}/booking<br>idempotency_key=uuid-v7
    Pro->>DB: INSERT booking_holds(user_id, created_at, expires_at)
    DB-->>Pro: hold_id (fail if already held / expired)

    Pro->>Pay: POST /payments/init {hold_id, IK}
    Pay->>S: Create PaymentIntent (idempotency_key = "book_"+IK)
    S-->>Pay: payment_intent_id, status=requires_confirmation
    Pay-->>Pro: payment_intent_id, client_secret

    Note over C,Pro: Client confirms payment via SDK/redirect

    C->>Pro: POST /bookings/finalize {hold_id, IK}
    Pro->>DB: Validate hold (exists, not expired, matches user/slot)

    rect rgb(245,245,245)
    Note over Pro,DB: Idempotent finalize
    Pro->>DB: INSERT INTO bookings(..., hold_id, created_at, updated_at)<br>ON CONFLICT (hold_id) DO UPDATE SET updated_at = EXCLUDED.updated_at<br>RETURNING booking_id
    DB-->>Pro: booking_id
    Pro->>DB: INSERT INTO booking_professional_availability(availability_id, booking_id, created_at)
    DB-->>Pro: ok
    end

    Pro-->>C: 200 OK {booking_id, payment_intent_id, status}
```
