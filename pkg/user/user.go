package user

import(
"encoding/json"
"errors"
"github.com/aws/aws-lambda-go/events"
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/service/dynamodb"
"github.com/aws/aws-sdk-go/service/dynamodb/attribute"
"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)


var(
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorInvalidUserData ="failed to fetch report"
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorInvalidEmail = "invalid email"
	ErrorCouldNotMarshalItem = "could not marshal item"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item"
	ErrorUserAlreadyExisted = "user.User already existed"
	ErrorUserDoesNotExist = "user.User does not exist"
)

type User struct{
	Email 			string 'json:"email"'
	FirstName 		string 'json:"firstName"'
	LastName		string 'json:"lastName"'
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error,){
	input := &dynamodb.GetItemInput{
		Key:map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		}
		TableName: aws.String(tableName),
	}
	result,err := dynaClient.GetItem(input)
	if err !=nil{
		return nil,errors.New(ErrorFailedToFetchRecord)
	}
	item  := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err !=nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	return item, nil
		
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User,err,){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName)
	}
	result,err:=dynaClient.Scan(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshalMap(result.Items,item)
	return item,nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)
(
	*User,
	error,
){
	var u  User

	if err := json.UnmarshalMap([]byte(req,body),&u); err!=nil{
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}
	//check if user already exist
	currentUser,_ :=FetchUser(u.email, tableName, dynaClient)
	if currentUser !=nil && len(currentUser.Email)!= 0{
		return nil, errors.New(ErrorUserAlreadyExisted)
	}

	av, err := dynamodbattribute.marshalMap(u)
	if err!=nil{
		return nil, errors.New(ErrorCouldNotDeleteItem)
	}

	input: &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err= dynaClient.PutItem(input)
	if err!=nil{
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)
(
	*User,
	error,
){
	var u User
	if err := json.Unmarshal([]byte(req.Body),&u); err!=nil{
		return nil, errors.New(ErrorInvalidEmail)
	}
	currentUser,_ := FetchUser(u.email,tableName,dynaClient)
	if currentUser !=nil && len(currentUser.Email)==0{
		return nil, errors.New(ErrorUserDoesNotExist)
	}
	av, err := dynamodbattribute.MarshalMap(u)
	if err!=nil{
		return nil, errors.New(ErrorCouldNotMarshalItem) 
	}
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	},

	_, err = dynaClient.PutItem(input)
	if err != nil{
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func DeleteUser(req events.APIGateWayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)error,{
	email := req.QueryStringParameters("email")
	input := &dynamodb.DeleteItemInput{
		Key:map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		}
		TableName: aws.String(tableName),
	}
	_,err:=dynaClient.DeleteItem(input)
	if err!=nil{
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}