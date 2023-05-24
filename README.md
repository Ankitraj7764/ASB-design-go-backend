# ASB-Design-Go
## CreateChallenge

The http request body contains all details of the challenge
- challenge name
- problem statement / description  
- image (in url form)
- score 
- next challenge (in url form)
- previous challenge (in url form)
- difficulty (easy/medium/hard)

This request body is decoded and all the data is saved in `challengeRequestBody`.

This `challengeRequestBody` is saved into the database in the `challenges` collection.

## SubmitSolution

The http request body contains
- user email,
- solution (which is in form of a string)
- corresponding problem name
This request body is decoded and saved to a variable `submittedSolutionRequestBody`. 

Then it save the `submittedSolutionRequestBody` to the `unsettled` collection.

## AcceptSolution

The http request contains
- user email
- solution url
- problem name 

This request body is decoded and saved in `pendingProblemRequestBody`.

Then we extract user using username.

Create a new variable named `solvedChallengesArr` which contains all the challenges that the user has solved.

Then we check if the user has already solved the problem or not.

If the problem is already solved by that user then we return `http.StatusAlreadyReported`.

If the problem is not solved we append the challenge in the `solvedChallengesArr`.

## Forgot Password

The request body only contains the email that has been entered by the user who is claiming that he / she forgot the password.
This request body is decoded and stored inside `req`.
`Receiver` variable contains the email that has been extracted from `req`.

Then we create a random code (6 digits).

We send this random code to the entered email id by calling the function `SendMail`.

Now in the database we set the recovery code as this random code.

## Reset Password

The http request for resetting password contain three parameters
- email
- code
- new password

First we check if the code is same as the code that has been sent to the user email by using the function `CompareCodes`.

Now if the code matches then we checek if the code has been sent to his email id within the speculated time( 15 minutes in our case ) or not.

If the code has been sent within 15 minutes then we update the password for the corresponding user with the new password.

## Send Mail

This function is used to send recovery code to the entered email address.
First it initializes a new gomail.

In the header it set the sender as the email id from which the message will be sent.

Corresponding to `From` it set the receiver email address.

"Subject" key contains the subject or heading of the mail.

`SetBody` function sets the body message which contains recovery code for that user.
Then it set SMTP server.
Then we send the email.
## Log in

The request body contains the email id and password entered by the user while logging in.

These responses are decoded and stored inside the variable named `credentials`.

Then we find the actual password by calling the function `FindPassword` and passing the email id of user. 

This function find the coresponding password of the user from the database.

As this password is in encrypted form we compare the passwords by calling the function `ComparePasswords`.

Now we generate access token and encrypted refresh token and return them as response.

Now this tokens should be saved in the cookie or local storage throughout the session.

## Assign Auth Tokens

First we create a new access token by calling the function `CreateNewJWTAccessToken` which takes email id as a parameter and generate an access token.

We also create a new refresh token by calling the function `CreateNewReferenceToken` and passing length of the token.

As the user is loggin in with their email id and password we can say that the user is not suspicious. 

So we don't have to keep record of any older tokens so we create and empty array of refresh tokens `newRefreshTokenArray`.
 Now we append the newly generated refresh token to the array.

Now we assign this `newRefreshTokenArray` corresponding the refreshtokens field for that particular user.

As this access token and refresh tokens can be used through out the session so we send the tokens as a response for storing them in the local storage.

## Log out 

We extract the email id from JWT token and corresponding to that email id we assign an empty array as refresh tokens for that user.




# Database Design

## user
```
type User struct {
	Id                      primitive.ObjectID `json:"id,omitempty"`
	Name                    string             `json:"name,omitempty" validate:"required"`
	EmailId                 string             `json:"email-id,omitempty" validate:"required" `
	Password                string             `json:"password,omitempty" validate:"required"`
	ProfilePictureURL       string             `json:"profile-picture-url,omitempty" validate:"required"`
	Score                   int                `json:"score,omitempty" validate:"required"`
	RecoveryCode            string             `json:"recovery-code"`
	CodeSendingTime         time.Time          `json:"code-sending-time"`
	SolvedChallenges        []Challange        `json:"solved-challenges"`
	RefreshTokens           []string           `json:"refresh-tokens"`
	CurrentAccessToken      string             `json:"current-access-token"`
	AccessTokenSendingTime  time.Time          `json:"access-token-sending-time"`
	RefreshTokenSendingTime time.Time          `json:"refresh-token-sending-time"`
	// SubmittedSolutions  []string
	// SubmittedChallenges []Challange
}
```
## Modules
```
type Modules struct {
	AllModules []Module
}
```

