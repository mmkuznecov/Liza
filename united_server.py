from flask import Flask
from flask import request
import json
from flask import jsonify
with open('/home/mikhail/server/document.json') as f:
    data = f.read()
    fgh = json.loads(data)
 
app = Flask(__name__)

@app.route('/json')
def json_info():
    user_id = request.args.get('id')
    #if key doesn't exist, returns None
    #framework = request.args['framework'] #if key doesn't exist, returns a 400, bad request error
    #website = request.args.get('website')

    return jsonify(fgh['users'][int(user_id)][user_id])

@app.route('/search_name')
def search_name():
	jsonFile = open("/home/mikhail/server/document.json", "r") # Open the JSON file for reading
	data = json.load(jsonFile) # Read the JSON into the buffer
	jsonFile.close()

	name = request.args.get('name')

	for i in range(len(data['users'])):
		if data['users'][i][str(i)][0]['name']==name:
			break
	return jsonify(i)

@app.route('/search')
def search():
	jsonFile = open("/home/mikhail/server/document.json", "r") # Open the JSON file for reading
	data = json.load(jsonFile) # Read the JSON into the buffer
	jsonFile.close()

	user = request.args.get('user')

	for i in range(len(data['users'])):
		if data['users'][i][str(i)][0]['name']==user:
			break
	return jsonify(data['users'][i][str(i)])

@app.route('/redact_profile')
def redact_profile():
	user_id = request.args.get('id')
	age = request.args.get('age')
	adress = request.args.get('adress')
	hours = request.args.get('hours')
	name = request.args.get('name')

	jsonFile = open("/home/mikhail/server/document.json", "r") # Open the JSON file for reading
	data = json.load(jsonFile) # Read the JSON into the buffer
	jsonFile.close()
	data['users'][user_id][str(user_id)][0].append([{'age':age,"adress": adress,'name':name,'hours':hours},])
	print(data)
	jsonFile = open("document.json", "w+")
	jsonFile.write(json.dumps(data))
	jsonFile.close()

	return jsonify(data)

@app.route('/reg')
def registration():
	tel = request.args.get('tel')
	email = request.args.get('email')
	password = request.args.get('password')
	name = request.args.get('name')
	surname = request.args.get('surname')
	otch = request.args.get('otch')

	
	jsonFile = open("/home/mikhail/server/document.json", "r") # Open the JSON file for reading
	data = json.load(jsonFile) # Read the JSON into the buffer
	jsonFile.close()


	data['users'].append({len(data['users']):[{'tel':tel,'email':email,'password':password,'name':surname+' '+name+' '+otch},]})
	#print(data)
	jsonFile = open("document.json", "w+")
	jsonFile.write(json.dumps(data))
	jsonFile.close()

	return str(len(data['users'])-1)


app.run('192.168.8.106')

@app.route('/stat')
def statistics():
	user_id = request.args.get('id')
	jsonFile = open("/home/mikhail/server/user_obraz.json", "r") # Open the JSON file for reading
	data = json.load(jsonFile) # Read the JSON into the buffer
	jsonFile.close()


	return jsonify(data['users'][int(user_id)][user_id][0]['hours']), #nify(len(data['users'][int(user_id)][user_id][0]['history_rating']))


@app.route('/json_event')
def json_event():
    event_id = request.args.get('event')
    #if key doesn't exist, returns None
    #framework = request.args['framework'] #if key doesn't exist, returns a 400, bad request error
    #website = request.args.get('website')

    return jsonify(fgh['events'][int(event_id)][event_id])

@app.route('/list_of_act')
def list_of_active():
	x=[]
	for i in range(len(fgh['events'])):
		if fgh['events'][i][str(i)][0]["status"]=="active":
			x.append(i)
	return jsonify(x)






app.run('0.0.0.0')