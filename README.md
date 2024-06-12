# Documentation 
## Go Chat Application

Initialize project for make chat application use go gorilla websocket

## Prerequisites

1. Go 1.22.2
2. PostgreSQL
3. Make sure your go installation are configured properly

## Set-up

1. Configure Project

    ```sh
    # Run scripts to make env from .env-example
    make env

    # Run scripts to install used package
    make setup
    ```

### Run Development

```sh
# Run Service
make server
```

## How To Use This Application ?
1. You must register first via Postman (Postman documentation has been provided).
2. When registration is successful, please log in with the email and password that you have registered.

## How To Use Direct Message Feature ?
1. After successful login, you will successfully get the user ID which has been stored in the cookie. After successful login, you will successfully get the user ID which has been stored in the cookie. Then please visit the web page [Link Here](http://localhost:8000/page/direct-message/1?receiver_id=2). The param (number 1) shows your user ID, and the query param receiver ID shows the user ID you want to send to start a direct message.
2. Then you will be able to start a chat with the user you are targeting. and the intended user can also open this web page [Link Here](http://localhost:8000/page/direct-message/2?receiver_id=1).

## How To Use Chat Room Feature ?
1. First, you have to create a chat room first. Use postman to create a chat room.
2. After successfully creating a chat room, you can invite other users who you will enter into the chat room. hit room chat join endpoint on postman to invite users.
3. You can see a list of chat rooms that you have. Hit the chat room list endpoint on postman, then you will see a list of chat rooms that you have
4. You can also leave a group that you have if you want. Hit leave room chat endpoint on postman. Then, you will be successful in leaving the group.
5. If you have created a chat room and have invited several of your friends. You can start a chat on this page [Link Here](http://localhost:8000/page/room-chat/1?room_chat_id=5). The param (number 1) shows your user ID, and the query param room chat ID shows the room chat ID that you have created. And then you can start a chat in this chat room and your friends in the group can read the messages you made in the chat room at this link [Link Here](http://localhost:8000/page/room-chat/2?room_chat_id=5).
**NOTE**
In the url param. Please enter your registered user ID

postman collection has been attached to the documentation folder