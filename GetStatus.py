import requests
import json
import re
import os

Headers = {"accept": "application/json"}

def AMP_Login(url, username, password):
    data = {  # This is use to fetch the sessionID token on AMP Server
        "username": username,
        "password": password,
        "token": "",
        "rememberMe": "true"
    }
    request = requests.post(url, headers=Headers, json=data)  # sending the request to AMP Server
    result = json.loads(request.text)  # converting the response from json format

    if result['success'] == True:  # make sure the request is success
        sessionID = result['sessionID']
    else:
        print('Login failed')
    return sessionID

def List_Instances(url, sessionID):
    sid = {  # SessionId we fetch from last request
        "SESSIONID": sessionID,
    }
    
    Instancesraw = requests.post(url, headers=Headers, json=sid)  # getting all the instances from AMP Server

    Instances = json.loads(Instancesraw.text)  # converting the response from json format
    try:
        Instances = Instances['result'][0]['AvailableInstances']
    except:
        print('Permission error')

    InstancesID = []  # Empty Array to store the InstanceID

    for instance in Instances:  # for loop to get all the InstanceID from the server
        InstancesInfo = {"InstanceID": None}  # defining the keys we want from the response
        InstancesInfo['InstanceID'] = instance['InstanceID']  # inserting the keys we want from the response
        InstancesID.append(InstancesInfo)  # Store the InstanceID in the InstancesID array
    return InstancesID

def Status_Instances(url, sessionID, InstanceID):
    InstancesStatuses = []  # Empty Array to store the Instance Status
    
    for status in InstanceID:  # with the instanceID we can get the status of all instances
        ID = status['InstanceID']  # for each instanceID we get the status
        InstanceStats = {  # body of the request
            "SESSIONID": sessionID,
            "InstanceId": ID
        }
        IStatsraw = requests.post(url, headers=Headers, json=InstanceStats)  # getting the status of the instance
        InstancesStatuses.append(json.loads(IStatsraw.text))  # storing the status in the InstancesStatuses array while converting the response from json format

    FullStatus = []  # Empty Array to store the Full Status

    for InstStatus in InstancesStatuses:  # for loop to format what data we want from the response
        if InstStatus['FriendlyName'] != 'ADS01':  # excluding the ADS01/controller instances
            Status = {"FriendlyName": None, 
                      "Game": None, 
                      "Running": None, 
                      "CPU Usage": None, 
                      "Memory Usage": None,
                      "Active Users": None, 
                      "Max Users": None
                      }  # defining the keys we want from the response
            Status['FriendlyName'] = InstStatus['FriendlyName']
            Status['Game'] = InstStatus['Module']
            Status['Running'] = InstStatus['Running']
            Status['CPU Usage'] = InstStatus['Metrics']['CPU Usage']['Percent']
            Status['Memory Usage'] = InstStatus['Metrics']['Memory Usage']['Percent']
            Status['Active Users'] = InstStatus['Metrics']['Active Users']['RawValue']
            Status['Max Users'] = InstStatus['Metrics']['Active Users']['MaxValue']
            FullStatus.append(Status)  # Store the Full Status in the FullStatus array
    return FullStatus

def Manage_Instances(url, sessionID, InstanceID, action ):
    match action:
        case 'start':
            return url
        case 'stop':
            return url
        case 'restart':
            return url

def main(args):  # start of the Function
    AMPUrl = args.get("AMPUrl", os.environ['AMPUrl'])
    AMPUser = args.get("AMPUser", os.environ['AMPUser'])
    AMPPass = args.get("AMPPass", os.environ['AMPPass'])
    urlpattern = '(http|https)://\S*'  # pattern use to validate if url have http or https
    IsHTTP = re.match(urlpattern, AMPUrl)  # match the url with the pattern
    if IsHTTP:  # if url have http or https just add the paths
        Login = AMPUrl + '/API/Core/Login'
        GInstances = AMPUrl + '/API/ADSModule/GetInstances'
        UrlStatInstance = AMPUrl + '/API/ADSModule/GetInstance'
    else:  # else add the protocol and paths
        Login = 'https://' + AMPUrl + '/API/Core/Login'
        GInstances = 'https://' + AMPUrl + '/API/ADSModule/GetInstances'
        UrlStatInstance = 'https://' + AMPUrl + '/API/ADSModule/GetInstance'
    
    Sessionid = AMP_Login(Login, AMPUser, AMPPass)  # fetching the sessionID from the AMP Server
    
    All_Instances = List_Instances(GInstances, Sessionid)  # fetching all the instances from the AMP Server
    
    Status = Status_Instances(UrlStatInstance, Sessionid, All_Instances)  # fetching the status of all the instances from the AMP Server
        
    

    return {'body': Status}  # return the FullStatus array