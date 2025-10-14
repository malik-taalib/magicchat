#!/bin/bash

echo "Creating users via API..."

curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username":"sarah_adventures","email":"sarah@magicchat.com","password":"password123","display_name":"Sarah Adventures"}' && echo ""

curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username":"chef_mike","email":"mike@magicchat.com","password":"password123","display_name":"Chef Mike"}' && echo ""

curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username":"fitness_emma","email":"emma@magicchat.com","password":"password123","display_name":"Emma Fitness"}' && echo ""

curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username":"comedy_jay","email":"jay@magicchat.com","password":"password123","display_name":"Jay Comedy"}' && echo ""

curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username":"tech_alex","email":"alex@magicchat.com","password":"password123","display_name":"Tech Alex"}' && echo ""

curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username":"music_lisa","email":"lisa@magicchat.com","password":"password123","display_name":"Lisa Music"}' && echo ""

echo "Done creating users!"
