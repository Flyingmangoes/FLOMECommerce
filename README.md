# Everything almost good

```
Need to implement this in order:
    1. Prevent overbought
    2. Reset password (Probably optional)
    3. REDIS caching 
    4. Frontend (React || Vue)
```

```
Implemented so far:
    1. BASIC CRUD system (such as repositories and controllers)
    2. Auth middleware
    3. CORS(untested since we don't have a web page yet)
    4. Rate Limiter
    5. Sudo system (Token for destructive request like update, delete, cancel order)
    6. Merchant User checker
    7. Stripe Payment System
    8. Emailing Feature using sendgrid API(For reset password and new account confirmation)
    9. Centralized Transaction Manager
    10. Automatic Error detection in middleware
    11. Custom logger that save the log into a *.log file using lumberjack.v2
    12. Database connection using PostgreSQL
    13. A password hashing service
    14. Email Template
    15. Store Search
    16. Pagination
    17. Product Search
    18. Prevent Insecure Direct Object Reference(IDOR) by define authorization for every type of user
    19. User Search
```