“Real-time Study Room” (WebSocket chat + presence + file drop)

A web app where people join a room with a code and can:

chat in real time

see who’s currently online (“presence”)

optionally send small files or links

optionally do “ping” / “raise hand” / reactions

STEPS:

**Run the Go backend **
'''
steps:
    cd server
    go mod download
    go run .
'''

**Run the TypeScript client (in a second terminal) **
'''
steps:
    cd web
    npm install
    npm run dev
'''
