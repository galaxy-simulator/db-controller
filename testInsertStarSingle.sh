#!/usr/bin/env bash

echo "Recreating the tables..."
curl -X POST http://localhost:8081/createNodesTable                                                                                
curl -X POST http://localhost:8081/createStarsTable                                                                                
echo "Done"

echo "Deleting preexisting tables..."
curl -X POST http://localhost:8081/deleteStars
curl -X POST http://localhost:8081/deleteNodes
echo "Done"

echo "Inserting all stars from teststars.csv..."
curl -X POST --data "filename=100" http://localhost:8081/insertStarList
echo "Done"

echo "Inserting a list of stars..."
echo "Done"

echo "Getting the forest representation of the tree..."
curl -X GET http://localhost:8081/genForestTree?index=1
echo "Done"

echo "Getting a list of all stars..."
curl -X GET http://localhost:8081/starslist/csv                                                                                    
echo "Done"
