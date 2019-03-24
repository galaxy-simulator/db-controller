#!/usr/bin/env bash

################################################################################
# init 
################################################################################

# create tables 
printf "Recreating the tables... "
curl -X POST http://db-controller.docker.localhost/createNodesTable
curl -X POST http://db-controller.docker.localhost/createStarsTable
printf "Done   "
read -n 1 -s -r -p "Press any key to continue (delete preexisting tables)"

# delete preexisting table entries
printf "\nDeleting preexisting tables... "
curl -X POST http://db-controller.docker.localhost/deleteStars
curl -X POST http://db-controller.docker.localhost/deleteNodes
printf "Done   "
read -n 1 -s -r -p "Press any key to continue (create a new tree)"

# create a new tree 
printf "\nCreating a new tree... "
curl -X POST --data "w=100000000000" http://db-controller.docker.localhost/newTree
printf "Done   "
read -n 1 -s -r -p "Press any key to continue (insert stars from 10.csv)"

################################################################################
# insert 
################################################################################

# insert all stars from csv
printf "\nInserting all stars from 10.csv..."
curl -X POST --data "filename=10" http://db-controller.docker.localhost/insertStarList
printf "Done   "
read -n 1 -s -r -p "Press any key to continue (update total mass) "

################################################################################
# update
################################################################################

printf "\nInserting all stars from 10.csv..."
curl -X POST --data "index=1" http://db-controller.docker.localhost/updateTotalMass
printf "Done   "
read -n 1 -s -r -p "Press any key to continue (update center of mass) "

printf "\nInserting all stars from 10.csv..."
curl -X POST --data "index=1" http://db-controller.docker.localhost/updateCenterOfMass
printf "Done   "
read -n 1 -s -r -p "Press any key to continue (forest representation) "

################################################################################
# forest 
################################################################################

printf "\nGetting the forest representation of the tree..."
curl -X GET http://db-controller.docker.localhost/genForestTree?index=1
printf "Done    "
read -n 1 -s -r -p "Press any key to continue (get a list of all stars)"

################################################################################
# list of stars 
################################################################################

printf "\nGetting a list of all stars... "
curl -X GET http://db-controller.docker.localhost/starslist/csv
printf "Done    "
read -n 1 -s -r -p "Press any key to continue"
