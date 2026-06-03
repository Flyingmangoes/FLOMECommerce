# Everything almost good

```
Need to implement this in order:
    1. Store search
    2. Prevent overbought. What does it mean?, if a user is checkouting an order we need to make sure
        the user only pay it one time, because if the server happen to be lagging the user might accidentaly 
        pay twice and that is bad for user. So we need to make sure every transaction only happen once.
    3. Product search & list
    4. Pagination using cursor and limit, even thought i have pagination system created already, i still 
        don't understand how it works and how to implement it to the code base so i need to work on it 
        once again (lol bro using ai)
    5  User search
    6. Reset password (Probably optional)
    7. Email template
    8. REDIS caching
    9. Frontend (React || Vue || Kotlin)
    10. Prevent Insecure Direct Object Reference(IDOR) by define authorization for every type of user
```

```
Implemented so far:
    1. BASIC CRUD system (such as repositories and controllers)
    2. Auth middleware
    3. CORS (untested since we don't have a web page yet)
    4. Rate Limiter that block an ip that spamming the request
    5. Sudo system (Basically an confirmation token for request like update, delete, etc)
    6. User type checker middleware to separate a normal user from a store owner user
        the difference is in the authorization both user have.
    7. Payment system using stripe (since this is scalable i can increase more payment feature in the future)
    8. Emailing Feature using sendgrid API(For reset password and new account confirmation)
    9. Centralized Tx manager (or just transaction wrapper) this one is for wrapping a function that required
        the process to either complete 100% or fail 100% no in between such as place order in the service that
        required multiple table accessed and modified.
    10. Automatic Error detection in middleware
    11. Custom logger that save the log into a *.log file using lumberjack.v2
    12. Database connection using PostgreSQL
    13. A password hashing service
```