# Documentation for Boarding API

The source code for this project is designed in a modular manner defining controllers and models modules. I described the routes and route handlers within the controllers and describe the data types and datastore functions within the models. 

#### Running the Application
---
- Begin by following the instructions of the README.md of "infrastructure" directory
- From the parent directory "sd-automation" change into "application" directory - `cd application` 
- Run the "golibraries.sh" script to install required project libraries in GOPATH
 
<b> Steps </b>  
  1. `./boardingapi`  #Run app locally then test. See Testing Endpoints below 
  2. `docker build . -t boardingapi:v1.0.0`   #Build the docker image
  3. `docker run -d --network host --name boardingapi boardingapi:v1.0.0` #Run the container in detached mode 
  4. `docker login`    #Provide your dockerhub user name and password
  5. `docker tag boardingapi:v1.0.0 knoxknot/boardingapi:v1.0.0` #Tag image for push to repository
  6. `docker push  knoxknot/boardingapi:v1.0.0`   #Push image to repository
   
#### Testing Endpoints  
---
`curl -X GET http://localhost:8080/api/v1/people`  # Display everyone on board.  
`curl -X GET http://localhost:8080/api/v1/people/5e1efe5ebd182108b22242d3`  # Display a unique person on board. Replace 5e1efe5ebd182108b22242d3 with selected uuid.  
`curl -X POST localhost:8080/api/v1/people -d '{"survived": true,"passengerClass": 1,"name": "Mr. Grace Chukwusom","sex": "male","age": 25,"siblingsOrSpousesAboard": 0,"parentsOrChildrenAboard": 0,"fare": 8.12}' -H "Content-Type: application/json"`  # Enter details of a person that boarded.  
`curl -i -X PUT http://localhost:8080/api/v1/people/5e1efe5ebd182108b22242d3 -d '{"survived": false, "sex": "female","siblingsOrSpousesAboard": 3,"age": 35}' -H "Content-Type: application/json"`    # Update details of person on board. Replace 5e1efe5ebd182108b22242d3 with selected uuid.  
`curl -X DELETE http://localhost:8080/api/v1/people/5e1efe5ebd182108b22242d3`    # Detele a person that unboarded. Replace 5e1efe5ebd182108b22242d3 with selected uuid.


