# fetch-rewards-receipt-processor
### Run instructions:
In commands.txt, there are some curl commands that I pasted for easy use. This is an example run:
```
go run receipt-processor.go
```

```
// This is to test the process receipts endpoint
// Replace 'data' with any other json file contents

curl http://localhost:5670/receipts/process --include --header "Content-Type: application/json" --request "POST" --data '{
    "retailer": "M&M Corner Market",
    "purchaseDate": "2022-03-20",
    "purchaseTime": "14:33",
    "items": [
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },{
        "shortDescription": "Gatorade",
        "price": "2.25"
      },{
        "shortDescription": "Gatorade",
        "price": "2.25"
      },{
        "shortDescription": "Gatorade",
        "price": "2.25"
      }
    ],
    "total": "9.00"
  }'


  // or
  curl http://localhost:5670/receipts/process --include --header "Content-Type: application/json" --request "POST" --data 'REPLACETHIS'

```
```
// Use the printed out id (should look like: "57497390-461d-4c48-ba81-0f5c582085c3")
// to replace the {id} and check other entries

curl http://localhost:5670/receipts/85b3063b-c31f-423d-a942-df5821d8f762/points -v

// or
curl http://localhost:5670/receipts/{id}/points -v
```


### Sources
I used these sources in addition to the Go documentation (https://pkg.go.dev/encoding/json) to build my project:

https://medium.com/@briankworld/building-a-restful-api-with-go-a-step-by-step-guide-d17e69f004a7

https://go.dev/doc/tutorial/web-service-gin

https://www.jetbrains.com/guide/go/tutorials/rest_api_series/stdlib/

https://pkg.go.dev/encoding/json#Unmarshal