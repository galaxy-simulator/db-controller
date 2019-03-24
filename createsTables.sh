#!/usr/bin/env bash

echo "Recreating the tables..."
curl -X POST http://localhost:8080/createNodesTable                                                                            
curl -X POST http://localhost:8080/createStarsTable                                                                            
    

