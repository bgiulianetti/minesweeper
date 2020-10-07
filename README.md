## Minesweeper API

## Decisions made
- The project was developed in Golang, with Go modules and [Gin](https://github.com/gin-gonic/gin).
- I decided to use in memory persistence, and also developed a fully MongoDB service persistance as well, with a MongoDB server up and running of my own. I used the [MongoDB Golang Driver](https://github.com/mongodb/mongo-go-driver)
- I decided to deploy the API in [Heroku.com](https://heroku.com), the main reason for this decision is that I have never used Heroku.com before and I was curious of what this platform as a service was about. I use Amazon AWS for all of my projects, but I wanted to use something differente just to give it a try and build more experience with other tools an services.

## List of Endpoints

- Base URL: `https://bgiulianetti-minesweeper.herokuapp.com/minesweeper`

### Get all games
- Path: `/games`
- Rest verb: GET
- Responses:
  - 500: Internal Server Error
    ```
    {
      "message": "Internal Server Error",
      "error": "internal_server_error",
      "status": 500
    }
    ```
  - 200: Lists of all the games of all the users
  ```
  [
    {
        "games": [
            {
                "game_id": 1602094036617367,
                "rows": 2,
                "columns": 2,
                "mines": 2,
                "start_time": "2020-10-07T18:07:16.617367339Z",
                "finish_time": "0001-01-01T00:00:00Z",
                "cells_revealed": 0,
                "status": "on_going",
                "board": [
                    [
                        {
                            "is_revealed": false,
                            "has_mine": false,
                            "sourrounded_by": 2,
                            "flag": ""
                        },
                        {
                            "is_revealed": false,
                            "has_mine": true,
                            "sourrounded_by": 2,
                            "flag": ""
                        }
                    ],
                    [
                        {
                            "is_revealed": false,
                            "has_mine": false,
                            "sourrounded_by": 2,
                            "flag": ""
                        },
                        {
                            "is_revealed": false,
                            "has_mine": true,
                            "sourrounded_by": 2,
                            "flag": ""
                        }
                    ]
                ]
            }
        ],
        "user_id": "bruno"
    },
    {
        "games": [
            {
                "game_id": 1602094100929445,
                "rows": 2,
                "columns": 2,
                "mines": 2,
                "start_time": "2020-10-07T18:08:20.929445615Z",
                "finish_time": "0001-01-01T00:00:00Z",
                "cells_revealed": 0,
                "status": "on_going",
                "board": [
                    [
                        {
                            "is_revealed": false,
                            "has_mine": false,
                            "sourrounded_by": 1,
                            "flag": ""
                        },
                        {
                            "is_revealed": false,
                            "has_mine": true,
                            "sourrounded_by": 1,
                            "flag": ""
                        }
                    ],
                    [
                        {
                            "is_revealed": false,
                            "has_mine": false,
                            "sourrounded_by": 1,
                            "flag": ""
                        },
                        {
                            "is_revealed": false,
                            "has_mine": false,
                            "sourrounded_by": 1,
                            "flag": ""
                        }
                    ]
                ]
            }
        ],
        "user_id": "federico"
    }
  ]```
### Create a new game
- Path: `/users/{username}/games`
- Rest verb: POST
- Request:
```
{
    "rows" : 2,
    "columns" : 2,
    "mines" : 1
}

```
- Responses:
  - 400: Bad Request
    - Rows and columns don't match
    ```
    {
        "message": "rows and columns must be equals",
        "error": "bad_request",
        "status": 400
    }
    ```
    - None or too many Mines
    ```
    {
      "message": "the number of mines must be at least one, and less or equal than total of cells in the game",
      "error": "bad_request",
      "status": 400
    }
    ```
    - Columns or rows greater than 30
    ```
    {
    "message": "columns must be greater than 0 and less or equal than 30",
      "error": "bad_request",
      "status": 400
      }
    ```
  - 201: Game Created
  ```
  {
    "game_id": 1602095465769850,
    "rows": 2,
    "columns": 2,
    "mines": 1,
    "start_time": "2020-10-07T18:31:05.769850542Z",
    "finish_time": "0001-01-01T00:00:00Z",
    "cells_revealed": 0,
    "status": "on_going",
    "board": [
        [
            {
                "is_revealed": false,
                "has_mine": true,
                "sourrounded_by": 1,
                "flag": ""
            },
            {
                "is_revealed": false,
                "has_mine": false,
                "sourrounded_by": 1,
                "flag": ""
            }
        ],
        [
            {
                "is_revealed": false,
                "has_mine": false,
                "sourrounded_by": 1,
                "flag": ""
            },
            {
                "is_revealed": false,
                "has_mine": false,
                "sourrounded_by": 1,
                "flag": ""
            }
        ]
    ]
  }```

### Reveal a cell
- Path: `/users/{username}/games/{gameid}/reveal`
- Rest verb: POST
- Request:
```
{
  "row" : 1,
  "column": 0
}
```
- Responses:
  - 400: Bad Request
    - Rows or column out of boundries
    ```
    {
      "message": "flag out of boundries (columns exceeded)",
      "error": "bad_request",
      "status": 400
    }
    ```
  - 200: Cell Revealed
  ```
    Same response as Game Created but with the given cell revealed
  ```
### Flag a cell
- Path: `/users/{username}/games/{gameid}/flag`
- Rest verb: POST
- Request:
```
{
  "row" : 1,
  "column": 0,
  "flag": "red_flag"
}
```
- Responses:
  - 400: Bad Request
    - Missing or invalid flag parameter
    ```
    {
      "message": "Available flag options: [question_mark, [red_flag]",
      "error": "bad_request",
      "status": 400
    }
    ```
    - Flag position out of boundries
    ```
    {
      "message": "flag out of boundries (columns exceeded)",
      "error": "out_of_boundries",
      "status": 400
    }
    ```
  - 200: Cell Flagged
  ```
    Same response as Game Created but with the given cell flagged (or unflagged)
  ```

### Win or Lose
If in some point a cell with a mine is revealed or all the blank cells in the board are revealed and all the mines are flagged, the game will change its status to "WON" or "LOSE", and it will populate the finish date and hour, and you won't be able to make any more changes to that game, otherwise the game status will be "on going"

### Additional endpoints

I Created some endpoints that in my opinion helped me to develop the API and validate its functionality and behavior. 
The endpoints are:

- Show Board solution: It will show the board solution in a matrix look, formatted in plain text.
  ```
    GET /users/{username}/games/{gameid}/solution
  ```
- Show Board status: It will show the board status in a matrix look, formatted in plain text. 
  ```
    GET /users/{username}/games/{gameid}/status
  ```
  - The cells not yet revealed will bw shown as squares
  - The cells revealed without mines around will be shown as underscore
  - The cells revealed with mines around will be shown with a number, indicating the numbers of mines around
  - The cells revealed with mines will be shown with '*'
- Get all games from a user
  ```
    GET /users/{username}/games
  ```
- Get a single game from a user
  ```
    GET /users/{username}/games/{gameid}
  ```
  - Delete all games
  ```
    DELETE /games
  ```


