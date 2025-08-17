curl -X POST -H "Content-Type: application/json" -d '{"email": "test@test.com", "password": "password", "name": "Test User"}' http://localhost:8080/api/v1/register

curl -X POST -H "Content-Type: application/json" -d '{"email": "test@test.com", "password": "password"}' http://localhost:8080/api/v1/login


curl -X POST -H "Content-Type: application/json" -d '{"name": "Московская область"}' http://localhost:8080/api/v1/regions
curl -X POST -H "Content-Type: application/json" -d '{"name": "Солнечногорский"}' http://localhost:8080/api/v1/regions/094a5f2f-bdb8-4861-aabb-2045fae11764/land-parcels
curl -X POST -H "Content-Type: application/json" -d '{"name": "Южная", "type": "Поликарбонат"}' http://localhost:8080/api/v1/land-parcels/8a570b35-16e3-4ed1-b093-cf1d54a0e203/greenhouses
curl -X POST -H "Content-Type: application/json" -d '{"name": "A-1", "size": "1x2m", "camera_url": "rtsp://1.2.3.4/stream1"}' http://localhost:8080/api/v1/greenhouses/72674c47-98be-468b-bb1e-952eff9b130a/plots


## Аренда грядки
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU0MjUxNDAsImlhdCI6MTc1NTQyNDI0MCwic3ViIjoiMWM2Y2IzMDMtMzMyNy00N2QzLThhODgtMjg1OWM4M2FjZTkzIn0.7m-4LW9c4txOW3UqfFpYX57m78E09dp-hIhOwJpDyOU" http://localhost:8080/api/v1/plots/3735ad75-7c19-47d3-9456-8762af6daaa1/lease