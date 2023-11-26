# Logging

## Logging strategy
1. User events logging will be done by the frontend so they can be user-oriented.
2. The backend will log system stats and all technical errors
3. The backend is responsible for logging
   - Expose a Logger to the frontend ?

## Errors
Technical errors must be logged in detail to make for easier debugging & quicker detection

## Events
Events to be logged :

### Backend
    - app startup & result
    - app closed
    - every error
    - every error ONCE

### Frontend
#### For each page
  - page visited
  - operation started
  - operation result
#### homepilot link opened