## Module
```
type Module struct {
	ModuleName string `json:"module-name,omitempty" validate:"required"`
	LevelList  []Level
}
```
## Level
 
```
type Level struct {
	LevelName    string `json:"level-name,omitempty" validate:"required"`
	CategoryList []Categories
}
```

## Categories

```
type Categories struct {
	CategoryName string `json:"category-name,omitempty" validate:"required"`
	ProblemList  []Challange
}
```

## Challenge

```
type Challange struct {
	ChallengeName        string `json:"challenge-name,omitempty" validate:"required"`
	Description          string `json:"description,omitempty" validate:"required"`
	ImageUrl             string `json:"image-url,omitempty" validate:"required"`
	Score                int    `json:"score-assigned,omitempty" validate:"required"`
	NextChallengeUrl     string `json:"next-challenge-url,omitempty" validate:"required"`
	PreviousChallengeUrl string `json:"previous-challenge-url,omitempty" validate:"required"`
	Difficulty           string `json:"difficulty,omitempty" validate:"required"`
	// easy/medium/hard
}
```
# Working
## localhost:8080/user
![image](https://user-images.githubusercontent.com/84634405/236602469-0b90f6d7-5b01-47a5-8abb-61b8b046e43d.png)
### output
```json
{
    "status": 201,
    "message": "success",
    "data": {
        "data": {
            "InsertedID": "6455e66a2e6160eaa35f575b"
        }
    }
}
```
![image](https://user-images.githubusercontent.com/84634405/236602489-356acf50-6e7b-43cd-b950-ceb427993d80.png)
## localhost:8080/login
### output
### using correct password ✅
![image](https://user-images.githubusercontent.com/84634405/236602621-a5560a9a-559f-4a0e-8420-c1ee5a1681eb.png)

```json
{
    "status": 200,
    "message": "success",
    "data": {
        "data": {
            "access-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODMzODc1NDUsImN1c3RvbS1jbGFpbXMiOnsiZW1haWwiOiJsb2dhbmRoYW5hbmpveUBnbWFpbC5jb20ifX0.koq0365HK8uI7yeqSPdP8p-fMIjkS1PxN5G00seE9Ek",
            "refresh-token": "4315652023-05-06 11:09:05"
        }
    }
}
```
### using wrong password ❌
![image](https://user-images.githubusercontent.com/84634405/236602727-0f2fca8a-ab27-4b47-bfb4-f54aeb332586.png)

## localhost:8080/log-out
![image](https://user-images.githubusercontent.com/84634405/236603149-27a76973-bc46-447e-915d-66c98a6f4ce1.png)

## localhost:8080/home
![image](https://user-images.githubusercontent.com/84634405/236603188-6b7eee0b-f902-469b-b102-2a6e71c4562e.png)

## localhost:8080/forgot-password
![image](https://user-images.githubusercontent.com/84634405/236603229-605f27e4-a8cb-4f94-ac3f-bdb11d61fa5c.png)
![image](https://user-images.githubusercontent.com/84634405/236603250-0164ca69-e7f5-44c0-811a-a63be84c7d33.png)

## localhost:8080/reset-password
![image](https://user-images.githubusercontent.com/84634405/236603302-ba7de814-ece2-4998-94cc-0b120cd7b0f3.png)
## localhost:8080/module
![image](https://user-images.githubusercontent.com/84634405/236603424-6edff545-d0f6-4ec2-9429-eb8c65e9b944.png)
![image](https://user-images.githubusercontent.com/84634405/236603541-d395bfcc-4021-420a-9335-cf1064d76a85.png)

## localhost:8080/level
![image](https://user-images.githubusercontent.com/84634405/236603484-7c938034-55f7-4675-a7c0-f764afaee1d3.png)
![image](https://user-images.githubusercontent.com/84634405/236603544-f438f87d-1908-439d-80b7-8a62dd068ffc.png)

## localhost:8080/add-challenge
![image](https://user-images.githubusercontent.com/84634405/236604972-84412245-47af-41e7-bde4-64fa691384c5.png)
![image](https://user-images.githubusercontent.com/84634405/236605024-191cbe51-4923-400f-8738-2bccdb0641d6.png)
## localhost:8080/category
![image](https://user-images.githubusercontent.com/84634405/236605700-f0fa1ee8-b463-4ce6-8720-e1cab1dbfe2d.png)
![image](https://user-images.githubusercontent.com/84634405/236605724-7e8ab6d0-28d1-4441-8101-ecc0bab01266.png)

## localhost:8080/user/6455e66a2e6160eaa35f575a
![image](https://user-images.githubusercontent.com/84634405/236605785-8ebef039-0825-4c1a-a3d5-2483953ea395.png)
![image](https://user-images.githubusercontent.com/84634405/236605789-e9914bb8-e583-4e59-8f69-90323cabfbb9.png)
## localhost:8080/get-all-users
![image](https://user-images.githubusercontent.com/84634405/236605843-e09ec8ba-f781-43a8-ab61-56d94678f900.png)
```json
{
    "status": 200,
    "message": "success",
    "data": {
        "data": [
            {
                "_id": "643563bb8f02a21c66f90d49",
                "codesendingtime": "2023-04-11T19:12:19.08+05:30",
                "emailid": "abckd1wds23gmsasdsail.com",
                "id": "643563bb8f02a21c66f90d48",
                "name": "Dhananjoy Dey test",
                "password": "$2a$04$JFtSy3LRS1VsAS4Ya8VWEeZzN.oQioYnId9ChempXJ1TS.DwMKTG2",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "score": 123,
                "solvedchallenges": null,
                "submittedchallenges": null,
                "submittedsolutions": null
            },
            {
                "_id": "643563c08f02a21c66f90d4b",
                "codesendingtime": "2023-04-11T19:12:24.872+05:30",
                "emailid": "abckd1wds23gmsasdsail.com",
                "id": "643563c08f02a21c66f90d4a",
                "name": "Dhananjoy Dey test",
                "password": "$2a$04$2aEB5WkCa9JlB60XUVdf/eWlI9TBAMfh0kBKWJjwGEnXWSIDrgnYq",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "score": 123,
                "solvedchallenges": null,
                "submittedchallenges": null,
                "submittedsolutions": null
            },
            {
                "_id": "643952a6e313c9c60291453d",
                "codesendingtime": "2023-04-14T18:48:30.143+05:30",
                "emailid": "abckd1wds23gmsasdsail.com",
                "id": "643952a6e313c9c60291453c",
                "name": "Dhananjoy Dey test",
                "password": "$2a$04$x1vDPrbCSWkxg6jlH1at3uEM5oQm/6TamGUbJOu48sWC44Gl5cEeq",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "score": 0,
                "solvedchallenges": null,
                "submittedchallenges": null,
                "submittedsolutions": null
            },
            {
                "_id": "6447ef4c56580ec192facbb7",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-04-25T20:52:40.007+05:30",
                "currentaccesstoken": "",
                "emailid": "shreyasidgp2@gmail.com",
                "id": "6447ef4c56580ec192facbb6",
                "name": "pupai",
                "password": "$2a$04$.lqieM6dd3rfHbT/J28BwOA/qDe9zHS9S1R1ag/Sqk/bYJZu//BOe",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "763721",
                "refreshtokens": [
                    "$2a$04$NjGns2KB0/TaEGai8lK6GeY76kfUAGcyoZznSagEYIPHJlVfNKJCu"
                ],
                "refreshtokensendingtime": "2023-04-25T20:53:45.625+05:30",
                "score": 0,
                "solvedchallenges": null
            },
            {
                "_id": "644a1a0af969e8119b156a06",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-05-06T11:24:27.616+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "644a1a0af969e8119b156a05",
                "name": "logan",
                "password": "$2a$04$PfhYHTBvP7zC3jndYJb9c.HELqDy.z1ELSSrPi2uhNE2KGke25GZS",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "988938",
                "refreshtokens": [
                    "$2a$04$IplPqLbUWyOTeP4B2oYkN.PqjNvayD/dBnp6GfsvBeC1V3be9urju"
                ],
                "refreshtokensendingtime": "2023-05-06T11:18:28.602+05:30",
                "score": 0,
                "solvedchallenges": null
            },
            {
                "_id": "644cf0d2a244ae1aa50c97c2",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-04-29T15:56:26.007+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "644cf0d2a244ae1aa50c97c1",
                "name": "logan",
                "password": "$2a$04$BBYWbJcVeT7kEDPSjwpgWen8AEbpAMJfpr0sT4cUuJGL6EyiGpm1W",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 0,
                "solvedchallenges": null
            },
            {
                "_id": "644cf19d0d57d43b7d84967d",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-04-29T15:59:49.3+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "644cf19d0d57d43b7d84967c",
                "name": "logan",
                "password": "$2a$04$GRxXw/zA2w0o5CDRWtwjQ.KHiHWS15khGMqkN58IamckomHig8qKO",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 0,
                "solvedchallenges": null
            },
            {
                "_id": "644cf21c21f9eb6095d97f60",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-04-29T16:01:56.835+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "644cf21c21f9eb6095d97f5f",
                "name": "logan",
                "password": "$2a$04$eNFPXK4pMi74A3YEBjtA0ufQvmJK.a6gQxkQK19xAZD6mpK3bYKea",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 0,
                "solvedchallenges": null
            },
            {
                "_id": "644cf22321f9eb6095d97f62",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-04-29T16:02:03.907+05:30",
                "currentaccesstoken": "",
                "emailid": "",
                "id": "644cf22321f9eb6095d97f61",
                "name": "DHANANJOY DEY",
                "password": "$2a$04$OQKz4CjV1AwabK92ir9qDu9FJxqcty4zsgq.nKdHaAhlGWxACKW/C",
                "profilepictureurl": "",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 0,
                "solvedchallenges": null
            },
            {
                "_id": "64533de996b168010ac33953",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-05-04T10:38:57.617+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "64533de996b168010ac33952",
                "name": "logan",
                "password": "$2a$04$6IAwb2yYeMkn.QkXzdS7mufgHd1yXF/nxfNPn5EpK0WEp6QWInFFm",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 23,
                "solvedchallenges": null
            },
            {
                "_id": "64533df796b168010ac33955",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-05-04T10:39:11.322+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "64533df796b168010ac33954",
                "name": "logan",
                "password": "$2a$04$ITkLaY9ScqOSb3EsgmfemefC4Xe7oKXW/6SToYaoQx2qkoUeXx4tG",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 124,
                "solvedchallenges": null
            },
            {
                "_id": "64533dff96b168010ac33957",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-05-04T10:39:19.07+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "64533dff96b168010ac33956",
                "name": "logan",
                "password": "$2a$04$QgG.AzAYIz6fEloj6Xn8m.IjZJmuVX1yTIZEKRM6vKr..5ZPkNSgu",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 125,
                "solvedchallenges": null
            },
            {
                "_id": "6455e66a2e6160eaa35f575b",
                "accesstokensendingtime": "0001-01-01T05:30:00+05:30",
                "codesendingtime": "2023-05-06T11:02:26.991+05:30",
                "currentaccesstoken": "",
                "emailid": "logandhananjoy@gmail.com",
                "id": "6455e66a2e6160eaa35f575a",
                "name": "sample name",
                "password": "$2a$04$s4AZQ9cgv7vWDWklgATBtuGtAwsCDYa9ICQigtJ5IK.993ZrW/5wy",
                "profilepictureurl": "random-profile-pic-url",
                "recoverycode": "-123",
                "refreshtokens": null,
                "refreshtokensendingtime": "0001-01-01T05:30:00+05:30",
                "score": 0,
                "solvedchallenges": null
            }
        ]
    }
}
```
## localhost:8080/get-leaderboard-details
![image](https://user-images.githubusercontent.com/84634405/236605889-80956a44-833f-4cc3-b43c-4021433fff6c.png)

```json
{
    "status": 200,
    "message": "success",
    "data": {
        "data": [
            {
                "id": "64533dff96b168010ac33956",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$QgG.AzAYIz6fEloj6Xn8m.IjZJmuVX1yTIZEKRM6vKr..5ZPkNSgu",
                "profile-picture-url": "random-profile-pic-url",
                "score": 125,
                "recovery-code": "-123",
                "code-sending-time": "2023-05-04T05:09:19.07Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "64533df796b168010ac33954",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$ITkLaY9ScqOSb3EsgmfemefC4Xe7oKXW/6SToYaoQx2qkoUeXx4tG",
                "profile-picture-url": "random-profile-pic-url",
                "score": 124,
                "recovery-code": "-123",
                "code-sending-time": "2023-05-04T05:09:11.322Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "643563bb8f02a21c66f90d48",
                "name": "Dhananjoy Dey test",
                "email-id": "abckd1wds23gmsasdsail.com",
                "password": "$2a$04$JFtSy3LRS1VsAS4Ya8VWEeZzN.oQioYnId9ChempXJ1TS.DwMKTG2",
                "profile-picture-url": "random-profile-pic-url",
                "score": 123,
                "recovery-code": "-123",
                "code-sending-time": "2023-04-11T13:42:19.08Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "643563c08f02a21c66f90d4a",
                "name": "Dhananjoy Dey test",
                "email-id": "abckd1wds23gmsasdsail.com",
                "password": "$2a$04$2aEB5WkCa9JlB60XUVdf/eWlI9TBAMfh0kBKWJjwGEnXWSIDrgnYq",
                "profile-picture-url": "random-profile-pic-url",
                "score": 123,
                "recovery-code": "-123",
                "code-sending-time": "2023-04-11T13:42:24.872Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "64533de996b168010ac33952",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$6IAwb2yYeMkn.QkXzdS7mufgHd1yXF/nxfNPn5EpK0WEp6QWInFFm",
                "profile-picture-url": "random-profile-pic-url",
                "score": 23,
                "recovery-code": "-123",
                "code-sending-time": "2023-05-04T05:08:57.617Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "643952a6e313c9c60291453c",
                "name": "Dhananjoy Dey test",
                "email-id": "abckd1wds23gmsasdsail.com",
                "password": "$2a$04$x1vDPrbCSWkxg6jlH1at3uEM5oQm/6TamGUbJOu48sWC44Gl5cEeq",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "-123",
                "code-sending-time": "2023-04-14T13:18:30.143Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "6447ef4c56580ec192facbb6",
                "name": "pupai",
                "email-id": "shreyasidgp2@gmail.com",
                "password": "$2a$04$.lqieM6dd3rfHbT/J28BwOA/qDe9zHS9S1R1ag/Sqk/bYJZu//BOe",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "763721",
                "code-sending-time": "2023-04-25T15:22:40.007Z",
                "solved-challenges": null,
                "refresh-tokens": [
                    "$2a$04$NjGns2KB0/TaEGai8lK6GeY76kfUAGcyoZznSagEYIPHJlVfNKJCu"
                ],
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "2023-04-25T15:23:45.625Z"
            },
            {
                "id": "644a1a0af969e8119b156a05",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$PfhYHTBvP7zC3jndYJb9c.HELqDy.z1ELSSrPi2uhNE2KGke25GZS",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "988938",
                "code-sending-time": "2023-05-06T05:54:27.616Z",
                "solved-challenges": null,
                "refresh-tokens": [
                    "$2a$04$IplPqLbUWyOTeP4B2oYkN.PqjNvayD/dBnp6GfsvBeC1V3be9urju"
                ],
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "2023-05-06T05:48:28.602Z"
            },
            {
                "id": "644cf0d2a244ae1aa50c97c1",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$BBYWbJcVeT7kEDPSjwpgWen8AEbpAMJfpr0sT4cUuJGL6EyiGpm1W",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "-123",
                "code-sending-time": "2023-04-29T10:26:26.007Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "644cf19d0d57d43b7d84967c",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$GRxXw/zA2w0o5CDRWtwjQ.KHiHWS15khGMqkN58IamckomHig8qKO",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "-123",
                "code-sending-time": "2023-04-29T10:29:49.3Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "644cf21c21f9eb6095d97f5f",
                "name": "logan",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$eNFPXK4pMi74A3YEBjtA0ufQvmJK.a6gQxkQK19xAZD6mpK3bYKea",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "-123",
                "code-sending-time": "2023-04-29T10:31:56.835Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "644cf22321f9eb6095d97f61",
                "name": "DHANANJOY DEY",
                "password": "$2a$04$OQKz4CjV1AwabK92ir9qDu9FJxqcty4zsgq.nKdHaAhlGWxACKW/C",
                "recovery-code": "-123",
                "code-sending-time": "2023-04-29T10:32:03.907Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            },
            {
                "id": "6455e66a2e6160eaa35f575a",
                "name": "sample name",
                "email-id": "logandhananjoy@gmail.com",
                "password": "$2a$04$s4AZQ9cgv7vWDWklgATBtuGtAwsCDYa9ICQigtJ5IK.993ZrW/5wy",
                "profile-picture-url": "random-profile-pic-url",
                "recovery-code": "-123",
                "code-sending-time": "2023-05-06T05:32:26.991Z",
                "solved-challenges": null,
                "refresh-tokens": null,
                "current-access-token": "",
                "access-token-sending-time": "0001-01-01T00:00:00Z",
                "refresh-token-sending-time": "0001-01-01T00:00:00Z"
            }
        ]
    }
}
```

## localhost:8080/find-module
![image](https://user-images.githubusercontent.com/84634405/236607369-7330c872-94f4-4d36-8c8a-6573ed118dbe.png)

## localhost:8080/find-level
![image](https://user-images.githubusercontent.com/84634405/236607446-f9c54a12-a2f2-480e-a10a-3a5f29ac8fbb.png)

## localhost:8080/get-all-levels
![image](https://user-images.githubusercontent.com/84634405/236607513-51ddf887-4a8d-4f86-ad30-1cc6f2d4f697.png)

## localhost:8080/get-all-challenges
![image](https://user-images.githubusercontent.com/84634405/236607554-108ae2e6-97ae-47ec-90a3-b6455736fe6a.png)

## localhost:8080/get-all-categories
![image](https://user-images.githubusercontent.com/84634405/236607604-65bd920f-6ad5-4159-9ecb-022f74a83094.png)

## localhost:8080/find-category
![image](https://user-images.githubusercontent.com/84634405/236607705-42550776-00cd-4f81-991e-de60bc4b3db6.png)

## localhost:8080/add-level-in-module
![image](https://user-images.githubusercontent.com/84634405/236607906-3fbc341c-e7ad-428f-af64-74a32be177de.png)

## localhost:8080/delete-level
![image](https://user-images.githubusercontent.com/84634405/236608064-1fa3bd89-51bb-4c7e-bb8e-9f3d110f175d.png)
![image](https://user-images.githubusercontent.com/84634405/236608067-e475c0ab-4a1b-426d-b544-17f519de7936.png)

## localhost:8080/delete-module
![image](https://user-images.githubusercontent.com/84634405/236608146-0dde8edc-f6ba-4ce9-8398-047d515ab368.png)
![image](https://user-images.githubusercontent.com/84634405/236608143-310d603f-f196-4c12-bcd4-815db5f092d2.png)

## localhost:8080/delete-category
![image](https://user-images.githubusercontent.com/84634405/236608273-a8eecef5-6491-4aaf-8226-5455ce91c567.png)

![image](https://user-images.githubusercontent.com/84634405/236608275-4cc6d68e-738b-49c7-b9a3-9137423c2b36.png)

## localhost:8080/submit-solution
![image](https://user-images.githubusercontent.com/84634405/236608430-75de495b-3ee9-493d-854f-e15bdc2d2cf3.png)
![image](https://user-images.githubusercontent.com/84634405/236608427-7da1e045-4069-4e79-a726-6029766909e9.png)

## localhost:8080/accept-solution

![image](https://user-images.githubusercontent.com/84634405/236609620-e5aac521-26aa-49a0-983a-ba2367c39ba3.png)
![image](https://user-images.githubusercontent.com/84634405/236609619-596cb12c-0ee7-4502-b8b1-f67cc93ff885.png)


```json
{
  "_id": {
    "$oid": "6455ffb5a9f13bc7256c26a4"
  },
  "id": {
    "$oid": "6455ffb5a9f13bc7256c26a3"
  },
  "name": "sample name 12:43",
  "emailid": "logandhananjoy1250@gmail.com",
  "password": "$2a$04$gU93uChQjnsiORZnp.5MQ.ZOvPXXqiPdjGEykp2RcX.zGRqGdpkqm",
  "profilepictureurl": "random-profile-pic-url",
  "score": 123,
  "recoverycode": "-123",
  "codesendingtime": {
    "$date": "2023-05-06T07:20:21.071Z"
  },
  "solvedchallenges": [
    {
      "challengename": "2-sum",
      "description": "find sum of 2 element",
      "imageurl": "test-url",
      "score": 123,
      "nextchallengeurl": "test-next",
      "previouschallengeurl": "test-previous",
      "difficulty": "medium"
    }
  ],
  "refreshtokens": null,
  "currentaccesstoken": "",
  "accesstokensendingtime": {
    "$date": {
      "$numberLong": "-62135596800000"
    }
  },
  "refreshtokensendingtime": {
    "$date": {
      "$numberLong": "-62135596800000"
    }
  }
}
```
## localhost:8080/get-all-submissions
![image](https://user-images.githubusercontent.com/84634405/236888664-e6787c00-6805-4564-8a84-99bbacbcf3ec.png)
```json
{
    "status": 200,
    "message": "success",
    "data": {
        "data": [
            {
                "_id": "64592d4134ae7d94ad03a945",
                "problemname": "test_challenge_1",
                "solution": "http://drive.google.com/fjiknf",
                "useremail": "logandhananjoy855gmail.com"
            },
            {
                "_id": "64592d6434ae7d94ad03a946",
                "problemname": "test_challenge_1",
                "solution": "http://drive.google.com/fjiknf",
                "useremail": "logandhananjoy854gmail.com"
            }
        ]
    }
}
```
## localhost:8080/get-all-modules
![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/1b4e64d5-4d51-4f6e-a778-9b6844f6f13d)

```
{
    "status": 200,
    "message": "success",
    "data": {
        "data": [
            {
                "_id": "64142b1a3b98e7d969150766",
                "levellist": null,
                "modulename": "new module test 1"
            },
            {
                "_id": "64142b5f1ff5fc2085084f1f",
                "levellist": null,
                "modulename": "new module test 1"
            },
            {
                "_id": "6414391d933646bd3b42bfdc",
                "levellist": null,
                "modulename": "new module test 1"
            }
        ]
    }
}
```
## localhost:8080/add-level-in-module
![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/4905ea58-94fc-4fe4-8839-157dfb249715)
```json
{
    "status": 201,
    "message": "success",
    "data": {
        "data": {
            "MatchedCount": 1,
            "ModifiedCount": 1,
            "UpsertedCount": 0,
            "UpsertedID": null
        }
    }
}
```

![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/9c47f65f-cbfa-441b-8759-f6e8fab34431)
```json
{
  "_id": {
    "$oid": "64625461e96f30a0f280d829"
  },
  "levelname": "test-level-1",
  "categorylist": [
    {
      "categoryname": "test-category-2",
      "problemlist": null
    }
  ]
}
```
## localhost:8080/add-challenge-in-category
![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/0ef440a4-9d40-4c23-861e-d24df9e654a3)
```json
{
  "_id": {
    "$oid": "641fedd8a0a7bf965d06e1f2"
  },
  "challengename": "test_challenge_1",
  "description": "description will contain problem description",
  "imageurl": "http://example-url-to-images",
  "score": 10,
  "nextchallengeurl": "http://url-to-next-challenge",
  "previouschallengeurl": "http://url-to-previous-challenge",
  "difficulty": "medium"
}
```
![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/c15c990a-1168-4d54-bc7e-d94b8d258bcf)


## localhost:8080/add-category-in-level
![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/faf469cd-52f3-4fc5-8220-a9d2eec2e703)
![image](https://github.com/asb-activity/ASB-design-go/assets/84634405/f46d9dba-6c79-4129-b024-cd26d92424bb)
```json
{
  "_id": {
    "$oid": "64625461e96f30a0f280d829"
  },
  "levelname": "test-level-1",
  "categorylist": [
    {
      "categoryname": "test-category-2",
      "problemlist": null
    },
    {
      "categoryname": "test-category-2",
      "problemlist": [
        {
          "challengename": "",
          "description": "description will contain problem description",
          "imageurl": "",
          "score": 0,
          "nextchallengeurl": "",
          "previouschallengeurl": "",
          "difficulty": "medium"
        },
        {
          "challengename": "test_challenge_1",
          "description": "description will contain problem description",
          "imageurl": "http://example-url-to-images",
          "score": 0,
          "nextchallengeurl": "http://url-to-next-challenge",
          "previouschallengeurl": "",
          "difficulty": "medium"
        },
        {
          "challengename": "test_challenge_1",
          "description": "description will contain problem description",
          "imageurl": "http://example-url-to-images",
          "score": 10,
          "nextchallengeurl": "http://url-to-next-challenge",
          "previouschallengeurl": "",
          "difficulty": "medium"
        }
      ]
    }
  ]
}
```

