# casbin-auth
A role based access control example using Echo and Casbin

## Test

curl -H "X-User: alice" http://localhost:9000/admin/dashboard
curl -H "X-User: bob" http://localhost:9000/admin/dashboard
curl -H "X-User: charlie" http://localhost:9000/admin/dashboard